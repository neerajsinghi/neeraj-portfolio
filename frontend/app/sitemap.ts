import type { MetadataRoute } from "next";
import { getBlogPosts, type BlogPost } from "../lib/blogs";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  let posts: BlogPost[] = [];
  try {
    posts = await getBlogPosts();
  } catch {
    posts = [];
  }
  return [
    {
      url: "https://neerajsinghi.com",
      lastModified: new Date(),
      changeFrequency: "monthly",
      priority: 1,
    },
    {
      url: "https://neerajsinghi.com/blogs",
      lastModified: new Date(),
      changeFrequency: "weekly",
      priority: 0.8,
    },
    ...posts.map((post) => ({
      url: `https://neerajsinghi.com/blogs/${post.slug}`,
      lastModified: new Date(post.updated_at),
      changeFrequency: "monthly" as const,
      priority: 0.7,
    })),
  ];
}