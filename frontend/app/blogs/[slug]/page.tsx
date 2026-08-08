import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import BlogHeader from "../../../components/BlogHeader";
import { getBlogPost } from "../../../lib/blogs";

type Props = { params: { slug: string } };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
    const post = await getBlogPost(params.slug);
    if (!post) return { title: "Blog not found | Neeraj Singhi" };
    return {
        title: `${post.title} | Neeraj Singhi`,
        description: post.description,
        alternates: { canonical: `/blogs/${post.slug}` },
        openGraph: {
            type: "article",
            url: `/blogs/${post.slug}`,
            title: post.title,
            description: post.description,
            publishedTime: post.published_at,
            modifiedTime: post.updated_at,
            tags: post.tags,
        },
    };
}

export default async function BlogDetailPage({ params }: Props) {
    const post = await getBlogPost(params.slug);
    if (!post) notFound();

    const articleData = {
        "@context": "https://schema.org",
        "@type": "BlogPosting",
        headline: post.title,
        description: post.description,
        datePublished: post.published_at,
        dateModified: post.updated_at,
        mainEntityOfPage: `https://neerajsinghi.com/blogs/${post.slug}`,
        author: { "@type": "Person", name: "Neeraj Singhi", url: "https://neerajsinghi.com" },
    };

    return (
        <>
            <BlogHeader />
            <main className="blog-page">
                <article className="wrap blog-article">
                    <Link className="blog-back" href="/blogs">← All blogs</Link>
                    <header className="blog-article-head">
                        <div className="blog-card-meta">
                            <time dateTime={post.published_at || post.updated_at}>{formatDate(post.published_at || post.updated_at)}</time>
                            <span>{post.tags.join(" · ")}</span>
                        </div>
                        <h1>{post.title}</h1>
                        <p>{post.description}</p>
                    </header>
                    <div className="blog-prose">
                        <ReactMarkdown
                            remarkPlugins={[remarkGfm]}
                            components={{
                                a: ({ href, children }) => {
                                    const external = href?.startsWith("http");
                                    return <a href={href} target={external ? "_blank" : undefined} rel={external ? "noopener noreferrer" : undefined}>{children}</a>;
                                },
                            }}
                        >
                            {post.content_markdown || ""}
                        </ReactMarkdown>
                    </div>
                    <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(articleData).replace(/</g, "\\u003c") }} />
                </article>
            </main>
        </>
    );
}

function formatDate(value: string) {
    return new Intl.DateTimeFormat("en-IN", { dateStyle: "long", timeZone: "UTC" }).format(new Date(value));
}