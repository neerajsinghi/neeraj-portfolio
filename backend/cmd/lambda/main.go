package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	anthropicprovider "neeraj-portfolio/backend/internal/llm/anthropic"
	"neeraj-portfolio/backend/internal/server"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var handler http.Handler

func init() {
	prov := anthropicprovider.New()
	handler = server.NewHandler(prov)
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, evt events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	body, err := decodeBody(evt.Body, evt.IsBase64Encoded)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusBadRequest, Body: "invalid request body"}, nil
	}

	path := evt.RawPath
	if path == "" {
		path = "/"
	}
	if evt.RawQueryString != "" {
		path += "?" + evt.RawQueryString
	}

	method := http.MethodGet
	if evt.RequestContext.HTTP.Method != "" {
		method = evt.RequestContext.HTTP.Method
	}

	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: "request build failed"}, nil
	}

	for k, v := range evt.Headers {
		req.Header.Set(k, v)
	}
	if evt.RequestContext.HTTP.SourceIP != "" {
		req.RemoteAddr = evt.RequestContext.HTTP.SourceIP
	}

	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)
	resp := rw.Result()
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: "response read failed"}, nil
	}

	headers := make(map[string]string, len(resp.Header))
	for k, vals := range resp.Header {
		if len(vals) > 0 {
			headers[k] = strings.Join(vals, ",")
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode:      resp.StatusCode,
		Headers:         headers,
		Body:            string(respBody),
		IsBase64Encoded: false,
	}, nil
}

func decodeBody(body string, isBase64 bool) ([]byte, error) {
	if !isBase64 {
		return []byte(body), nil
	}
	return base64.StdEncoding.DecodeString(body)
}
