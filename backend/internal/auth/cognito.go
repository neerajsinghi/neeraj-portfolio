package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"neeraj-portfolio/backend/internal/blog"
)

var ErrUnauthorized = errors.New("unauthorized")

type Authenticator interface {
	Authenticate(*http.Request) (blog.Principal, error)
}

type Cognito struct {
	issuer     string
	clientID   string
	httpClient *http.Client
	mu         sync.RWMutex
	keys       map[string]*rsa.PublicKey
	keysUntil  time.Time
}

type jwksResponse struct {
	Keys []struct {
		KID string `json:"kid"`
		KTY string `json:"kty"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

func NewCognito(region, userPoolID, clientID string) (*Cognito, error) {
	if region == "" || userPoolID == "" || clientID == "" {
		return nil, errors.New("COGNITO_REGION, COGNITO_USER_POOL_ID, and COGNITO_CLIENT_ID are required")
	}
	return &Cognito{
		issuer:     "https://cognito-idp." + region + ".amazonaws.com/" + userPoolID,
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		keys:       make(map[string]*rsa.PublicKey),
	}, nil
}

func (cognito *Cognito) Authenticate(request *http.Request) (blog.Principal, error) {
	header := request.Header.Get("Authorization")
	tokenString, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || tokenString == "" {
		return blog.Principal{}, ErrUnauthorized
	}
	claims := jwt.MapClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(cognito.issuer),
		jwt.WithExpirationRequired(),
	)
	token, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, ErrUnauthorized
		}
		return cognito.key(request.Context(), kid)
	})
	if err != nil || !token.Valid {
		return blog.Principal{}, ErrUnauthorized
	}
	if claims["token_use"] != "access" || claims["client_id"] != cognito.clientID {
		return blog.Principal{}, ErrUnauthorized
	}
	subject, _ := claims["sub"].(string)
	if subject == "" {
		return blog.Principal{}, ErrUnauthorized
	}
	email, _ := claims["email"].(string)
	if email == "" {
		email, _ = claims["username"].(string)
	}
	return blog.Principal{Subject: subject, Email: email, Roles: claimRoles(claims["cognito:groups"])}, nil
}

func claimRoles(value any) map[blog.Role]bool {
	roles := make(map[blog.Role]bool)
	var groups []string
	switch typed := value.(type) {
	case []any:
		for _, group := range typed {
			if name, ok := group.(string); ok {
				groups = append(groups, name)
			}
		}
	case []string:
		groups = typed
	}
	for _, group := range groups {
		switch blog.Role(group) {
		case blog.RoleAdmin, blog.RoleEditor:
			roles[blog.Role(group)] = true
		}
	}
	return roles
}

func (cognito *Cognito) key(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	cognito.mu.RLock()
	key, fresh := cognito.keys[kid], time.Now().Before(cognito.keysUntil)
	cognito.mu.RUnlock()
	if key != nil && fresh {
		return key, nil
	}
	if err := cognito.refreshKeys(ctx); err != nil {
		return nil, err
	}
	cognito.mu.RLock()
	defer cognito.mu.RUnlock()
	key = cognito.keys[kid]
	if key == nil {
		return nil, ErrUnauthorized
	}
	return key, nil
}

func (cognito *Cognito) refreshKeys(ctx context.Context) error {
	cognito.mu.Lock()
	defer cognito.mu.Unlock()
	if time.Now().Before(cognito.keysUntil) && len(cognito.keys) > 0 {
		return nil
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cognito.issuer+"/.well-known/jwks.json", nil)
	if err != nil {
		return err
	}
	response, err := cognito.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch Cognito keys: HTTP %d", response.StatusCode)
	}
	var document jwksResponse
	if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
		return err
	}
	keys := make(map[string]*rsa.PublicKey, len(document.Keys))
	for _, jwk := range document.Keys {
		if jwk.KTY != "RSA" {
			continue
		}
		modulus, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			continue
		}
		exponentBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			continue
		}
		exponent := 0
		for _, value := range exponentBytes {
			exponent = exponent<<8 + int(value)
		}
		keys[jwk.KID] = &rsa.PublicKey{N: new(big.Int).SetBytes(modulus), E: exponent}
	}
	if len(keys) == 0 {
		return errors.New("Cognito returned no usable signing keys")
	}
	cognito.keys, cognito.keysUntil = keys, time.Now().Add(6*time.Hour)
	return nil
}
