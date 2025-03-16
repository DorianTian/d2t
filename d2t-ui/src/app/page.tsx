"use client";
import { useState } from "react";
import ReactMarkdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { tomorrow } from "react-syntax-highlighter/dist/cjs/styles/prism";
import Link from "next/link";

// Define a type for the result rows
type ResultRow = Record<string, string | number | boolean | null>;

// Define API response type
interface ApiResponse {
  sql: string;
  results: ResultRow[];
  analysis: string;
}

// Example prompts for guidance
const examplePrompts = [
  "帮我找到订单数量为100的产品都有哪一些",
  "查询所有销售额超过1000元的订单",
  "统计每个月的销售总额",
  "找出库存少于20件的产品",
  "哪些顾客在过去3个月有过购买记录",
];

export default function Home() {
  const [query, setQuery] = useState("");
  const [sqlQuery, setSqlQuery] = useState("");
  const [results, setResults] = useState<ResultRow[]>([]);
  const [analysis, setAnalysis] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showExamples, setShowExamples] = useState(false);
  const [showLoadingPopup, setShowLoadingPopup] = useState(false);

  // Function to handle the conversion
  const handleConvert = async () => {
    if (!query.trim()) {
      setError("Please enter a question first");
      return;
    }

    setIsLoading(true);
    setError(null);
    setShowLoadingPopup(true);

    try {
      const response = await fetch("/api/askQA", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ question: query }),
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      const data: ApiResponse = await response.json();

      setSqlQuery(data.sql);
      setResults(data.results || []);
      setAnalysis(data.analysis);
    } catch (err) {
      console.error("Error calling API:", err);
      setError(
        `Failed to process query: ${
          err instanceof Error ? err.message : "Unknown error"
        }`
      );
      // Keep previous results if any
    } finally {
      setIsLoading(false);
      setShowLoadingPopup(false);
    }
  };

  // Set the query from example
  const setExampleQuery = (example: string) => {
    setQuery(example);
    setShowExamples(false);
  };

  // Footer styles to ensure it's fixed at the bottom
  const footerStyle = {
    position: "fixed",
    bottom: 0,
    left: 0,
    width: "100%",
    backgroundColor: "white",
    boxShadow: "0 -2px 10px rgba(0, 0, 0, 0.1)",
    borderTop: "1px solid #e5e7eb",
    padding: "1rem 0",
    zIndex: 50,
  } as React.CSSProperties;

  // Navigation styles to ensure it's fixed at the top
  const navStyle = {
    position: "fixed",
    top: 0,
    left: 0,
    width: "100%",
    zIndex: 40,
    backgroundColor: "#1f2937", // bg-gray-800
    boxShadow: "0 2px 10px rgba(0, 0, 0, 0.1)",
  } as React.CSSProperties;

  return (
    <div
      className="min-h-screen bg-gray-100"
      style={{ paddingBottom: "60px", paddingTop: "64px" }}
    >
      {/* Loading Popup */}
      {showLoadingPopup && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg shadow-xl text-center">
            <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-blue-500 border-solid mx-auto mb-4"></div>
            <p className="text-lg font-medium text-gray-700">
              正在deepseek查询中
            </p>
          </div>
        </div>
      )}

      {/* Navigation */}
      <nav style={navStyle}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <span className="text-xl font-bold text-white">D2T</span>
              </div>
              <div className="ml-10 flex items-baseline space-x-4">
                <Link
                  href="/"
                  className="px-3 py-2 rounded-md text-sm font-medium bg-gray-900 text-white"
                >
                  首页
                </Link>
                <Link
                  href="/test"
                  className="px-3 py-2 rounded-md text-sm font-medium text-gray-300 hover:bg-gray-700 hover:text-white"
                >
                  测试
                </Link>
              </div>
            </div>
          </div>
        </div>
      </nav>

      {/* Header */}
      <header className="bg-white shadow-sm py-4">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900 text-center">D2T</h1>
          <p className="text-gray-600 text-center mt-1">
            A tool for converting natural language to SQL queries
          </p>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8 w-full">
        {/* Query input */}
        <div className="mb-8 bg-white rounded-lg shadow p-6">
          <div className="w-full relative">
            <div className="flex items-center mb-2">
              <label className="text-sm font-medium text-gray-600 mr-2">
                输入问题:
              </label>
              <button
                type="button"
                className="text-blue-600 text-sm hover:text-blue-800 focus:outline-none"
                onClick={() => setShowExamples(!showExamples)}
              >
                查看示例
              </button>
            </div>

            {showExamples && (
              <div className="absolute z-10 mt-1 w-full bg-white border border-gray-300 rounded-md shadow-lg">
                <ul className="py-1 max-h-60 overflow-auto">
                  {examplePrompts.map((example, index) => (
                    <li
                      key={index}
                      className="px-4 py-2 hover:bg-gray-100 cursor-pointer text-gray-700"
                      onClick={() => setExampleQuery(example)}
                    >
                      {example}
                    </li>
                  ))}
                </ul>
              </div>
            )}

            <input
              type="text"
              placeholder="输入您的问题，例如：帮我找到订单数量为100的产品"
              className="w-full p-3 border-2 border-gray-400 rounded-md focus:border-blue-600 focus:outline-none bg-gray-50 text-gray-900 placeholder-gray-500"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              disabled={isLoading}
              onKeyDown={(e) =>
                e.key === "Enter" && !isLoading && handleConvert()
              }
              onFocus={() => setError(null)}
            />
            {error && <p className="text-red-500 mt-2 text-sm">{error}</p>}
          </div>
          <button
            className={`mt-4 px-6 py-3 rounded-md text-white font-medium transition-colors ${
              isLoading
                ? "bg-gray-400 cursor-not-allowed"
                : "bg-blue-600 hover:bg-blue-700"
            }`}
            onClick={handleConvert}
            disabled={isLoading}
          >
            {isLoading ? "Processing..." : "Convert"}
          </button>
        </div>

        {/* SQL Query and Analysis - Side by side */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* SQL Query */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              SQL Query
            </h2>
            <div className="w-full border border-gray-200 rounded-md overflow-hidden">
              {isLoading ? (
                <p className="text-gray-500 p-4">Generating SQL query...</p>
              ) : sqlQuery ? (
                <div className="bg-gray-800 rounded-md">
                  <SyntaxHighlighter
                    language="sql"
                    style={tomorrow}
                    customStyle={{
                      margin: 0,
                      padding: "16px",
                      borderRadius: "0.375rem",
                      fontSize: "0.95rem",
                      lineHeight: "1.5",
                    }}
                  >
                    {sqlQuery}
                  </SyntaxHighlighter>
                </div>
              ) : (
                <p className="text-gray-500 p-4">No SQL query generated yet</p>
              )}
            </div>
          </div>

          {/* Analysis */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              Analysis
            </h2>
            <div className="w-full border border-gray-200 rounded-md overflow-auto min-h-[200px] max-h-[400px] p-4 prose prose-sm prose-strong:text-gray-900 prose-headings:text-gray-900 prose-p:text-gray-800">
              {isLoading ? (
                <p className="text-gray-500 text-center py-8">
                  Generating analysis...
                </p>
              ) : (
                <ReactMarkdown
                  components={{
                    code: ({ className, children, ...props }) => {
                      const match = /language-(\w+)/.exec(className || "");
                      return match ? (
                        <SyntaxHighlighter
                          // @ts-expect-error - Type compatibility issue between react-markdown and react-syntax-highlighter
                          style={tomorrow}
                          language={match[1]}
                          PreTag="div"
                          customStyle={{
                            borderRadius: "0.375rem",
                            margin: "1rem 0",
                          }}
                          {...props}
                        >
                          {String(children).replace(/\n$/, "")}
                        </SyntaxHighlighter>
                      ) : (
                        <code
                          className={
                            className +
                            " px-1 py-0.5 bg-gray-100 rounded text-red-600 font-semibold"
                          }
                          {...props}
                        >
                          {children}
                        </code>
                      );
                    },
                    p: ({ children }) => (
                      <p className="text-gray-800">{children}</p>
                    ),
                    li: ({ children }) => (
                      <li className="text-gray-800">{children}</li>
                    ),
                    em: ({ children }) => (
                      <em className="text-gray-800 font-italic">{children}</em>
                    ),
                    strong: ({ children }) => (
                      <strong className="text-gray-900 font-bold">
                        {children}
                      </strong>
                    ),
                    a: ({ href, children }) => (
                      <a
                        href={href}
                        className="text-blue-700 hover:text-blue-900 underline"
                      >
                        {children}
                      </a>
                    ),
                    h1: ({ children }) => (
                      <h1 className="text-gray-900 font-bold">{children}</h1>
                    ),
                    h2: ({ children }) => (
                      <h2 className="text-gray-900 font-bold">{children}</h2>
                    ),
                    h3: ({ children }) => (
                      <h3 className="text-gray-900 font-bold">{children}</h3>
                    ),
                    h4: ({ children }) => (
                      <h4 className="text-gray-900 font-bold">{children}</h4>
                    ),
                    table: ({ children }) => (
                      <table className="border-collapse border border-gray-300 my-4">
                        {children}
                      </table>
                    ),
                    th: ({ children }) => (
                      <th className="border border-gray-300 px-4 py-2 text-gray-900 bg-gray-100">
                        {children}
                      </th>
                    ),
                    td: ({ children }) => (
                      <td className="border border-gray-300 px-4 py-2 text-gray-800">
                        {children}
                      </td>
                    ),
                  }}
                >
                  {analysis}
                </ReactMarkdown>
              )}
            </div>
          </div>
        </div>

        {/* Results - Full width */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Results</h2>
          <div className="w-full overflow-auto">
            {isLoading ? (
              <p className="text-gray-500 text-center py-8">
                Fetching results...
              </p>
            ) : results.length > 0 ? (
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-gray-700 text-white">
                    {Object.keys(results[0] || {}).map((key) => (
                      <th
                        key={key}
                        className="p-3 text-left font-semibold border border-gray-600"
                      >
                        {key}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {results.map((row, rowIndex) => (
                    <tr
                      key={rowIndex}
                      className={rowIndex % 2 === 0 ? "bg-white" : "bg-gray-50"}
                    >
                      {Object.entries(row).map(([key, value]) => (
                        <td
                          key={`${rowIndex}-${key}`}
                          className="p-3 border border-gray-300 text-gray-700"
                        >
                          {String(value)}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : sqlQuery ? (
              <div className="text-center py-12 border-2 border-dashed border-gray-300 rounded-lg">
                <svg
                  className="mx-auto h-12 w-12 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                  />
                </svg>
                <h3 className="mt-2 text-sm font-medium text-gray-900">
                  查询成功，但未找到符合条件的数据
                </h3>
                <p className="mt-1 text-sm text-gray-500">
                  您的查询已成功执行，但没有符合条件的数据记录。
                </p>
                <p className="mt-1 text-sm text-gray-500">
                  您可以尝试调整查询条件或查看SQL语句以了解更多信息。
                </p>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-8">
                请输入问题开始查询
              </p>
            )}
          </div>
        </div>
      </main>

      {/* Footer with inline styles to ensure fixed position */}
      <footer style={footerStyle}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <p className="text-center text-gray-600 font-medium">
            &copy; 2025 D2T - Natural Language to SQL Query Tool
          </p>
        </div>
      </footer>
    </div>
  );
}
