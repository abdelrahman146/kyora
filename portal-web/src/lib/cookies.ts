export function setCookie(name: string, value: string, days = 365): void {
  const expires = new Date();
  expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
  document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Lax`;
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
