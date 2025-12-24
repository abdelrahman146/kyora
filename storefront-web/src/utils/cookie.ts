export function getCookie(name: string): string | undefined {
  const all = typeof document !== "undefined" ? document.cookie : "";
  if (!all) return undefined;
  const parts = all.split(";");
  for (const p of parts) {
    const [k, ...rest] = p.trim().split("=");
    if (k === name) return decodeURIComponent(rest.join("="));
  }
  return undefined;
}

export function setCookie(name: string, value: string, days = 365): void {
  const expiresAt = new Date(Date.now() + days * 24 * 60 * 60 * 1000);
  const secure =
    typeof location !== "undefined" && location.protocol === "https:";
  const cookie = [
    `${name}=${encodeURIComponent(value)}`,
    `Expires=${expiresAt.toUTCString()}`,
    "Path=/",
    "SameSite=Lax",
    secure ? "Secure" : "",
  ]
    .filter(Boolean)
    .join("; ");
  document.cookie = cookie;
}
