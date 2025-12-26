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
