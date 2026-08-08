"use client";

import { useEffect, useState } from "react";
import { completeLogin } from "../../../lib/auth";

export default function CallbackPage() {
    const [error, setError] = useState("");

    useEffect(() => {
        const code = new URLSearchParams(window.location.search).get("code");
        if (!code) {
            setError("Authorization code is missing.");
            return;
        }
        completeLogin(code)
            .then(() => window.location.replace("/"))
            .catch((reason: Error) => setError(reason.message));
    }, []);

    return <main className="auth-state"><div className="auth-mark">NS</div><p>{error || "Securing your editorial session…"}</p>{error && <a href="/">Return to sign in</a>}</main>;
}