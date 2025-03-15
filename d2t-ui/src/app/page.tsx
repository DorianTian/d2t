"use client";
import { useState } from "react";
import ReactMarkdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { tomorrow } from "react-syntax-highlighter/dist/cjs/styles/prism";

// Define a type for the result rows
type ResultRow = Record<string, string | number | boolean | null>;

// Define API response type
interface ApiResponse {
  sql: string;
  results: ResultRow[];
  analysis: string;
}

export default function Home() {
  const [query, setQuery] = useState("");
  const [sqlQuery, setSqlQuery] = useState("");
  const [results, setResults] = useState<ResultRow[]>([]);
  const [analysis, setAnalysis] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Function to handle the conversion
  const handleConvert = async () => {
    if (!query.trim()) {
      setError("Please enter a question first");
      return;
    }

    setIsLoading(true);
    setError(null);

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
      setResults(data.results);
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
    }
  };

  // Footer styles to ensure it's fixed at the bottom
  const footerStyle = {
    position: 'fixed',
    bottom: 0,
    left: 0,
    width: '100%',
    backgroundColor: 'white',
    boxShadow: '0 -2px 10px rgba(0, 0, 0, 0.1)',
    borderTop: '1px solid #e5e7eb',
    padding: '1rem 0',
    zIndex: 50,
  } as React.CSSProperties;

  return (
    <div className="min-h-screen bg-gray-100" style={{ paddingBottom: '60px' }}>
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
          <div className="w-full">
            <input
              type="text"
              placeholder="Enter your question"
              className="w-full p-3 border-2 border-gray-400 rounded-md focus:border-blue-600 focus:outline-none bg-gray-50 text-gray-900 placeholder-gray-500"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              disabled={isLoading}
              onKeyDown={(e) => e.key === 'Enter' && !isLoading && handleConvert()}
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
            <h2 className="text-xl font-semibold text-gray-800 mb-4">SQL Query</h2>
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
                      padding: '16px',
                      borderRadius: '0.375rem',
                      fontSize: '0.95rem',
                      lineHeight: '1.5'
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
            <h2 className="text-xl font-semibold text-gray-800 mb-4">Analysis</h2>
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
                            borderRadius: '0.375rem',
                            margin: '1rem 0'
                          }}
                          {...props}
                        >
                          {String(children).replace(/\n$/, "")}
                        </SyntaxHighlighter>
                      ) : (
                        <code className={className + " px-1 py-0.5 bg-gray-100 rounded text-red-600 font-semibold"} {...props}>
                          {children}
                        </code>
                      );
                    },
                    p: ({children}) => <p className="text-gray-800">{children}</p>,
                    li: ({children}) => <li className="text-gray-800">{children}</li>,
                    em: ({children}) => <em className="text-gray-800 font-italic">{children}</em>,
                    strong: ({children}) => <strong className="text-gray-900 font-bold">{children}</strong>,
                    a: ({href, children}) => <a href={href} className="text-blue-700 hover:text-blue-900 underline">{children}</a>,
                    h1: ({children}) => <h1 className="text-gray-900 font-bold">{children}</h1>,
                    h2: ({children}) => <h2 className="text-gray-900 font-bold">{children}</h2>,
                    h3: ({children}) => <h3 className="text-gray-900 font-bold">{children}</h3>,
                    h4: ({children}) => <h4 className="text-gray-900 font-bold">{children}</h4>,
                    table: ({children}) => <table className="border-collapse border border-gray-300 my-4">{children}</table>,
                    th: ({children}) => <th className="border border-gray-300 px-4 py-2 text-gray-900 bg-gray-100">{children}</th>,
                    td: ({children}) => <td className="border border-gray-300 px-4 py-2 text-gray-800">{children}</td>,
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
            ) : (
              <p className="text-gray-500 text-center py-8">
                No results to display
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
