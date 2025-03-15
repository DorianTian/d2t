import { NextRequest, NextResponse } from "next/server";

// 获取API基础URL，优先使用环境变量，否则默认为相对路径
const getApiBaseUrl = () => {
  // 如果环境变量存在且不为空，则使用环境变量
  if (process.env.NEXT_PUBLIC_API_BASE_URL) {
    return process.env.NEXT_PUBLIC_API_BASE_URL;
  }

  // 在生产环境中如果没有设置环境变量，通常是在同一个域中，可以使用相对路径
  // 或者使用Docker Compose中定义的服务名称
  // 这里我们使用go-backend服务名称，当在Docker中运行时
  if (process.env.NODE_ENV === "production") {
    return "http://go-backend:8080";
  }

  // 默认情况下，回退到localhost
  return "http://localhost:8080";
};

/**
 * POST handler for the askQA API
 * This endpoint acts as a proxy to the external service
 */
export async function POST(request: NextRequest) {
  try {
    // Get the request body
    const body = await request.json();

    // 构建API URL
    const apiBaseUrl = getApiBaseUrl();
    console.log(
      "Using Docker Compose service name for API base URL",
      apiBaseUrl
    );

    const apiUrl = `${apiBaseUrl}/api/askQA`;

    // Log the request for debugging
    console.log(`Forwarding request to ${apiUrl}:`, body);

    // Forward the request to the external service - no transformation needed
    // since the backend expects 'question' parameter which is already correct
    const response = await fetch(apiUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body), // Forward the original body without transformation
    });

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
    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Proxy error:", error);
    return NextResponse.json(
      { error: "Failed to proxy the request to the external service" },
      { status: 500 }
    );
  }
}
