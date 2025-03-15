import { NextResponse } from "next/server";

/**
 * GET handler for health check endpoint
 * Used by Docker and AWS for health monitoring
 */
export async function GET() {
  return NextResponse.json({ status: 'ok', timestamp: new Date().toISOString() }, { status: 200 });
} 