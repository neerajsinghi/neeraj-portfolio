import { revalidateTag } from "next/cache";
import { NextRequest, NextResponse } from "next/server";
import { timingSafeEqual } from "crypto";

// Called by the backend right after a publish-affecting change (create,
// update, delete, or the scheduled-publish cron) so the blog list/detail
// pages update immediately instead of waiting on the 60s time window.
export async function POST(request: NextRequest) {
  const expected = process.env.REVALIDATE_TOKEN;
  const provided = request.headers.get("authorization")?.replace(/^Bearer\s+/i, "") ?? "";
  if (!expected || !tokensMatch(provided, expected)) {
    return NextResponse.json({ error: "valid revalidate token required" }, { status: 401 });
  }

  revalidateTag("blogs", "max");
  return NextResponse.json({ revalidated: true });
}

function tokensMatch(a: string, b: string): boolean {
  const bufferA = Buffer.from(a);
  const bufferB = Buffer.from(b);
  if (bufferA.length !== bufferB.length) return false;
  return timingSafeEqual(bufferA, bufferB);
}
