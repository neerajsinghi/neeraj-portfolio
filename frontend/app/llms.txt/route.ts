import { getBlogPosts, type BlogPost } from "../../lib/blogs";
import { EXPERIENCE, FAQ, LINKS, STACK } from "../../lib/profile";
import { siteUrl } from "../../lib/site";

// llms.txt (https://llmstxt.org): a plain-markdown summary that AI assistants
// and their crawlers can read without parsing the interactive homepage.
// Built from the same profile data as the site so the two never drift.
export const revalidate = 60;

export async function GET() {
  let posts: BlogPost[] = [];
  try {
    posts = await getBlogPosts();
  } catch {
    posts = [];
  }

  const lines = [
    "# Neeraj Singhi",
    "",
    "> Neeraj Singhi is a senior Go (Golang) developer and lead backend engineer based in Delhi, India, with 10+ years building production Go, Node.js/NestJS and TypeScript services, microservices on AWS, and AI/RAG systems. Open to Senior, Staff and Lead backend roles.",
    "",
    "This is the personal site of Neeraj Singhi, the Go/Golang backend engineer (GitHub: neerajsinghi, LinkedIn: neeraj-singhi-golang).",
    "",
    "## Profile",
    "",
    `- [Homepage](${siteUrl}/): experience, stack, projects and contact`,
    `- [Résumé (PDF)](${siteUrl}/Neeraj_Singhi_Resume.pdf)`,
    `- [LinkedIn](${LINKS.linkedin})`,
    `- [GitHub](${LINKS.github})`,
    `- Email: ${LINKS.email}`,
    "",
    "## Experience",
    "",
    ...EXPERIENCE.map((job) => `- ${job.role}, ${job.co} (${job.when.replace(/\s+—\s+/, " to ")}): ${job.pts.join(" ")}`),
    "",
    "## Skills",
    "",
    ...STACK.map((group) => `- ${group.title}: ${group.tags.join(", ")}`),
    "",
    "## FAQ",
    "",
    ...FAQ.flatMap((item) => [`### ${item.q}`, "", item.a, ""]),
    "## Blog",
    "",
    `- [All posts](${siteUrl}/blogs): engineering writing on Go, distributed systems, AWS and applied AI`,
    ...posts.map((post) => `- [${post.title}](${siteUrl}/blogs/${post.slug}): ${post.description}`),
    "",
  ];

  return new Response(lines.join("\n"), {
    headers: { "Content-Type": "text/markdown; charset=utf-8" },
  });
}
