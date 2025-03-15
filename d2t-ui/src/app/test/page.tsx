"use client";
import { useState } from "react";

// Define response type interfaces
interface DebugResponse {
  message: string;
  timestamp: string;
  environment: string | undefined;
  api_base_url: string;
}

interface AskQAResponse {
  sql?: string;
  results?: Array<Record<string, unknown>>;
  analysis?: string;
  error?: string;
}

// Union type for all possible responses
type ApiResponse = DebugResponse | AskQAResponse;

export default function TestPage() {
  const [result, setResult] = useState<ApiResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 测试debug端点
  const testDebug = async () => {
    setLoading(true);
    setError(null);
    
    try {
      console.log("浏览器: 发送请求到 /api/debug");
      const response = await fetch("/api/debug");
      console.log("浏览器: 收到响应", response.status);
      
      const data = await response.json();
      console.log("浏览器: 解析响应数据", data);
      setResult(data);
    } catch (err) {
      console.error("浏览器: 请求错误", err);
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  };

  // 测试askQA端点
  const testAskQA = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const body = { question: "Test question" };
      console.log("浏览器: 发送请求到 /api/askQA", body);
      
      const response = await fetch("/api/askQA", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });
      
      console.log("浏览器: 收到响应", response.status);
      const data = await response.json();
      console.log("浏览器: 解析响应数据", data);
      setResult(data);
    } catch (err) {
      console.error("浏览器: 请求错误", err);
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">API 测试页面</h1>
      
      <div className="space-x-4 mb-6">
        <button 
          onClick={testDebug}
          disabled={loading}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400"
        >
          测试 Debug API
        </button>
        
        <button 
          onClick={testAskQA}
          disabled={loading}
          className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:bg-gray-400"
        >
          测试 AskQA API
        </button>
      </div>
      
      {loading && <p className="text-gray-600">加载中...</p>}
      
      {error && (
        <div className="p-4 bg-red-100 border border-red-400 text-red-700 rounded mb-4">
          <p className="font-bold">错误:</p>
          <p>{error}</p>
        </div>
      )}
      
      {result && (
        <div className="mt-4">
          <h2 className="text-xl font-bold mb-2">响应结果:</h2>
          <pre className="bg-gray-100 p-4 rounded overflow-auto max-h-96">
            {JSON.stringify(result, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
} 