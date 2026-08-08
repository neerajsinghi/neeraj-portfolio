const storageKey = "portfolio-admin-session";
const verifierKey = "portfolio-admin-pkce";

type Session = {
  accessToken: string;
  expiresAt: number;
};

function config() {
  const domain = process.env.NEXT_PUBLIC_COGNITO_DOMAIN;
  const clientID = process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID;
  const redirectURI = process.env.NEXT_PUBLIC_COGNITO_REDIRECT_URI;
  if (!domain || !clientID || !redirectURI) throw new Error("Cognito environment is incomplete");
  return { domain: domain.replace(/\/$/, ""), clientID, redirectURI };
}

function base64URL(bytes: Uint8Array) {
  return btoa(String.fromCharCode(...bytes)).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

export async function beginLogin() {
  const { domain, clientID, redirectURI } = config();
  const verifierBytes = crypto.getRandomValues(new Uint8Array(64));
  const verifier = base64URL(verifierBytes);
  const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(verifier));
  sessionStorage.setItem(verifierKey, verifier);
  const query = new URLSearchParams({
    client_id: clientID,
    response_type: "code",
    scope: "openid email profile",
    redirect_uri: redirectURI,
    code_challenge_method: "S256",
    code_challenge: base64URL(new Uint8Array(digest)),
  });
  window.location.assign(`${domain}/oauth2/authorize?${query}`);
}

export async function completeLogin(code: string) {
  const { domain, clientID, redirectURI } = config();
  const verifier = sessionStorage.getItem(verifierKey);
  if (!verifier) throw new Error("Login verifier expired; start again");
  const body = new URLSearchParams({
    grant_type: "authorization_code",
    client_id: clientID,
    code,
    redirect_uri: redirectURI,
    code_verifier: verifier,
  });
  const response = await fetch(`${domain}/oauth2/token`, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body,
  });
  if (!response.ok) throw new Error("Cognito rejected the authorization code");
  const tokens = (await response.json()) as { access_token: string; expires_in: number };
  const session: Session = { accessToken: tokens.access_token, expiresAt: Date.now() + tokens.expires_in * 1000 };
  sessionStorage.setItem(storageKey, JSON.stringify(session));
  sessionStorage.removeItem(verifierKey);
}

export function getSession(): Session | null {
  const raw = sessionStorage.getItem(storageKey);
  if (!raw) return null;
  try {
    const session = JSON.parse(raw) as Session;
    if (!session.accessToken || session.expiresAt <= Date.now() + 30_000) {
      sessionStorage.removeItem(storageKey);
      return null;
    }
    return session;
  } catch {
    sessionStorage.removeItem(storageKey);
    return null;
  }
}

export function getRoles(accessToken: string): string[] {
  try {
    const payload = JSON.parse(atob(accessToken.split(".")[1].replace(/-/g, "+").replace(/_/g, "/"))) as { "cognito:groups"?: string[] };
    return payload["cognito:groups"] || [];
  } catch {
    return [];
  }
}

export function logout() {
  const { domain, clientID, redirectURI } = config();
  sessionStorage.removeItem(storageKey);
  const query = new URLSearchParams({ client_id: clientID, logout_uri: new URL("/", redirectURI).toString() });
  window.location.assign(`${domain}/logout?${query}`);
}