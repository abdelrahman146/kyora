import type { Problem } from "./types";

export class ApiError extends Error {
  public readonly status: number;
  public readonly problem?: Problem;

  constructor(message: string, status: number, problem?: Problem) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.problem = problem;
  }
}

function joinUrl(base: string, path: string): string {
  const b = base.replace(/\/$/, "");
  const p = path.startsWith("/") ? path : `/${path}`;
  return `${b}${p}`;
}

export function apiBaseUrl(): string {
  const env = import.meta.env.VITE_API_BASE_URL as string | undefined;
  return (env && env.trim()) || "http://localhost:8080";
}

export async function apiFetch<T>(
  path: string,
  init: RequestInit & { signal?: AbortSignal } = {}
): Promise<T> {
  const url = joinUrl(apiBaseUrl(), path);
  const headers = new Headers(init.headers);
  if (!headers.has("Accept")) headers.set("Accept", "application/json");

  const res = await fetch(url, { ...init, headers });
  const contentType = res.headers.get("content-type") || "";

  const isJson = contentType.includes("application/json");
  if (res.ok) {
    if (res.status === 204) return undefined as T;
    if (!isJson) {
      const text = await res.text();
      return text as unknown as T;
    }
    return (await res.json()) as T;
  }

  let problem: Problem | undefined;
  let message = `Request failed (${res.status})`;

  try {
    if (isJson) {
      problem = (await res.json()) as Problem;
      if (typeof problem?.title === "string" && problem.title.trim()) {
        message = problem.title;
      }
      if (typeof problem?.detail === "string" && problem.detail.trim()) {
        message = `${message}: ${problem.detail}`;
      }
    } else {
      const text = await res.text();
      if (text.trim()) message = text;
    }
  } catch {
    // swallow parse errors
  }

  throw new ApiError(message, res.status, problem);
}
