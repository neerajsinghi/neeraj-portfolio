import { API_BASE } from "./profile";
import type { Repo } from "../types";

export async function fetchGithubRepos(): Promise<{ user: string; repos: Repo[] }> {
  const res = await fetch(`${API_BASE}/api/github`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

export async function streamChat(
  messages: { role: string; content: string }[],
  onEvent: (event: string, data: unknown) => void
): Promise<void> {
  const res = await fetch(`${API_BASE}/api/chat`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ messages }),
  });
  if (!res.ok || !res.body) throw new Error(`HTTP ${res.status}`);

  const reader = res.body.getReader();
  const decoder = new TextDecoder();
  let buf = "";
  while (true) {
    const { value, done } = await reader.read();
    if (done) break;
    buf += decoder.decode(value, { stream: true });
    let idx: number;
    while ((idx = buf.indexOf("\n\n")) >= 0) {
      const frame = buf.slice(0, idx);
      buf = buf.slice(idx + 2);
      const lines = frame.split("\n");
      let ev = "message";
      let dataStr = "";
      for (const ln of lines) {
        if (ln.startsWith("event:")) ev = ln.slice(6).trim();
        else if (ln.startsWith("data:")) dataStr += ln.slice(5).trim();
      }
      if (!dataStr) continue;
      let data: unknown = {};
      try {
        data = JSON.parse(dataStr);
      } catch {
        continue;
      }
      onEvent(ev, data);
    }
  }
}
