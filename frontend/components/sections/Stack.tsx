import { STACK } from "../../lib/profile";

export default function Stack() {
  return (
    <section className="wrap" id="stack">
      <div className="sec-head">
        <span className="sec-no">02</span>
        <h2>The stack</h2>
      </div>
      <div className="stack">
        {STACK.map((g) => (
          <div className="grp" key={g.title}>
            <h3>{g.title}</h3>
            <div className="tags">
              {g.tags.map((t) => (
                <span className="tag" key={t}>
                  {t}
                </span>
              ))}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
