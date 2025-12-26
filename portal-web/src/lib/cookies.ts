/**
 * Sets a cookie with secure options
 * In production, includes Secure flag for HTTPS-only transmission
 *
 * @param name - Cookie name
 * @param value - Cookie value
 * @param days - Expiration in days (default: 365)
 */
export function setCookie(name: string, value: string, days = 365): void {
  const expires = new Date();
  expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);

  // Add Secure flag in production for HTTPS-only transmission
  const secureFlag = import.meta.env.PROD ? ";Secure" : "";

  document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Lax${secureFlag}`;
}

export function getCookie(name: string): string | null {
  const nameEQ = `${name}=`;
  const ca = document.cookie.split(";");

  for (const cookie of ca) {
    let c = cookie;
    while (c.startsWith(" ")) {
      c = c.substring(1, c.length);
    }
    if (c.startsWith(nameEQ)) {
      return c.substring(nameEQ.length, c.length);
    }
  }

  return null;
}

export function deleteCookie(name: string): void {
  document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`;
}
