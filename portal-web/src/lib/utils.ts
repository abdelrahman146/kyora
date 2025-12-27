// Utility functions and helpers

/**
 * Combines class names conditionally
 * Similar to clsx but lightweight
 */
export const cn = (
  ...classes: (string | undefined | null | false)[]
): string => {
  return classes.filter(Boolean).join(" ");
};

export const formatCurrency = (amount: number, currency = "SAR"): string => {
  return new Intl.NumberFormat("ar-SA", {
    style: "currency",
    currency,
  }).format(amount);
};

export const formatDate = (date: Date | string, locale = "ar"): string => {
  const dateObj = typeof date === "string" ? new Date(date) : date;
  return new Intl.DateTimeFormat(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(dateObj);
};

export const debounce = <T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: ReturnType<typeof setTimeout>;
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn(...args), delay);
  };
};

export const formatCountdownDuration = (totalSeconds: number): string => {
  if (!Number.isFinite(totalSeconds)) return "0s";

  const secondsInt = Math.max(0, Math.floor(totalSeconds));
  const hours = Math.floor(secondsInt / 3600);
  const minutes = Math.floor((secondsInt % 3600) / 60);
  const seconds = secondsInt % 60;

  if (hours > 0) {
    return `${String(hours)}hr:${String(minutes).padStart(2, "0")}m:${String(
      seconds
    ).padStart(2, "0")}s`;
  }

  if (minutes > 0) {
    return `${String(minutes)}m:${String(seconds).padStart(2, "0")}s`;
  }

  return `${String(seconds)}s`;
};
