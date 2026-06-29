import { LINKS } from "../../lib/profile";

export default function Contact() {
  return (
    <section className="wrap" id="contact">
      <div className="sec-head">
        <span className="sec-no">05</span>
        <h2>Get in touch</h2>
      </div>
      <div className="contact-grid">
        <a className="clink" href={`mailto:${LINKS.email}`}>
          <span className="ic">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <rect x="3" y="5" width="18" height="14" rx="2" />
              <path d="m3 7 9 6 9-6" />
            </svg>
          </span>
          <span>
            <span className="k">email</span>
            <br />
            <span className="v">{LINKS.email}</span>
          </span>
        </a>
        <a className="clink" href={LINKS.linkedin} target="_blank" rel="noopener noreferrer">
          <span className="ic">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M4.98 3.5A2.5 2.5 0 1 0 5 8.5a2.5 2.5 0 0 0 0-5zM3 9h4v12H3zM9 9h3.8v1.7h.05c.53-1 1.83-2.06 3.77-2.06 4.03 0 4.78 2.65 4.78 6.1V21H21v-5.6c0-1.34-.02-3.06-1.86-3.06-1.87 0-2.16 1.46-2.16 2.96V21H13z" />
            </svg>
          </span>
          <span>
            <span className="k">linkedin</span>
            <br />
            <span className="v">in/neeraj-singhi-golang</span>
          </span>
        </a>
        <a className="clink" href={LINKS.github} target="_blank" rel="noopener noreferrer">
          <span className="ic">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2a10 10 0 0 0-3.16 19.49c.5.09.68-.22.68-.48l-.01-1.7c-2.78.6-3.37-1.34-3.37-1.34-.45-1.16-1.11-1.47-1.11-1.47-.91-.62.07-.6.07-.6 1 .07 1.53 1.03 1.53 1.03.9 1.53 2.36 1.09 2.93.83.09-.65.35-1.1.63-1.35-2.22-.25-4.55-1.11-4.55-4.94 0-1.09.39-1.98 1.03-2.68-.1-.25-.45-1.27.1-2.64 0 0 .84-.27 2.75 1.02a9.6 9.6 0 0 1 5 0c1.91-1.29 2.75-1.02 2.75-1.02.55 1.37.2 2.39.1 2.64.64.7 1.03 1.59 1.03 2.68 0 3.84-2.34 4.69-4.57 4.94.36.31.68.92.68 1.85l-.01 2.74c0 .27.18.58.69.48A10 10 0 0 0 12 2z" />
            </svg>
          </span>
          <span>
            <span className="k">github</span>
            <br />
            <span className="v">github.com/neerajsinghi</span>
          </span>
        </a>
        <a className="clink" href="/Neeraj_Singhi_Resume.pdf" download>
          <span className="ic">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <path d="M14 2v6h6M9 15h6M9 11h2" />
            </svg>
          </span>
          <span>
            <span className="k">résumé</span>
            <br />
            <span className="v">Download PDF</span>
          </span>
        </a>
      </div>
    </section>
  );
}
