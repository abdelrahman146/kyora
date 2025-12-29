export function normalizePhoneCode(raw: string): string {
  const trimmed = raw.trim().replace(/\s+/g, "");
  if (trimmed === "") return "";
  return trimmed.startsWith("+") ? trimmed : `+${trimmed}`;
}

export function normalizePhoneNumber(raw: string): string {
  return raw.replace(/[^0-9]/g, "");
}

export function buildE164Phone(
  phoneCodeRaw: string,
  phoneNumberRaw: string
): {
  phoneCode: string;
  phoneNumber: string;
  e164: string;
} {
  const phoneCode = normalizePhoneCode(phoneCodeRaw);
  const phoneNumber = normalizePhoneNumber(phoneNumberRaw);
  return {
    phoneCode,
    phoneNumber,
    e164: `${phoneCode}${phoneNumber}`,
  };
}

export function parseE164Phone(
  phoneCode: string,
  phoneNumber: string
): {
  phoneCode: string;
  phoneNumber: string;
} {
  return {
    phoneCode: normalizePhoneCode(phoneCode),
    phoneNumber: normalizePhoneNumber(phoneNumber),
  };
}
