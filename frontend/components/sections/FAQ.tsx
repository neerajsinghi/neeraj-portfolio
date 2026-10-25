import { FAQ as ITEMS } from "../../lib/profile";

export default function FAQ() {
  return (
    <section className="wrap" id="faq">
      <div className="sec-head">
        <span className="sec-no">05</span>
        <h2>Quick answers</h2>
      </div>
      <div className="faq">
        {ITEMS.map((item) => (
          <details className="faq-item" key={item.q}>
            <summary>
              <h3>{item.q}</h3>
            </summary>
            <p>{item.a}</p>
          </details>
        ))}
      </div>
    </section>
  );
}
