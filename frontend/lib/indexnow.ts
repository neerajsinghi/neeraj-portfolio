import type { BlogPost } from "./blogs";
import { siteUrl } from "./site";

// IndexNow (https://www.indexnow.org) tells Bing, Yandex, Seznam and others
// about changed URLs immediately; Bing's index also feeds ChatGPT search and
// Copilot. The key is public by design and must match the file served at
// /<key>.txt from public/.
const indexNowKey = "aa1f87905476d04fba78bdb170e832c8";
const apiBase = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
const recentWindowMs = 24 * 60 * 60 * 1000;

// Best effort: submits the homepage, /blogs and any post updated in the last
// day. Never throws. Only runs against the production domain.
export async function notifyIndexNow(): Promise<void> {
  if (!apiBase.startsWith("https://")) return;
  try {
    const response = await fetch(`${apiBase}/api/v1/blogs?limit=100`, { cache: "no-store" });
    const posts = response.ok ? ((await response.json()) as { posts: BlogPost[] }).posts : [];
    const cutoff = Date.now() - recentWindowMs;
    const urlList = [
      `${siteUrl}/`,
      `${siteUrl}/blogs`,
      ...posts.filter((post) => new Date(post.updated_at).getTime() >= cutoff).map((post) => `${siteUrl}/blogs/${post.slug}`),
    ];
    const result = await fetch("https://api.indexnow.org/indexnow", {
      method: "POST",
      headers: { "Content-Type": "application/json; charset=utf-8" },
      body: JSON.stringify({
        host: new URL(siteUrl).host,
        key: indexNowKey,
        keyLocation: `${siteUrl}/${indexNowKey}.txt`,
        urlList,
      }),
      signal: AbortSignal.timeout(10_000),
    });
    if (!result.ok && result.status !== 202) {
      console.error(`indexnow: unexpected status ${result.status}`);
    }
  } catch (error) {
    console.error("indexnow: submit failed", error);
  }
}
