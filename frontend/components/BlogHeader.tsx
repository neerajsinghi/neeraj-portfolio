import Link from "next/link";

export default function BlogHeader() {
    return (
        <header>
            <div className="wrap nav">
                <Link className="brand" href="/">
                    <span className="blink">▸</span> <b>neeraj</b>singhi<span className="blink">_</span>
                </Link>
                <nav className="nav-links blog-nav" aria-label="Blog navigation">
                    <Link href="/">Portfolio</Link>
                    <Link href="/blogs">Blogs</Link>
                </nav>
            </div>
        </header>
    );
}