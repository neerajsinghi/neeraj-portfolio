import AgentConsole from "../AgentConsole";
import { LINKS } from "../../lib/profile";

export default function Hero() {
  return (
    <section className="wrap hero" id="top">
      <div>
        <div className="eyebrow">Backend · Distributed Systems · AI</div>
        <h1>
          Neeraj Singhi
          <br />
          <span className="grad">
            builds the systems
            <br />
            behind the product.
          </span>
        </h1>
        <div className="role">
          Senior backend engineer · <b>10+ years</b> shipping production <b>Go</b> services,
          <br />
          SDKs &amp; containerized systems on <b>AWS</b> — now leading <b>AI / RAG</b> work.
        </div>
        <p className="lede">
          From license-enforcement and tamper detection at Dell to an AI-enabled marketplace today. Ask the agent
          anything about my work — it runs on a Go backend that retrieves from my real résumé and GitHub.
        </p>
        <div className="cta">
          <a className="btn primary" href="#agent">
            Ask my AI agent
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M5 12h14M13 6l6 6-6 6" />
            </svg>
          </a>
          <a className="btn" href="/Neeraj_Singhi_Resume.pdf" download>
            Résumé
          </a>
          <a className="btn" href={LINKS.github} target="_blank" rel="noopener noreferrer">
            GitHub
          </a>
          <a className="btn" href={LINKS.linkedin} target="_blank" rel="noopener noreferrer">
            LinkedIn
          </a>
        </div>
      </div>

      <AgentConsole />
    </section>
  );
}
