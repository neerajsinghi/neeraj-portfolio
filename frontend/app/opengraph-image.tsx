import { ImageResponse } from "next/og";

export const alt = "Neeraj Singhi, Senior Go Backend and AI Engineer";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
    return new ImageResponse(
        (
            <div
                style={{
                    alignItems: "flex-start",
                    background: "#08101d",
                    color: "#f2f6ff",
                    display: "flex",
                    flexDirection: "column",
                    height: "100%",
                    justifyContent: "center",
                    padding: "72px 84px",
                    width: "100%",
                }}
            >
                <div style={{ color: "#45e0c1", display: "flex", fontSize: 25, letterSpacing: 0 }}>
                    BACKEND · DISTRIBUTED SYSTEMS · AI
                </div>
                <div style={{ display: "flex", fontSize: 72, fontWeight: 700, letterSpacing: 0, marginTop: 34 }}>
                    Neeraj Singhi
                </div>
                <div style={{ color: "#aebbd1", display: "flex", fontSize: 43, letterSpacing: 0, marginTop: 18 }}>
                    Senior Go Backend &amp; AI Engineer
                </div>
                <div style={{ borderTop: "2px solid #25334a", color: "#f4bf5f", display: "flex", fontSize: 25, marginTop: 52, paddingTop: 24, width: "100%" }}>
                    10+ years · Go · AWS · RAG · Secure systems
                </div>
            </div>
        ),
        size
    );
}