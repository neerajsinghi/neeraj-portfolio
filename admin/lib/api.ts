export type Status = "draft" | "scheduled" | "published" | "archived";

export type BlogPost = {
  id: string;
  slug: string;
  title: string;
  description: string;
  content_markdown: string;
  linkedin_post?: string;
  social_post?: string;
  tags: string[];
  status: Status;
  version: number;
  updated_at: string;
  scheduled_at?: string;
  published_at?: string;
  publish_devto?: boolean;
  publish_linkedin?: boolean;
  devto_url?: string;
  devto_published_at?: string;
  linkedin_url?: string;
  linkedin_published_at?: string;
};

export type BlogInput = Pick<BlogPost, "slug" | "title" | "description" | "content_markdown" | "linkedin_post" | "social_post" | "status" | "version"> & {
  tags: string[];
  scheduled_at?: string;
  publish_devto?: boolean;
  publish_linkedin?: boolean;
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
  const result = await request<{ posts: BlogPost[] }>("/api/v1/admin/blogs?limit=100", token);
  return result.posts;
}

export function createPost(token: string, input: BlogInput) {
  return request<BlogPost>("/api/v1/admin/blogs", token, { method: "POST", body: JSON.stringify(input) });
}

export function updatePost(token: string, id: string, input: BlogInput) {
  return request<BlogPost>(`/api/v1/admin/blogs/${id}`, token, { method: "PUT", body: JSON.stringify(input) });
}

export function deletePost(token: string, id: string) {
  return request<void>(`/api/v1/admin/blogs/${id}`, token, { method: "DELETE" });
}

export type ExternalPlatform = "devto" | "linkedin";

export type PublishResult = {
  platform: ExternalPlatform;
  id?: string;
  url?: string;
};

export async function publishExternally(token: string, input: BlogInput, platforms: ExternalPlatform[]) {
  const response = await request<{ results: PublishResult[] }>("/api/v1/admin/publish", token, {
    method: "POST",
    body: JSON.stringify({
      title: input.title,
      slug: input.slug,
      description: input.description,
      article_markdown: input.content_markdown,
      linkedin_post: input.linkedin_post || "",
      social_post: input.social_post || "",
      tags: input.tags,
      canonical_url: `https://neerajsinghi.com/blogs/${input.slug}`,
      platforms,
    }),
  });
  return response.results;
}