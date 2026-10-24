import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import BlogHeader from "../../../components/BlogHeader";
import { getBlogPost } from "../../../lib/blogs";
import { siteUrl } from "../../../lib/site";

type Props = { params: Promise<{ slug: string }> };

// Without generateStaticParams, a dynamic segment's notFound() result can get
// cached indefinitely by the Full Route Cache; these force it to respect the
// same 60s window (and on-demand revalidateTag calls) as the data fetch.
export const revalidate = 60;
export const dynamicParams = true;

export async function generateMetadata({ params }: Props): Promise<Metadata> {
    const { slug } = await params;
    const post = await getBlogPost(slug);
    if (!post) return { title: "Blog not found | Neeraj Singhi" };
    return {
        title: `${post.title} | Neeraj Singhi`,
        description: post.description,
        keywords: post.tags,
        authors: [{ name: "Neeraj Singhi", url: siteUrl }],
        alternates: { canonical: `/blogs/${post.slug}` },
        openGraph: {
            type: "article",
            url: `/blogs/${post.slug}`,
            siteName: "Neeraj Singhi",
            locale: "en_IN",
            title: post.title,
            description: post.description,
            publishedTime: post.published_at,
            modifiedTime: post.updated_at,
            authors: [siteUrl],
            tags: post.tags,
        },
        // Set explicitly so posts don't inherit the homepage's Twitter card.
        twitter: {
            card: "summary_large_image",
            title: post.title,
            description: post.description,
        },
    };
}

export default async function BlogDetailPage({ params }: Props) {
    const { slug } = await params;
    const post = await getBlogPost(slug);
    if (!post) notFound();

    const postUrl = `${siteUrl}/blogs/${post.slug}`;
    const articleData = {
        "@context": "https://schema.org",
        "@graph": [
            {
                "@type": "BlogPosting",
                "@id": `${postUrl}#article`,
                url: postUrl,
                headline: post.title,
                description: post.description,
                image: `${postUrl}/opengraph-image`,
                datePublished: post.published_at || post.updated_at,
                dateModified: post.updated_at,
                inLanguage: "en",
                keywords: post.tags.join(", "),
                mainEntityOfPage: postUrl,
                isPartOf: { "@id": `${siteUrl}/#website` },
                author: { "@type": "Person", "@id": `${siteUrl}/#person`, name: "Neeraj Singhi", url: siteUrl },
                publisher: { "@type": "Person", "@id": `${siteUrl}/#person`, name: "Neeraj Singhi", url: siteUrl },
            },
            {
                "@type": "BreadcrumbList",
                itemListElement: [
                    { "@type": "ListItem", position: 1, name: "Home", item: siteUrl },
                    { "@type": "ListItem", position: 2, name: "Blogs", item: `${siteUrl}/blogs` },
                    { "@type": "ListItem", position: 3, name: post.title, item: postUrl },
                ],
            },
        ],
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