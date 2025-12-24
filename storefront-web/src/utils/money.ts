import Decimal from "decimal.js-light";

Decimal.set({ precision: 40, rounding: Decimal.ROUND_HALF_UP });

export function money(value: string | number | Decimal): Decimal {
  if (value instanceof Decimal) return value;
  return new Decimal(value);
}

export function formatMoney(
  value: Decimal,
  currency: string,
  locale: string
): string {
  const amount = value.toNumber();
  try {
    return new Intl.NumberFormat(locale, {
      style: "currency",
      currency: currency.toUpperCase(),
      maximumFractionDigits: 2,
    }).format(amount);
  } catch {
    return `${value.toFixed(2)} ${currency.toUpperCase()}`;
  }
}
