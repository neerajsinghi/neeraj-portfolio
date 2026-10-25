import type { Metadata } from "next";
import Topology from "../components/Topology";
import Hero from "../components/sections/Hero";
import About from "../components/sections/About";
import Stack from "../components/sections/Stack";
import Experience from "../components/sections/Experience";
import Projects from "../components/sections/Projects";
import FAQ from "../components/sections/FAQ";
import Contact from "../components/sections/Contact";
import { FAQ as FAQ_ITEMS } from "../lib/profile";
import { siteUrl } from "../lib/site";

// Google ignores the keywords meta tag, but Bing and some job/people search
// tools still read it. Ranking itself comes from the visible copy and schema.
export const metadata: Metadata = {
  keywords: [
    "Neeraj Singhi",
    "Neeraj",
    "neerajsinghi",
    "Go developer",
    "Golang developer",
    "Golang engineer",
    "Senior Golang developer",
    "Go backend engineer",
    "Backend developer",
    "Node.js developer",
    "NestJS developer",
    "TypeScript developer",
    "Senior software engineer",
    "Staff software engineer",
    "Lead backend engineer",
    "Tech lead",
    "Microservices engineer",
    "AWS backend engineer",
    "AI engineer",
    "RAG engineer",
    "Freelance Golang developer",
    "Golang developer India",
    "Golang developer Delhi",
  ],
};

const profile = {
  "@context": "https://schema.org",
  "@type": "ProfilePage",
  "@id": `${siteUrl}/#profile`,
  url: siteUrl,
  name: "Neeraj Singhi - Senior Go (Golang) Developer & Backend Lead",
  isPartOf: { "@id": `${siteUrl}/#website` },
  mainEntity: {
    "@type": "Person",
    "@id": `${siteUrl}/#person`,
    name: "Neeraj Singhi",
    givenName: "Neeraj",
    familyName: "Singhi",
    alternateName: ["Neeraj", "neerajsinghi"],
    url: siteUrl,
    image: `${siteUrl}/opengraph-image`,
    jobTitle: ["Lead Backend Engineer", "Senior Software Engineer", "Senior Go (Golang) Developer"],
    description:
      "Senior Go (Golang) developer and lead backend engineer with 10+ years of Go, Node.js and TypeScript, building microservices, AWS platforms and AI/RAG systems.",
    email: "mailto:nsinghi2011@gmail.com",
    homeLocation: {
      "@type": "Place",
      name: "Delhi, India",
    },
    nationality: { "@type": "Country", name: "India" },
    sameAs: [
      "https://github.com/neerajsinghi",
      "https://www.linkedin.com/in/neeraj-singhi-golang",
      "https://tal.club/meet/neeraj-singhi",
    ],
    hasOccupation: {
      "@type": "Occupation",
      name: "Senior Backend Software Engineer",
      alternateName: ["Golang Developer", "Go Backend Engineer", "Node.js Developer", "Staff Software Engineer", "Tech Lead"],
      occupationLocation: { "@type": "Country", name: "India" },
      skills: "Go, Golang, Node.js, NestJS, TypeScript, JavaScript, gRPC, REST, microservices, AWS, Docker, Kubernetes, Terraform, MongoDB, MySQL, Redis, OpenAI, RAG, AI agents",
    },
    knowsAbout: [
      "Go",
      "Golang",
      "Node.js",
      "NestJS",
      "TypeScript",
      "Distributed systems",
      "Microservices",
      "gRPC",
      "Backend engineering",
      "Amazon Web Services",
      "Docker",
      "Kubernetes",
      "Terraform",
      "MongoDB",
      "Redis",
      "Retrieval-augmented generation",
      "AI agents",
      "Technical leadership",
    ],
  },
};

const faqData = {
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "@id": `${siteUrl}/#faq`,
  about: { "@id": `${siteUrl}/#person` },
  mainEntity: FAQ_ITEMS.map((item) => ({
    "@type": "Question",
    name: item.q,
    acceptedAnswer: { "@type": "Answer", text: item.a },
  })),
};

export default function Page() {
  return (
    <>
      <Topology />
      <div className="veil" />

      <header>
        <div className="wrap nav">
          <div className="brand">
            <span className="blink">▸</span> <b>neeraj</b>singhi<span className="blink">_</span>
          </div>
          <nav className="nav-links">
            <a href="#about">About</a>
            <a href="#stack">Stack</a>
            <a href="#work">Work</a>
            <a href="#projects">Projects</a>
            <a href="/blogs">Blogs</a>
            <a href="#agent">Ask AI</a>
            <span className="pill">
              <span className="dot" />
              Open to Senior / Staff / Lead roles
            </span>
          </nav>
        </div>
      </header>

      <main>
        <Hero />
        <About />
        <Stack />
        <Experience />
        <Projects />
        <FAQ />
        <Contact />
      </main>

      <footer>
        <div className="wrap">
          <span>© {new Date().getFullYear()} Neeraj Singhi · Delhi, India</span>
          <span>Golang developer · Node.js · TypeScript · RAG + tool-use (MCP pattern) · Next.js</span>
        </div>
      </footer>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(profile).replace(/</g, "\\u003c") }}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(faqData).replace(/</g, "\\u003c") }}
      />
    </>
  );
}
