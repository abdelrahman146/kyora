function normalizeWhatsappNumber(input: string): string {
  const trimmed = (input || "").trim();
  // WhatsApp expects digits only (no leading +)
  return trimmed.replace(/[^0-9]/g, "");
}

export function buildWhatsAppLink(numberRaw: string, message: string): string {
  const number = normalizeWhatsappNumber(numberRaw);
  const text = encodeURIComponent(message);

  if (!number) {
    // Fallback to WhatsApp web without phone prefill
    return `https://wa.me/?text=${text}`;
  }

  return `https://wa.me/${number}?text=${text}`;
}

export function isProbablyMobileDevice(): boolean {
  if (typeof navigator === "undefined") return false;
  return /Android|iPhone|iPad|iPod|Mobile/i.test(navigator.userAgent);
}
