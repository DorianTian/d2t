import { NextResponse } from "next/server";

// 用于确认调试API路由模块已加载
console.log("Debug API Route Module Loaded");

/**
 * GET handler for the debug API
 * This endpoint helps test if API routes are working correctly
 */
export async function GET() {
  console.log("Debug GET request received");
  
  return NextResponse.json({
    message: "Debug endpoint is working",
    timestamp: new Date().toISOString(),
    environment: process.env.NODE_ENV,
    api_base_url: process.env.NEXT_PUBLIC_API_BASE_URL || "not set"
  });
} 