import type { StorefrontTheme } from "../api/types";

const BRAND_VARS = [
  "--color-primary",
  "--color-primary-content",
  "--color-secondary",
  "--color-secondary-content",
  "--color-accent",
  "--color-accent-content",
  "--color-base-100",
  "--color-base-200",
  "--color-base-300",
  "--color-base-content",
  "--font-family",
  "--font-family-heading",
] as const;

function parseHexColor(
  input: string
): { r: number; g: number; b: number } | undefined {
  const v = input.trim();
  if (!v.startsWith("#")) return undefined;
  const hex = v.slice(1);
  if (![3, 6].includes(hex.length)) return undefined;

  const full =
    hex.length === 3
      ? `${hex[0]}${hex[0]}${hex[1]}${hex[1]}${hex[2]}${hex[2]}`
      : hex;

  const n = Number.parseInt(full, 16);
  if (Number.isNaN(n)) return undefined;
  return {
    r: (n >> 16) & 255,
    g: (n >> 8) & 255,
    b: n & 255,
  };
}

function relativeLuminance({
  r,
  g,
  b,
}: {
  r: number;
  g: number;
  b: number;
}): number {
  const srgb = [r, g, b].map((x) => x / 255);
  const lin = srgb.map((c) =>
    c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4)
  );
  return 0.2126 * lin[0] + 0.7152 * lin[1] + 0.0722 * lin[2];
}

function bestContentColor(bg: string): string {
  const rgb = parseHexColor(bg);
  if (!rgb) return "oklch(98% 0 0)";
  const lum = relativeLuminance(rgb);
  return lum > 0.55 ? "oklch(20% 0 0)" : "oklch(98% 0 0)";
}

function setVar(name: string, value: string | undefined): void {
  if (!value) return;
  const v = value.trim();
  if (!v) return;
  document.documentElement.style.setProperty(name, v);
}

function resetBrandTheme(): void {
  for (const v of BRAND_VARS) {
    document.documentElement.style.removeProperty(v);
  }
}

export function applyBrandTheme(theme: StorefrontTheme | undefined): void {
  // Always reset first so switching storefronts doesn't keep the previous storefront's colors.
  resetBrandTheme();

  // No overrides: keep the default daisyUI theme.
  if (!theme) return;

  // daisyUI color variables (semantic). If a field is missing we keep the current theme defaults.
  setVar("--color-primary", theme.primaryColor);
  setVar("--color-secondary", theme.secondaryColor);
  setVar("--color-accent", theme.accentColor);

  if (theme.primaryColor)
    setVar("--color-primary-content", bestContentColor(theme.primaryColor));
  if (theme.secondaryColor)
    setVar("--color-secondary-content", bestContentColor(theme.secondaryColor));
  if (theme.accentColor)
    setVar("--color-accent-content", bestContentColor(theme.accentColor));

  // Base surfaces
  if (theme.backgroundColor) {
    setVar("--color-base-100", theme.backgroundColor);
    // Keep base surfaces consistent even if only base-100 is provided.
    // This reduces "mismatched" borders and section backgrounds.
    setVar(
      "--color-base-200",
      `color-mix(in oklab, ${theme.backgroundColor} 94%, black 6%)`
    );
    setVar(
      "--color-base-300",
      `color-mix(in oklab, ${theme.backgroundColor} 88%, black 12%)`
    );
  }
  setVar("--color-base-content", theme.textColor);

  // Typography (optional): prefer system fonts unless explicitly configured.
  if (theme.fontFamily) setVar("--font-family", theme.fontFamily);
  if (theme.headingFontFamily)
    setVar("--font-family-heading", theme.headingFontFamily);
}
