import { NextRequest, NextResponse } from "next/server";

// 自定义日志函数，确保日志显示
function devLog(...args: unknown[]): void {
  // 强制将日志写入stderr，这在某些情况下更可靠
  process.stderr.write(`[askQA Route] ${args.map(arg => 
    typeof arg === 'object' ? JSON.stringify(arg) : String(arg)
  ).join(' ')}\n`);
  
  // 也尝试常规控制台日志
  console.log(...args);
}

// 用于确认API路由模块已加载的初始化日志
devLog("API Route Module Loaded: askQA");

// 获取API基础URL，优先使用环境变量，否则默认为相对路径
const getApiBaseUrl = () => {
  devLog("Environment variables:", {
    NEXT_PUBLIC_API_BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL,
    NODE_ENV: process.env.NODE_ENV
  });

  // 如果环境变量存在且不为空，则使用环境变量
  if (process.env.NEXT_PUBLIC_API_BASE_URL) {
    devLog("Using environment variable for API base URL:", process.env.NEXT_PUBLIC_API_BASE_URL);
    return process.env.NEXT_PUBLIC_API_BASE_URL;
  }

  // 在生产环境中如果没有设置环境变量，通常是在同一个域中，可以使用相对路径
  // 或者使用Docker Compose中定义的服务名称
  // 这里我们使用go-backend服务名称，当在Docker中运行时
  if (process.env.NODE_ENV === "production") {
    devLog("Using Docker Compose service name for API base URL: http://go-backend:8080");
    return "http://go-backend:8080";
  }

  // 默认情况下，回退到localhost
  devLog("Using default localhost API base URL: http://localhost:8080");
  return "http://localhost:8080";
};

/**
 * POST handler for the askQA API
 * This endpoint acts as a proxy to the external service
 */
export async function POST(request: NextRequest) {
  devLog("POST request received at /api/askQA");
  
  try {
    // Get the request body
    devLog("Parsing request body...");
    const body = await request.json();
    devLog("Request body parsed:", body);

    // 构建API URL
    const apiBaseUrl = getApiBaseUrl();
    devLog("API Base URL determined:", apiBaseUrl);

    const apiUrl = `${apiBaseUrl}/api/askQA`;
    devLog("Full API URL:", apiUrl);

    // Log the request for debugging
    devLog(`Forwarding request to ${apiUrl}:`, body);

    // Forward the request to the external service - no transformation needed
    // since the backend expects 'question' parameter which is already correct
    devLog("Sending fetch request to backend...");
    const response = await fetch(apiUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body), // Forward the original body without transformation
    });
    devLog("Received response from backend, status:", response.status);

    // If the external service returns an error
    if (!response.ok) {
      const errorText = await response.text();
      console.error(
        `External service error: ${response.status} ${response.statusText}`,
        errorText
      );
      return NextResponse.json(
        {
          error: `External service error: ${response.status} ${response.statusText}`,
        },
        { status: response.status }
      );
    }

    // Return the response from the external service
    devLog("Parsing JSON response from backend...");
    const data = await response.json();
    devLog("Successfully parsed backend response");
    return NextResponse.json(data);
  } catch (error) {
    console.error("Proxy error:", error);
    return NextResponse.json(
      { error: "Failed to proxy the request to the external service" },
      { status: 500 }
    );
  }
}
