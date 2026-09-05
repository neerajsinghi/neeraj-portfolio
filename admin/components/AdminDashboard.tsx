"use client";

import { useEffect, useState } from "react";
import { Archive, BookOpen, CalendarClock, Check, Download, ExternalLink, FilePlus2, LogIn, LogOut, PanelLeft, Save, Send, Trash2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { createPost, deletePost, listPosts, publishExternally, updatePost, type BlogInput, type BlogPost, type ExternalPlatform, type Status } from "../lib/api";
import { beginLogin, getRoles, getSession, logout } from "../lib/auth";

const emptyPost: BlogInput = { slug: "", title: "", description: "", content_markdown: "", linkedin_post: "", social_post: "", tags: [], status: "draft", version: 0, scheduled_at: "", publish_devto: false, publish_linkedin: false };

export default function AdminDashboard() {
    const [token, setToken] = useState("");
    const [roles, setRoles] = useState<string[]>([]);
    const [posts, setPosts] = useState<BlogPost[]>([]);
    const [selectedID, setSelectedID] = useState<string | null>(null);
    const [draft, setDraft] = useState<BlogInput>(emptyPost);
    const [mode, setMode] = useState<"edit" | "preview">("edit");
    const [message, setMessage] = useState("");
    const [busy, setBusy] = useState(false);
    const [ready, setReady] = useState(false);
    const [destinations, setDestinations] = useState({ website: true, devto: false, linkedin: false });
    const [statusFilter, setStatusFilter] = useState<Status | "all">("all");
    const [sortBy, setSortBy] = useState<"updated_desc" | "updated_asc" | "title_asc" | "title_desc">("updated_desc");
    const isAdmin = roles.includes("admin");
    const selected = posts.find((post) => post.id === selectedID);
    const locked = Boolean(selected && selected.status !== "draft" && !isAdmin);
    const visiblePosts = posts
        .filter((post) => statusFilter === "all" || post.status === statusFilter)
        .sort((a, b) => {
            switch (sortBy) {
                case "updated_asc":
                    return new Date(a.updated_at).getTime() - new Date(b.updated_at).getTime();
                case "title_asc":
                    return a.title.localeCompare(b.title);
                case "title_desc":
                    return b.title.localeCompare(a.title);
                default:
                    return new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime();
            }
        });

    useEffect(() => {
        const session = getSession();
        if (!session) {
            setReady(true);
            return;
        }
        setToken(session.accessToken);
        setRoles(getRoles(session.accessToken));
        listPosts(session.accessToken)
            .then(setPosts)
            .catch(handleError)
            .finally(() => setReady(true));
    }, []);

    function handleError(error: unknown) {
        const text = error instanceof Error ? error.message : "Request failed";
        if (text === "SESSION_EXPIRED") {
            setToken("");
            setMessage("Your session expired. Sign in again.");
        } else {
            setMessage(text);
        }
    }

    function selectPost(post: BlogPost) {
        setSelectedID(post.id);
        setDraft({
            slug: post.slug, title: post.title, description: post.description,
            content_markdown: post.content_markdown, linkedin_post: post.linkedin_post || "", social_post: post.social_post || "",
            tags: post.tags, status: post.status, version: post.version,
            scheduled_at: post.scheduled_at || "",
            publish_devto: post.publish_devto || false, publish_linkedin: post.publish_linkedin || false,
        });
        // Approving a post should re-publish externally wherever it was intended but hasn't gone out yet.
        setDestinations({ website: true, devto: Boolean(post.publish_devto) && !post.devto_url, linkedin: Boolean(post.publish_linkedin) && !post.linkedin_url });
        setMessage("");
    }

    function newPost() {
        setSelectedID(null);
        setDraft(emptyPost);
        setDestinations({ website: true, devto: false, linkedin: false });
        setMode("edit");
        setMessage("");
    }

    async function save(status: Status = "draft"): Promise<BlogPost | undefined> {
        if (!token) return undefined;
        if (status === "scheduled" && !draft.scheduled_at) {
            setMessage("Choose a future publication date and time before scheduling.");
            return;
        }
        setBusy(true);
        setMessage("");
        try {
            const input = { ...draft, status, scheduled_at: status === "scheduled" ? draft.scheduled_at : undefined };
            const saved = selectedID ? await updatePost(token, selectedID, input) : await createPost(token, input);
            setPosts((current) => [saved, ...current.filter((post) => post.id !== saved.id)]);
            selectPost(saved);
            setMessage(status === "published" ? "Published successfully." : status === "scheduled" ? `Scheduled for ${formatDateTime(saved.scheduled_at)}.` : status === "archived" ? "Archived." : "Draft saved.");
            return saved;
        } catch (error) {
            handleError(error);
        } finally {
            setBusy(false);
        }
        return undefined;
    }

    async function publishSelected() {
        if (!token || !isAdmin) return;
        const external = (["devto", "linkedin"] as ExternalPlatform[]).filter((platform) => destinations[platform]);
        if (!destinations.website && external.length === 0) {
            setMessage("Select at least one publishing destination.");
            return;
        }
        if (external.length > 0 && !destinations.website && selected?.status !== "published") {
            setMessage("Publish on the personal site first so external posts have a canonical URL.");
            return;
        }
        const labels = [destinations.website ? "neerajsinghi.com" : "", ...external.map((platform) => platform === "devto" ? "DEV" : "LinkedIn")].filter(Boolean);
        if (!window.confirm(`Publish this article to ${labels.join(", ")}? External publishing cannot be undone from this console.`)) return;
        setBusy(true);
        setMessage("");
        try {
            let publishInput = draft;
            if (destinations.website) {
                const saved = await save("published");
                if (!saved) return;
                publishInput = { ...draft, version: saved.version, status: saved.status };
            }
            const results = external.length > 0 ? await publishExternally(token, publishInput, external) : [];
            const publishedLabels = [destinations.website ? "neerajsinghi.com" : "", ...results.map((result) => result.platform === "devto" ? "DEV" : "LinkedIn")].filter(Boolean);
            setMessage(`Published to ${publishedLabels.join(", ")}.`);
        } catch (error) {
            handleError(error);
        } finally {
            setBusy(false);
        }
    }

    function exportForMedium() {
        const frontmatter = `---\ntitle: "${draft.title.replaceAll('"', '\\"')}"\ndescription: "${draft.description.replaceAll('"', '\\"')}"\ncanonical_url: "https://neerajsinghi.com/blogs/${draft.slug}"\ntags: [${draft.tags.map((tag) => `"${tag}"`).join(", ")}]\n---\n\n`;
        const link = document.createElement("a");
        link.href = URL.createObjectURL(new Blob([frontmatter + draft.content_markdown], { type: "text/markdown" }));
        link.download = `${draft.slug || "article"}.md`;
        link.click();
        URL.revokeObjectURL(link.href);
    }

    async function remove() {
        if (!token || !selectedID || !isAdmin || !window.confirm("Permanently delete this post? Revisions remain in the audit collection.")) return;
        setBusy(true);
        try {
            await deletePost(token, selectedID);
            setPosts((current) => current.filter((post) => post.id !== selectedID));
            newPost();
            setMessage("Post deleted.");
        } catch (error) {
            handleError(error);
        } finally {
            setBusy(false);
        }
    }

    if (!ready) return <main className="auth-state"><div className="auth-mark">NS</div><p>Loading editorial workspace…</p></main>;
    if (!token) return (
        <main className="login-page">
            <section className="login-panel">
                <div className="auth-mark">NS</div>
                <span className="kicker">Private workspace</span>
                <h1>Editorial console</h1>
                <p>Draft, review, and publish engineering articles for neerajsinghi.com.</p>
                {message && <div className="notice error">{message}</div>}
                <button className="command primary" onClick={() => void beginLogin()}><LogIn size={17} /> Sign in with Cognito</button>
            </section>
        </main>
    );

    return (
        <div className="admin-shell">
            <header className="admin-header">
                <div className="admin-brand"><span>NS</span><div><b>Editorial console</b><small>neerajsinghi.com</small></div></div>
                <div className="admin-session"><span className="role-badge">{isAdmin ? "Admin" : "Editor"}</span><button className="icon-button" onClick={logout} title="Sign out"><LogOut size={18} /></button></div>
            </header>

            <aside className="post-rail">
                <div className="rail-head"><span><PanelLeft size={16} /> Posts</span><button className="icon-button" onClick={newPost} title="New post"><FilePlus2 size={18} /></button></div>
                <div className="rail-filters">
                    <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as Status | "all")} aria-label="Filter by status">
                        <option value="all">All statuses</option>
                        <option value="draft">Draft</option>
                        <option value="scheduled">Scheduled</option>
                        <option value="published">Published</option>
                        <option value="archived">Archived</option>
                    </select>
                    <select value={sortBy} onChange={(event) => setSortBy(event.target.value as typeof sortBy)} aria-label="Sort posts">
                        <option value="updated_desc">Newest updated</option>
                        <option value="updated_asc">Oldest updated</option>
                        <option value="title_asc">Title A–Z</option>
                        <option value="title_desc">Title Z–A</option>
                    </select>
                </div>
                <div className="post-list">
                    {visiblePosts.map((post) => (
                        <button className={`post-row ${selectedID === post.id ? "active" : ""}`} key={post.id} onClick={() => selectPost(post)}>
                            <span className={`status-dot ${post.status}`} />
                            <span><b>{post.title}</b><small>{post.status === "scheduled" ? `scheduled ${formatDateTime(post.scheduled_at)}` : post.status} · v{post.version}</small></span>
                        </button>
                    ))}
                    {posts.length === 0 && <p className="rail-empty">No posts yet.</p>}
                    {posts.length > 0 && visiblePosts.length === 0 && <p className="rail-empty">No posts match this filter.</p>}
                </div>
            </aside>

            <main className="editor-workspace">
                <div className="editor-toolbar">
                    <div className="mode-switch" role="tablist" aria-label="Editor mode">
                        <button className={mode === "edit" ? "active" : ""} onClick={() => setMode("edit")}>Edit</button>
                        <button className={mode === "preview" ? "active" : ""} onClick={() => setMode("preview")}>Read &amp; review</button>
                    </div>
                    <div className="editor-actions">
                        {selected && isAdmin && <button className="icon-button danger" onClick={() => void remove()} title="Delete post" disabled={busy}><Trash2 size={17} /></button>}
                        {selected && (selected.status === "published" || (selected.status === "scheduled" && selected.scheduled_at && new Date(selected.scheduled_at) <= new Date())) && <a className="icon-button" href={`https://neerajsinghi.com/blogs/${selected.slug}`} target="_blank" rel="noreferrer" title="Read published article"><ExternalLink size={17} /></a>}
                        {selected?.devto_url && <a className="icon-button" href={selected.devto_url} target="_blank" rel="noreferrer" title="View on DEV">DEV</a>}
                        {selected?.linkedin_url && <a className="icon-button" href={selected.linkedin_url} target="_blank" rel="noreferrer" title="View on LinkedIn">in</a>}
                        {selected && isAdmin && selected.status !== "archived" && <button className="command" onClick={() => void save("archived")} disabled={busy}><Archive size={16} /> Archive</button>}
                        <button className="command" onClick={() => void save("draft")} disabled={busy || locked}><Save size={16} /> Save draft</button>
                        {isAdmin && <button className="command schedule" onClick={() => void save("scheduled")} disabled={busy || !draft.scheduled_at}><CalendarClock size={16} /> Schedule</button>}
                        {isAdmin && <button className="command primary" onClick={() => void publishSelected()} disabled={busy}><Send size={16} /> Publish selected</button>}
                    </div>
                </div>

                {message && <div className={`notice ${message.includes("success") || message.includes("saved") ? "success" : ""}`}><Check size={15} /> {message}</div>}
                {locked && <div className="notice">Published posts are read-only for editors. An admin can update or unpublish them.</div>}

                {mode === "edit" ? (
                    <div className="editor-form">
                        <label>Title<input value={draft.title} onChange={(event) => setDraft({ ...draft, title: event.target.value })} disabled={locked} maxLength={120} /></label>
                        <div className="field-row">
                            <label>Slug<input value={draft.slug} onChange={(event) => setDraft({ ...draft, slug: event.target.value })} disabled={locked} placeholder="reliable-go-services" /></label>
                            <label>Tags<input value={draft.tags.join(", ")} onChange={(event) => setDraft({ ...draft, tags: event.target.value.split(",").map((tag) => tag.trim()).filter(Boolean) })} disabled={locked} placeholder="go, reliability" /></label>
                        </div>
                        {isAdmin && <label className="schedule-field">Publication schedule<input type="datetime-local" value={toLocalDateTime(draft.scheduled_at)} min={toLocalDateTime(new Date(Date.now() + 60_000).toISOString())} onChange={(event) => setDraft({ ...draft, scheduled_at: event.target.value ? new Date(event.target.value).toISOString() : "" })} /><small>Uses your local time and publishes automatically at the selected moment.</small></label>}
                        {isAdmin && <div className="auto-publish-field">
                            <small>When this goes live, also auto-publish to:</small>
                            <label><input type="checkbox" checked={draft.publish_devto || false} onChange={(event) => setDraft({ ...draft, publish_devto: event.target.checked })} /> DEV</label>
                            <label><input type="checkbox" checked={draft.publish_linkedin || false} onChange={(event) => setDraft({ ...draft, publish_linkedin: event.target.checked })} /> LinkedIn</label>
                        </div>}
                        <label>Description<textarea value={draft.description} onChange={(event) => setDraft({ ...draft, description: event.target.value })} disabled={locked} rows={3} maxLength={180} /></label>
                        <label className="content-field">Markdown<textarea value={draft.content_markdown} onChange={(event) => setDraft({ ...draft, content_markdown: event.target.value })} disabled={locked} spellCheck /></label>
                        <div className="field-row social-copy">
                            <label>LinkedIn post<textarea value={draft.linkedin_post || ""} onChange={(event) => setDraft({ ...draft, linkedin_post: event.target.value })} disabled={locked} rows={8} maxLength={2500} /></label>
                            <label>Short social post<textarea value={draft.social_post || ""} onChange={(event) => setDraft({ ...draft, social_post: event.target.value })} disabled={locked} rows={8} maxLength={280} /></label>
                        </div>
                        {isAdmin && <section className="publish-panel" aria-label="Publishing destinations">
                            <div><b>Publishing destinations</b><small>External posts use the personal-site article as their canonical URL.</small></div>
                            <label><input type="checkbox" checked={destinations.website} onChange={(event) => setDestinations({ ...destinations, website: event.target.checked })} /> Personal site</label>
                            <label><input type="checkbox" checked={destinations.devto} onChange={(event) => setDestinations({ ...destinations, devto: event.target.checked })} /> DEV</label>
                            <label><input type="checkbox" checked={destinations.linkedin} onChange={(event) => setDestinations({ ...destinations, linkedin: event.target.checked })} /> LinkedIn</label>
                            <button className="command" type="button" onClick={exportForMedium} disabled={!draft.content_markdown}><Download size={16} /> Export for Medium</button>
                        </section>}
                    </div>
                ) : (
                    <article className="preview-pane">
                        <span className="preview-kicker"><BookOpen size={15} /> Read and review before publishing</span>
                        <h1>{draft.title || "Untitled article"}</h1>
                        <p className="preview-description">{draft.description}</p>
                        {draft.scheduled_at && <p className="schedule-summary"><CalendarClock size={15} /> Planned for {formatDateTime(draft.scheduled_at)}</p>}
                        <ReactMarkdown remarkPlugins={[remarkGfm]}>{draft.content_markdown || "Start writing to see the preview."}</ReactMarkdown>
                    </article>
                )}
            </main>
        </div>
    );
}

function toLocalDateTime(value?: string) {
    if (!value) return "";
    const date = new Date(value);
    const local = new Date(date.getTime() - date.getTimezoneOffset() * 60_000);
    return local.toISOString().slice(0, 16);
}

function formatDateTime(value?: string) {
    if (!value) return "the selected time";
    return new Intl.DateTimeFormat(undefined, { dateStyle: "medium", timeStyle: "short" }).format(new Date(value));
}