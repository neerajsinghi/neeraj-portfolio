import GithubStrip from "../GithubStrip";
import { PROJECTS } from "../../lib/profile";

export default function Projects() {
  return (
    <section className="wrap" id="projects">
      <div className="sec-head">
        <span className="sec-no">04</span>
        <h2>Projects</h2>
      </div>
      <div className="projects">
        {PROJECTS.map((p) => (
          <div className="card" key={p.name}>
            <div className="ptag">{p.tag}</div>
            <h3>{p.name}</h3>
            <p>{p.desc}</p>
            <div className="st">
              {p.stack.map((s) => (
                <span key={s}>{s}</span>
              ))}
            </div>
          </div>
        ))}
      </div>
      <GithubStrip />
    </section>
  );
}
