"use client";

import { useEffect, useState } from "react";
import { fetchGithubRepos } from "../lib/api";
import type { Repo } from "../types";

export default function GithubStrip() {
  const [repos, setRepos] = useState<Repo[] | null>(null);
  const [user, setUser] = useState("neerajsinghi");

  useEffect(() => {
    let alive = true;
    fetchGithubRepos()
      .then((d) => {
        if (!alive) return;
        if (d.user) setUser(d.user);
        setRepos(Array.isArray(d.repos) ? d.repos : []);
      })
      .catch(() => alive && setRepos([]));
    return () => {
      alive = false;
    };
  }, []);

  return (
    <>
      <div className="gh-head">
        <span className="eyebrow">Live from GitHub</span>
        <span className="live">
          {repos === null ? (
            <>fetching…</>
          ) : repos.length ? (
            <>
              <span className="gh-dot" /> {repos.length} public repos · @{user}
            </>
          ) : (
            <span className="gh-note">@{user} — no public repos to show (or rate-limited).</span>
          )}
        </span>
      </div>
      {repos && repos.length > 0 && (
        <div className="gh-grid">
          {repos.slice(0, 6).map((r) => (
            <a key={r.name} className="repo" href={r.html_url} target="_blank" rel="noopener noreferrer">
              <div className="rn">▸ {r.name}</div>
              <div className="rd">{r.description || "—"}</div>
              <div className="rm">
                {r.language && <span className="lang">{r.language}</span>}
                <span>★ {r.stargazers_count}</span>
                <span>⑂ {r.forks_count}</span>
              </div>
            </a>
          ))}
        </div>
      )}
    </>
  );
}
