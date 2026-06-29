import Topology from "../components/Topology";
import Hero from "../components/sections/Hero";
import About from "../components/sections/About";
import Stack from "../components/sections/Stack";
import Experience from "../components/sections/Experience";
import Projects from "../components/sections/Projects";
import Contact from "../components/sections/Contact";

export default function Page() {
  return (
    <>
      <Topology />
      <div className="veil" />

      <header>
        <div className="wrap nav">
          <div className="brand">
            <span className="blink">▸</span> <b>neeraj</b>singhi<span className="blink">_</span>
          </div>
          <nav className="nav-links">
            <a href="#about">About</a>
            <a href="#stack">Stack</a>
            <a href="#work">Work</a>
            <a href="#projects">Projects</a>
            <a href="#agent">Ask AI</a>
            <span className="pill">
              <span className="dot" />
              Open to Senior / Staff roles
            </span>
          </nav>
        </div>
      </header>

      <main>
        <Hero />
        <About />
        <Stack />
        <Experience />
        <Projects />
        <Contact />
      </main>

      <footer>
        <div className="wrap">
          <span>© {new Date().getFullYear()} Neeraj Singhi · Delhi, India</span>
          <span>Go backend · RAG + tool-use (MCP pattern) · Next.js</span>
        </div>
      </footer>
    </>
  );
}
