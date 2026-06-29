import { FACTS } from "../../lib/profile";

export default function About() {
  return (
    <section className="wrap" id="about">
      <div className="sec-head">
        <span className="sec-no">01</span>
        <h2>About</h2>
      </div>
      <div className="about-grid">
        <div className="about">
          <p>
            I'm a senior backend engineer with over a decade building production systems that hold up under real
            load — distributed services, SDKs, and the infrastructure around them.
          </p>
          <p>
            My work tends to sit where backend depth meets reliability and security: at Dell I built core parts of
            an internal license-enforcement service, including usage metering and tamper detection backed by
            encryption and cryptographic hashing, and ported a performance-critical C component to a single pure-Go
            artifact. More recently I've gone deep on AI integration — OpenAI, retrieval-augmented generation, and
            tool-based (MCP-style) agent workflows.
          </p>
          <p>
            Today I lead a small team delivering an AI-enabled consumer marketplace, owning architecture and
            technical direction end to end. I care about the team as much as the code: setting standards, mentoring,
            and turning vague requirements into shipped systems.
          </p>
        </div>
        <div className="facts">
          {FACTS.map(([k, v]) => (
            <div className="fact" key={k}>
              <span className="k">{k}</span>
              <span className="v">{v}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
