import { ImageResponse } from "next/og";
import { getBlogPost } from "../../../lib/blogs";

export const alt = "Engineering article by Neeraj Singhi";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

type Props = { params: Promise<{ slug: string }> };

export default async function BlogOpenGraphImage({ params }: Props) {
    const { slug } = await params;
    let title = "Engineering notes";
    let tags: string[] = [];
    try {
        const post = await getBlogPost(slug);
        if (post) {
            title = post.title;
            tags = post.tags.slice(0, 4);
        }
    } catch {
        // Fall back to the generic card if the blog API is unavailable.
    }

    return new ImageResponse(
        (
            <div
                style={{
                    background: "#08101d",
                    color: "#f2f6ff",
                    display: "flex",
                    flexDirection: "column",
                    height: "100%",
                    justifyContent: "space-between",
                    padding: "72px 84px",
                    width: "100%",
                }}
            >
                <div style={{ color: "#45e0c1", display: "flex", fontSize: 25, letterSpacing: 0 }}>
                    NEERAJSINGHI.COM / BLOGS
                </div>
                <div style={{ display: "flex", fontSize: title.length > 90 ? 50 : 62, fontWeight: 700, letterSpacing: 0, lineHeight: 1.15 }}>
                    {title}
                </div>
                <div style={{ borderTop: "2px solid #25334a", color: "#f4bf5f", display: "flex", fontSize: 25, paddingTop: 24, width: "100%" }}>
                    {["Neeraj Singhi", ...tags].join(" · ")}
                </div>
            </div>
        ),
        size
    );
}
