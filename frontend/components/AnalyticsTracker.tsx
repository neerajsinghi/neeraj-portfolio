"use client";

import { useEffect } from "react";
import { trackEvent } from "../lib/analytics";

function trackLink(link: HTMLAnchorElement) {
    const rawHref = link.getAttribute("href") || "";
    const location = link.closest("header") ? "header" : link.closest("#top") ? "hero" : link.closest("#contact") ? "contact" : link.closest("#projects") ? "projects" : "page";

    if (rawHref.startsWith("#")) {
        trackEvent("navigation_click", { link_target: rawHref.slice(1), link_location: location });
        return;
    }

    if (rawHref.startsWith("mailto:")) {
        trackEvent("contact_click", { contact_method: "email", link_location: location });
        return;
    }

    if (link.hasAttribute("download") || rawHref.toLowerCase().endsWith(".pdf")) {
        trackEvent("resume_download", { file_name: rawHref.split("/").pop() || "resume.pdf", link_location: location });
        return;
    }

    const url = new URL(link.href, window.location.href);
    if (url.origin !== window.location.origin) {
        const destination = url.hostname.includes("linkedin.com")
            ? "linkedin"
            : url.hostname.includes("github.com")
                ? "github"
                : url.hostname;
        trackEvent("outbound_link_click", {
            link_destination: destination,
            link_location: location,
            link_path: destination === "github" ? url.pathname : "",
        });
    }
}

export default function AnalyticsTracker() {
    useEffect(() => {
        function onClick(event: MouseEvent) {
            const target = event.target;
            if (!(target instanceof Element)) return;
            const link = target.closest("a");
            if (link instanceof HTMLAnchorElement) trackLink(link);
        }

        document.addEventListener("click", onClick);
        return () => document.removeEventListener("click", onClick);
    }, []);

    useEffect(() => {
        const viewed = new Set<string>();
        const sections = document.querySelectorAll<HTMLElement>("main section[id]");

        function recordVisibleSections() {
            for (const section of sections) {
                if (viewed.has(section.id)) continue;
                const rect = section.getBoundingClientRect();
                const visibleHeight = Math.min(rect.bottom, window.innerHeight) - Math.max(rect.top, 0);
                if (visibleHeight < Math.min(rect.height, window.innerHeight) * 0.25) continue;
                viewed.add(section.id);
                trackEvent("section_view", { section_name: section.id });
            }
        }

        recordVisibleSections();
        window.addEventListener("scroll", recordVisibleSections, { passive: true });
        return () => window.removeEventListener("scroll", recordVisibleSections);
    }, []);

    useEffect(() => {
        const milestones = [25, 50, 75, 90];
        const reached = new Set<number>();

        function onScroll() {
            const scrollable = document.documentElement.scrollHeight - window.innerHeight;
            if (scrollable <= 0) return;
            const depth = Math.round((window.scrollY / scrollable) * 100);
            for (const milestone of milestones) {
                if (depth < milestone || reached.has(milestone)) continue;
                reached.add(milestone);
                trackEvent("scroll_depth", { percent_scrolled: milestone });
            }
        }

        window.addEventListener("scroll", onScroll, { passive: true });
        return () => window.removeEventListener("scroll", onScroll);
    }, []);

    return null;
}