import type { Metadata } from "next";
import { Inter, Space_Grotesk, JetBrains_Mono } from "next/font/google";
import Script from "next/script";
import "./globals.css";

const sans = Inter({ subsets: ["latin"], variable: "--font-sans" });
const disp = Space_Grotesk({ subsets: ["latin"], weight: ["500", "600", "700"], variable: "--font-disp" });
const mono = JetBrains_Mono({ subsets: ["latin"], weight: ["400", "500"], variable: "--font-mono" });
const siteUrl = "https://neerajsinghi.com";
const googleAnalyticsId = "G-6N1XEY3L6R";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: "Senior Go Backend & AI Engineer | Neeraj Singhi",
  description:
    "Neeraj Singhi is a senior Go backend and AI engineer with 10+ years building distributed systems, AWS platforms, RAG applications, and secure services.",
  alternates: {
    canonical: "/",
  },
  openGraph: {
    type: "profile",
    url: siteUrl,
    title: "Senior Go Backend & AI Engineer | Neeraj Singhi",
    description:
      "10+ years building distributed systems, AWS platforms, secure Go services, and AI/RAG applications.",
    siteName: "Neeraj Singhi",
    locale: "en_IN",
  },
  twitter: {
    card: "summary_large_image",
    title: "Senior Go Backend & AI Engineer | Neeraj Singhi",
    description:
      "10+ years building distributed systems, AWS platforms, secure Go services, and AI/RAG applications.",
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

const structuredData = {
  "@context": "https://schema.org",
  "@type": "ProfilePage",
  "@id": `${siteUrl}/#profile`,
  url: siteUrl,
  name: "Neeraj Singhi - Senior Go Backend & AI Engineer",
  mainEntity: {
    "@type": "Person",
    "@id": `${siteUrl}/#person`,
    name: "Neeraj Singhi",
    url: siteUrl,
    jobTitle: "Senior Go Backend & AI Engineer",
    email: "mailto:nsinghi2011@gmail.com",
    homeLocation: {
      "@type": "Place",
      name: "Delhi, India",
    },
    sameAs: [
      "https://github.com/neerajsinghi",
      "https://www.linkedin.com/in/neeraj-singhi-golang",
    ],
    knowsAbout: [
      "Go",
      "Distributed systems",
      "Backend engineering",
      "Amazon Web Services",
      "Retrieval-augmented generation",
      "AI agents",
      "Kubernetes",
      "Terraform",
    ],
  },
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${sans.variable} ${disp.variable} ${mono.variable}`}>
      <body>
        {children}
        <Script
          src={`https://www.googletagmanager.com/gtag/js?id=${googleAnalyticsId}`}
          strategy="afterInteractive"
        />
        <Script id="google-analytics" strategy="afterInteractive">
          {`
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', '${googleAnalyticsId}');
          `}
        </Script>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData).replace(/</g, "\\u003c") }}
        />
      </body>
    </html>
  );
}
