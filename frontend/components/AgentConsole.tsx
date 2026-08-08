"use client";

import { useEffect, useRef, useState } from "react";
import { CHIPS } from "../lib/profile";
import { streamChat } from "../lib/api";
import { trackEvent } from "../lib/analytics";
import type { ChatItem } from "../types";

function boldify(line: string, keyBase: string) {
  return line.split(/(\*\*[^*]+\*\*)/g).map((p, i) =>
    p.startsWith("**") && p.endsWith("**") ? (
      <strong key={keyBase + i}>{p.slice(2, -2)}</strong>
    ) : (
      <span key={keyBase + i}>{p}</span>
    )
  );
}
function Rich({ text }: { text: string }) {
  return (
    <>
      {text.split(/\n{2,}/).map((para, i) => (
        <p key={i}>
          {para.split("\n").map((line, j) => (
            <span key={j}>
              {j > 0 && <br />}
              {boldify(line, `${i}-${j}-`)}
            </span>
          ))}
        </p>
      ))}
    </>
  );
}

const GREETING =
  "Hi — I'm **Neeraj's portfolio agent**. I answer from his real résumé, GitHub and LinkedIn using retrieval + tools, all served by a Go backend. Ask me anything, or tap a question below.";

export default function AgentConsole() {
  const [items, setItems] = useState<ChatItem[]>([{ kind: "bot", text: GREETING }]);
  const [input, setInput] = useState("");
  const [busy, setBusy] = useState(false);
  const streamRef = useRef<HTMLDivElement>(null);
  const histRef = useRef<{ role: string; content: string }[]>([]);
  const botTextRef = useRef("");
  const streamErrorRef = useRef(false);

  useEffect(() => {
    const el = streamRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [items]);

  function dropTyping(arr: ChatItem[]) {
    const i = arr.findIndex((x) => x.kind === "typing");
    if (i >= 0) arr.splice(i, 1);
    return arr;
  }
  function insertBeforeTyping(arr: ChatItem[], item: ChatItem) {
    const i = arr.findIndex((x) => x.kind === "typing");
    if (i < 0) arr.push(item);
    else arr.splice(i, 0, item);
    return arr;
  }

  function onEvent(event: string, data: unknown) {
    const eventData = data as Record<string, unknown>;
    if (event === "tool") {
      trackEvent("agent_tool_used", { tool_name: String(eventData.name ?? "unknown") });
    } else if (event === "text") {
      botTextRef.current += (botTextRef.current ? "\n\n" : "") + String(eventData.text ?? "");
    } else if (event === "error") {
      streamErrorRef.current = true;
    }

    setItems((prev) => {
      const arr = [...prev];
      if (event === "tool") {
        const input = eventData.input as Record<string, unknown> | undefined;
        const arg = input && Object.keys(input).length ? JSON.stringify(input) : "{}";
        insertBeforeTyping(arr, { kind: "trace", tool: String(eventData.name ?? ""), arg });
      } else if (event === "sources") {
        for (let i = arr.length - 1; i >= 0; i--) {
          if (arr[i].kind === "trace") {
            arr[i] = { ...(arr[i] as { kind: "trace"; tool: string; arg: string }), sources: (eventData.sources as string[]) || [] };
            break;
          }
        }
      } else if (event === "text") {
        const bi = arr.findIndex((x) => x.kind === "bot" && (x as { kind: "bot"; running?: boolean }).running);
        if (bi >= 0) arr[bi] = { kind: "bot", text: botTextRef.current, running: true };
        else insertBeforeTyping(arr, { kind: "bot", text: botTextRef.current, running: true });
      } else if (event === "error") {
        dropTyping(arr);
        insertBeforeTyping(arr, {
          kind: "bot",
          text: "**The agent hit an error.** " + (String(eventData.message ?? "unknown")) + "\n\nCheck the Go backend logs and your `ANTHROPIC_API_KEY`.",
        });
      }
      return arr;
    });
  }

  async function ask(text: string, inputMethod: "suggestion" | "typed") {
    const userText = text.trim();
    if (busy || !userText) return;
    const startedAt = performance.now();
    trackEvent("agent_question_submit", {
      input_method: inputMethod,
      question_length: userText.length,
    });
    setBusy(true);
    setInput("");
    botTextRef.current = "";
    streamErrorRef.current = false;
    setItems((p) => [...p, { kind: "user", text: userText }, { kind: "typing" }]);
    histRef.current = [...histRef.current, { role: "user", content: userText }];

    try {
      await streamChat(histRef.current, onEvent);
      setItems((p) => dropTyping([...p]));
      if (botTextRef.current.trim())
        histRef.current = [...histRef.current, { role: "assistant", content: botTextRef.current }];
      const duration = Math.round(performance.now() - startedAt);
      if (streamErrorRef.current) {
        trackEvent("agent_response_error", { input_method: inputMethod, duration_ms: duration });
      } else if (botTextRef.current.trim()) {
        trackEvent("agent_response_success", {
          input_method: inputMethod,
          response_length: botTextRef.current.length,
          duration_ms: duration,
        });
      } else {
        trackEvent("agent_response_empty", { input_method: inputMethod, duration_ms: duration });
      }
    } catch (err: any) {
      trackEvent("agent_response_error", {
        input_method: inputMethod,
        duration_ms: Math.round(performance.now() - startedAt),
      });
      setItems((p) =>
        dropTyping([...p]).concat({
          kind: "bot",
          text:
            "**The agent backend isn't reachable.** " +
            (err?.message || "network error") +
            "\n\nStart it with `go run .` in `/backend` and set `NEXT_PUBLIC_API_BASE` to its URL.",
        })
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="console" id="agent">
      <div className="console-bar">
        <span className="tdot r" />
        <span className="tdot y" />
        <span className="tdot g" />
        <span className="console-title">
          neeraj-agent · <span className="on">rag + tools online</span>
        </span>
      </div>

      <div className="stream" ref={streamRef}>
        {items.map((it, i) => {
          if (it.kind === "user") return <div key={i} className="msg user">{it.text}</div>;
          if (it.kind === "bot") return <div key={i} className="msg bot"><Rich text={it.text} /></div>;
          if (it.kind === "typing")
            return (
              <div key={i} className="typing"><i /><i /><i /></div>
            );
          return (
            <div key={i} className="trace">
              <span>
                <span className="tool">▸ {it.tool}</span>
                <span className="arg">({it.arg})</span>
              </span>
              {it.sources && (
                <span className="src">↳ retrieved: [{it.sources.length ? it.sources.join(", ") : "—"}]</span>
              )}
            </div>
          );
        })}
      </div>

      <div className="chips">
        {CHIPS.map((c) => (
          <button key={c} className="chip" onClick={() => ask(c, "suggestion")} disabled={busy}>
            {c}
          </button>
        ))}
      </div>

      <div className="composer">
        <input
          aria-label="Ask Neeraj's portfolio agent"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && ask(input, "typed")}
          placeholder="Ask about Neeraj's experience…"
          autoComplete="off"
        />
        <button className="send" onClick={() => ask(input, "typed")} disabled={busy} title="Send" aria-label="Send question">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M22 2 11 13M22 2l-7 20-4-9-9-4 20-7z" />
          </svg>
        </button>
      </div>
    </div>
  );
}
