import type { Metadata } from "next";
import { Inter, Space_Grotesk, JetBrains_Mono } from "next/font/google";
import Script from "next/script";
import AnalyticsTracker from "../components/AnalyticsTracker";
import { siteUrl } from "../lib/site";
import "./globals.css";

const sans = Inter({ subsets: ["latin"], variable: "--font-sans" });
const disp = Space_Grotesk({ subsets: ["latin"], weight: ["500", "600", "700"], variable: "--font-disp" });
const mono = JetBrains_Mono({ subsets: ["latin"], weight: ["400", "500"], variable: "--font-mono" });
const googleAnalyticsId = "G-6N1XEY3L6R";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: "Neeraj Singhi | Senior Go (Golang) Developer & Backend Lead",
  description:
    "Neeraj Singhi: senior Golang developer and backend lead. 10+ years of Go, Node.js and TypeScript microservices on AWS, plus AI/RAG. Open to Staff and Lead roles.",
  alternates: {
    canonical: "/",
  },
  openGraph: {
    type: "profile",
    url: siteUrl,
    title: "Neeraj Singhi | Senior Go (Golang) Developer & Backend Lead",
    description:
      "Senior Go (Golang) developer and backend lead: 10+ years of Go, Node.js and TypeScript microservices, AWS, and AI/RAG systems.",
    siteName: "Neeraj Singhi",
    locale: "en_IN",
  },
  twitter: {
    card: "summary_large_image",
    title: "Neeraj Singhi | Senior Go (Golang) Developer & Backend Lead",
    description:
      "Senior Go (Golang) developer and backend lead: 10+ years of Go, Node.js and TypeScript microservices, AWS, and AI/RAG systems.",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-image-preview": "large",
      "max-snippet": -1,
      "max-video-preview": -1,
    },
  },
};

const website = {
  "@context": "https://schema.org",
  "@type": "WebSite",
  "@id": `${siteUrl}/#website`,
  url: siteUrl,
  name: "Neeraj Singhi",
  alternateName: ["neerajsinghi", "neerajsinghi.com", "Neeraj Singhi Portfolio"],
  inLanguage: "en",
  publisher: { "@type": "Person", "@id": `${siteUrl}/#person`, name: "Neeraj Singhi", url: siteUrl },
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${sans.variable} ${disp.variable} ${mono.variable}`}>
      <body>
        {children}
        <AnalyticsTracker />
        <Script
          src={`https://www.googletagmanager.com/gtag/js?id=${googleAnalyticsId}`}
          strategy="lazyOnload"
        />
        <Script id="google-analytics" strategy="lazyOnload">
          {`
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', '${googleAnalyticsId}');
          `}
        </Script>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(website).replace(/</g, "\\u003c") }}
        />
      </body>
    </html>
  );
}
