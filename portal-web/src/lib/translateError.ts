import type { TFunction } from "i18next";
import type { ErrorResult } from "./errorParser";

/**
 * Translates an ErrorResult using the i18next translation function
 * Falls back to the backend error message or a generic error if translation key is missing
 *
 * @param errorResult - The parsed error result from parseProblemDetails
 * @param t - The i18next translation function (from useTranslation hook)
 * @returns Localized error message
 *
 * @example
 * ```tsx
 * const { t } = useTranslation();
 * const errorResult = await parseProblemDetails(error);
 * const message = translateError(errorResult, t);
 * toast.error(message);
 * ```
 */
export function translateError(errorResult: ErrorResult, t: TFunction): string {
  // Try to get translation with interpolation params
  const translated = t(errorResult.key, {
    defaultValue: errorResult.fallback ?? t("errors:generic.unexpected"),
    ...errorResult.params,
  });

  return translated;
}

/**
 * Shorthand for translating errors inline
 * Useful in catch blocks where you have both error and translation function
 *
 * @example
 * ```tsx
 * try {
 *   await authApi.login(credentials);
 * } catch (error) {
 *   const message = await translateErrorAsync(error, t);
 *   toast.error(message);
 * }
 * ```
 */
export async function translateErrorAsync(
  error: unknown,
  t: TFunction
): Promise<string> {
  const { parseProblemDetails } = await import("./errorParser");
  const errorResult = await parseProblemDetails(error);
  return translateError(errorResult, t);
}
