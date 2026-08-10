import type { Metadata } from "next";
import Link from "next/link";
import BlogHeader from "../../components/BlogHeader";
import { getBlogPosts, type BlogPost } from "../../lib/blogs";

export const metadata: Metadata = {
    title: "Engineering Blogs | Neeraj Singhi",
    description: "Practical writing about Go, backend systems, AWS, distributed architecture, and applied AI.",
    alternates: { canonical: "/blogs" },
};

// Forces this route to always respect the 60s window (and on-demand
// revalidateTag calls) instead of being eligible for indefinite Full Route
// Cache static optimization.
export const revalidate = 60;

export default async function BlogsPage() {
    let posts: BlogPost[] = [];
    let unavailable = false;
    try {
        posts = await getBlogPosts();
    } catch {
        unavailable = true;
    }

    return (
        <>
            <BlogHeader />
            <main className="blog-page">
                <div className="wrap">
                    <div className="blog-intro">
                        <span className="eyebrow">Engineering notes</span>
                        <h1>Systems, software, and the decisions between them.</h1>
                        <p>Practical writing about Go, distributed backends, cloud infrastructure, and applied AI.</p>
                    </div>

                    {unavailable ? (
                        <div className="blog-empty">The blog service is temporarily unavailable.</div>
                    ) : posts.length === 0 ? (
                        <div className="blog-empty">The first article is being prepared.</div>
                    ) : (
                        <div className="blog-list">
                            {posts.map((post) => (
                                <article className="blog-card" key={post.slug}>
                                    <div className="blog-card-meta">
                                        <time dateTime={post.published_at || post.updated_at}>{formatDate(post.published_at || post.updated_at)}</time>
                                        <span>{post.tags.join(" · ")}</span>
                                    </div>
                                    <h2><Link href={`/blogs/${post.slug}`}>{post.title}</Link></h2>
                                    <p>{post.description}</p>
                                    <Link className="blog-read" href={`/blogs/${post.slug}`}>Read article <span aria-hidden="true">→</span></Link>
                                </article>
                            ))}
                        </div>
                    )}
                </div>
            </main>
        </>
    );
}

function formatDate(value: string) {
    return new Intl.DateTimeFormat("en-IN", { dateStyle: "medium", timeZone: "UTC" }).format(new Date(value));
}