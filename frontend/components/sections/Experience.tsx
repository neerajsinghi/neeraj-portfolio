import { EXPERIENCE } from "../../lib/profile";

export default function Experience() {
  return (
    <section className="wrap" id="work">
      <div className="sec-head">
        <span className="sec-no">03</span>
        <h2>Experience</h2>
      </div>
      <div className="tl">
        {EXPERIENCE.map((j, i) => (
          <div className={"job" + (j.dim ? " dim" : "")} key={i}>
            <div className="job-top">
              <h3>{j.role}</h3>
              <span className="when">{j.when}</span>
            </div>
            <div className="co">{j.co}</div>
            <ul>
              {j.pts.map((p, k) => (
                <li key={k}>{p}</li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </section>
  );
}
