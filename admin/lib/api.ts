export type Status = "draft" | "scheduled" | "published" | "archived";

export type BlogPost = {
  id: string;
  slug: string;
  title: string;
  description: string;
  content_markdown: string;
  tags: string[];
  status: Status;
  version: number;
  updated_at: string;
  scheduled_at?: string;
  published_at?: string;
};

export type BlogInput = Pick<BlogPost, "slug" | "title" | "description" | "content_markdown" | "status" | "version"> & {
  tags: string[];
  scheduled_at?: string;
};

const apiBase = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

async function request<T>(path: string, token: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBase}${path}`, {
    ...init,
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json", ...init?.headers },
  });
  if (response.status === 401) throw new Error("SESSION_EXPIRED");
  if (!response.ok) {
    const body = (await response.json().catch(() => ({}))) as { error?: string };
    throw new Error(body.error || `Request failed with ${response.status}`);
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

export async function listPosts(token: string) {
  const result = await request<{ posts: BlogPost[] }>("/api/admin/blogs?limit=100", token);
  return result.posts;
}

export function createPost(token: string, input: BlogInput) {
  return request<BlogPost>("/api/admin/blogs", token, { method: "POST", body: JSON.stringify(input) });
}

export function updatePost(token: string, id: string, input: BlogInput) {
  return request<BlogPost>(`/api/admin/blogs/${id}`, token, { method: "PUT", body: JSON.stringify(input) });
}

export function deletePost(token: string, id: string) {
  return request<void>(`/api/admin/blogs/${id}`, token, { method: "DELETE" });
}