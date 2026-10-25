"use client";

import { useEffect, useRef } from "react";

export default function Topology() {
  const ref = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const c = ref.current;
    if (!c) return;
    const x = c.getContext("2d");
    if (!x) return;

    const reduce = matchMedia("(prefers-reduced-motion: reduce)").matches;
    const N = 46;
    let w = 0,
      h = 0,
      raf = 0,
      nodes: { x: number; y: number; vx: number; vy: number }[] = [];

    function init() {
      nodes = Array.from({ length: N }, () => ({
        x: Math.random() * w,
        y: Math.random() * h,
        vx: (Math.random() - 0.5) * 0.18,
        vy: (Math.random() - 0.5) * 0.18,
      }));
    }
    function size() {
      w = c!.width = innerWidth;
      // The canvas is position: fixed, so it only ever needs to cover the viewport.
      h = c!.height = innerHeight;
      init();
    }
    function frame() {
      x!.clearRect(0, 0, w, h);
      for (let i = 0; i < N; i++) {
        const a = nodes[i];
        if (!reduce) {
          a.x += a.vx;
          a.y += a.vy;
          if (a.x < 0 || a.x > w) a.vx *= -1;
          if (a.y < 0 || a.y > h) a.vy *= -1;
        }
        for (let j = i + 1; j < N; j++) {
          const b = nodes[j];
          const dx = a.x - b.x,
            dy = a.y - b.y,
            d = Math.hypot(dx, dy);
          if (d < 140) {
            x!.strokeStyle = `rgba(79,209,197,${0.1 * (1 - d / 140)})`;
            x!.lineWidth = 1;
            x!.beginPath();
            x!.moveTo(a.x, a.y);
            x!.lineTo(b.x, b.y);
            x!.stroke();
          }
        }
        x!.fillStyle = "rgba(154,166,255,.45)";
        x!.beginPath();
        x!.arc(a.x, a.y, 1.3, 0, 7);
        x!.fill();
      }
      if (!reduce) raf = requestAnimationFrame(frame);
    }

    let started = false;
    const start = () => {
      if (started) return;
      started = true;
      size();
      frame();
    };
    // Decorative only: wait until the page is idle so it never competes with first paint.
    const idle = window.requestIdleCallback ?? ((cb: () => void) => window.setTimeout(cb, 1500));
    const idleId = idle(start);

    const onResize = () => {
      if (!started) return;
      cancelAnimationFrame(raf);
      size();
      frame();
    };
    const onVisibility = () => {
      if (!started || reduce) return;
      cancelAnimationFrame(raf);
      if (!document.hidden) raf = requestAnimationFrame(frame);
    };
    addEventListener("resize", onResize);
    document.addEventListener("visibilitychange", onVisibility);
    return () => {
      (window.cancelIdleCallback ?? window.clearTimeout)(idleId);
      cancelAnimationFrame(raf);
      removeEventListener("resize", onResize);
      document.removeEventListener("visibilitychange", onVisibility);
    };
  }, []);

  return <canvas id="topo" ref={ref} />;
}
