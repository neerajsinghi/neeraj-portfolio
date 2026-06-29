import type { Metadata } from "next";
import { Inter, Space_Grotesk, JetBrains_Mono } from "next/font/google";
import "./globals.css";

const sans = Inter({ subsets: ["latin"], variable: "--font-sans" });
const disp = Space_Grotesk({ subsets: ["latin"], weight: ["500", "600", "700"], variable: "--font-disp" });
const mono = JetBrains_Mono({ subsets: ["latin"], weight: ["400", "500"], variable: "--font-mono" });

export const metadata: Metadata = {
  title: "Neeraj Singhi — Backend & AI Engineer",
  description:
    "Neeraj Singhi — senior backend engineer (Go, distributed systems, AWS, AI/RAG). Ask the AI agent anything about his work.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${sans.variable} ${disp.variable} ${mono.variable}`}>
      <body>{children}</body>
    </html>
  );
}
