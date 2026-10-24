import type { MetadataRoute } from "next";
import { getBlogPosts, type BlogPost } from "../lib/blogs";
import { siteUrl } from "../lib/site";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  let posts: BlogPost[] = [];
  try {
    posts = await getBlogPosts();
  } catch {
    posts = [];
  }
  // Use the newest post date rather than "now" so lastmod only changes when
  // content does; an always-fresh lastmod teaches crawlers to ignore it.
  const latest = posts.reduce<Date | undefined>((max, post) => {
    const date = new Date(post.updated_at);
    return !max || date > max ? date : max;
  }, undefined);
  return [
    {
      url: siteUrl,
      lastModified: latest,
      changeFrequency: "monthly",
      priority: 1,
    },
    {
      url: `${siteUrl}/blogs`,
      lastModified: latest,
      changeFrequency: "weekly",
      priority: 0.8,
    },
    ...posts.map((post) => ({
      url: `${siteUrl}/blogs/${post.slug}`,
      lastModified: new Date(post.updated_at),
      changeFrequency: "monthly" as const,
      priority: 0.7,
    })),
  ];
}
