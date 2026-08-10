export type BlogPost = {
  slug: string;
  title: string;
  description: string;
  content_markdown?: string;
  tags: string[];
  published_at?: string;
  updated_at: string;
};

const apiBase = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export async function getBlogPosts(): Promise<BlogPost[]> {
  const response = await fetch(`${apiBase}/api/v1/blogs?limit=100`, { next: { revalidate: 60, tags: ["blogs"] } });
  if (!response.ok) {
    throw new Error(`Blog API returned ${response.status}`);
  }
  const body = (await response.json()) as { posts: BlogPost[] };
  return body.posts;
}

export async function getBlogPost(slug: string): Promise<BlogPost | null> {
  const response = await fetch(`${apiBase}/api/v1/blogs/${encodeURIComponent(slug)}`, { next: { revalidate: 60, tags: ["blogs"] } });
  if (response.status === 404) return null;
  if (!response.ok) {
    throw new Error(`Blog API returned ${response.status}`);
  }
  return (await response.json()) as BlogPost;
}