import { useMemo, useState, useEffect, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Building2, Globe, DollarSign } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { FormInput, FormSelect } from "@/components";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { translateErrorAsync } from "@/lib/translateError";
import { useMetadataStore } from "@/stores/metadataStore";

/**
 * Business Setup Step - Step 3 of Onboarding
 *
 * Features:
 * - Business name input
 * - Business descriptor (slug) validation
 * - Country selection
 * - Currency selection
 * - Real-time descriptor availability check
 *
 * Flow:
 * 1. User provides business details
 * 2. POST /v1/onboarding/business
 * 3. Navigate to /onboarding/payment (paid) or /onboarding/complete (free)
 */

export default function BusinessSetupPage() {
  const { t, i18n } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const {
    sessionToken,
    stage,
    isPaidPlan,
    loadSessionFromStorage,
  } = useOnboarding();

  const [businessName, setBusinessName] = useState("");
  const [descriptor, setDescriptor] = useState("");
  const [country, setCountry] = useState("");
  const [currency, setCurrency] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [descriptorError, setDescriptorError] = useState("");
  const [isDescriptorManuallyEdited, setIsDescriptorManuallyEdited] = useState(false);

  const countries = useMetadataStore((s) => s.countries);
  const countriesStatus = useMetadataStore((s) => s.status);
  const loadCountries = useMetadataStore((s) => s.loadCountries);

  const isArabic = i18n.language.toLowerCase().startsWith("ar");

  useEffect(() => {
    void loadCountries();
  }, [loadCountries]);

  const countryByCode = useMemo(() => {
    const map = new Map<string, (typeof countries)[number]>();
    for (const c of countries) {
      map.set(c.code, c);
    }
    return map;
  }, [countries]);

  const countryOptions = useMemo(() => {
    return countries.map((c) => {
      const label = `${c.flag ? `${c.flag} ` : ""}${isArabic ? c.nameAr : c.name}`;
      return {
        value: c.code,
        label,
        icon: <Globe className="w-4 h-4" />,
      };
    });
  }, [countries, isArabic]);

  const currencyOptions = useMemo(() => {
    const seen = new Set<string>();
    const options: { value: string; label: string; icon: ReactNode }[] = [];
    for (const c of countries) {
      if (!c.currencyCode || seen.has(c.currencyCode)) continue;
      seen.add(c.currencyCode);
      options.push({
        value: c.currencyCode,
        label: c.currencyLabel || c.currencyCode,
        icon: <DollarSign className="w-4 h-4" />,
      });
    }
    return options;
  }, [countries]);

  // Restore session from localStorage on mount
  useEffect(() => {
    const restoreSession = async () => {
      if (!sessionToken) {
        const hasSession = await loadSessionFromStorage();
        if (!hasSession) {
          await navigate("/onboarding/plan", { replace: true });
        }
      }
    };
    void restoreSession();
  }, [sessionToken, loadSessionFromStorage, navigate]);

  // Redirect if not verified
  useEffect(() => {
    if (!sessionToken || stage !== "identity_verified") {
      void navigate("/onboarding/verify", { replace: true });
    }
  }, [sessionToken, stage, navigate]);

  // Auto-generate descriptor from business name (only if not manually edited)
  useEffect(() => {
    if (businessName && !isDescriptorManuallyEdited) {
      const generated = businessName
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, "")
        .replace(/\s+/g, "-")
        .slice(0, 50);
      setDescriptor(generated);
    }
  }, [businessName, isDescriptorManuallyEdited]);

  // Validate descriptor format
  useEffect(() => {
    if (descriptor) {
      if (descriptor.length < 3) {
        setDescriptorError(t("onboarding:business.descriptorTooShort"));
      } else if (!/^[a-z0-9-]+$/.test(descriptor)) {
        setDescriptorError(t("onboarding:business.descriptorInvalidFormat"));
      } else {
        setDescriptorError("");
      }
    } else {
      setDescriptorError("");
    }
  }, [descriptor, t]);

  const submitBusiness = async () => {
    setError("");

    if (!businessName.trim()) {
      setError(t("onboarding:business.nameRequired"));
      return;
    }

    if (!descriptor.trim()) {
      setError(t("onboarding:business.descriptorRequired"));
      return;
    }

    if (descriptorError) {
      setError(descriptorError);
      return;
    }

    if (!country) {
      setError(t("onboarding:business.countryRequired"));
      return;
    }

    if (!currency) {
      setError(t("onboarding:business.currencyRequired"));
      return;
    }

    if (!sessionToken) return;

    try {
      setIsSubmitting(true);

      const response = await onboardingApi.setBusiness({
        sessionToken,
        name: businessName.trim(),
        descriptor: descriptor.trim(),
        country,
        currency,
      });

      // Backend updates the session automatically, just navigate based on response
      // Navigate based on the response stage
      if (response.stage === "ready_to_commit") {
        // Free plan - go directly to complete
        void navigate("/onboarding/complete");
      } else if (response.stage === "business_staged" || isPaidPlan) {
        // Paid plan - go to payment
        void navigate("/onboarding/payment");
      } else {
        // Fallback to complete
        void navigate("/onboarding/complete");
      }
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    void submitBusiness();
  };

  return (
    <OnboardingLayout currentStep={3} totalSteps={5}>
      <div className="max-w-2xl mx-auto">
        <div className="card bg-base-100 border border-base-300 shadow-xl">
          <div className="card-body">
            {/* Header */}
            <div className="text-center mb-6">
              <div className="flex justify-center mb-4">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                  <Building2 className="w-8 h-8 text-primary" />
                </div>
              </div>
              <h2 className="text-2xl font-bold">
                {t("onboarding:business.title")}
              </h2>
              <p className="text-base-content/70 mt-2">
                {t("onboarding:business.subtitle")}
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Business Name */}
              <FormInput
                label={t("onboarding:business.name")}
                type="text"
                value={businessName}
                onChange={(e) => {
                  setBusinessName(e.target.value);
                }}
                placeholder={t("onboarding:business.namePlaceholder")}
                required
                disabled={isSubmitting}
                startIcon={<Building2 className="w-5 h-5" />}
                helperText={t("onboarding:business.nameHint")}
              />

              {/* Business Descriptor */}
              <FormInput
                label={t("onboarding:business.descriptor")}
                type="text"
                value={descriptor}
                onChange={(e) => {
                  setDescriptor(e.target.value.toLowerCase());
                  setIsDescriptorManuallyEdited(true);
                }}
                onBlur={() => {
                  // If descriptor is empty after blur, allow auto-sync again
                  if (!descriptor.trim()) {
                    setIsDescriptorManuallyEdited(false);
                  }
                }}
                placeholder={t("onboarding:business.descriptorPlaceholder")}
                pattern="[a-z0-9-]+"
                minLength={3}
                maxLength={50}
                required
                disabled={isSubmitting}
                error={descriptorError}
                helperText={!descriptorError ? t("onboarding:business.descriptorHint") : undefined}
              />

              {/* Country & Currency Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Country */}
                <FormSelect<string>
                  label={t("onboarding:business.country")}
                  options={countryOptions}
                  value={country}
                  onChange={(value) => {
                    const code = value as string;
                    setCountry(code);
                    const selected = countryByCode.get(code);
                    if (selected?.currencyCode) {
                      setCurrency(selected.currencyCode);
                    }
                  }}
                  required
                  disabled={isSubmitting || countriesStatus !== "loaded"}
                  placeholder={t("onboarding:business.selectCountry")}
                  searchable
                />

                {/* Currency */}
                <FormSelect<string>
                  label={t("onboarding:business.currency")}
                  options={currencyOptions}
                  value={currency}
                  onChange={(value) => {
                    setCurrency(value as string);
                  }}
                  required
                  disabled={isSubmitting || countriesStatus !== "loaded"}
                  placeholder={t("onboarding:business.selectCurrency")}
                  searchable
                />
              </div>

              {error && (
                <div className="alert alert-error">
                  <span className="text-sm">{error}</span>
                </div>
              )}

              {/* Actions */}
              <div>
                <button
                  type="submit"
                  className="btn btn-primary w-full"
                  disabled={isSubmitting || !!descriptorError}
                >
                  {isSubmitting ? (
                    <>
                      <span className="loading loading-spinner loading-sm"></span>
                      {t("common:loading")}
                    </>
                  ) : (
                    t("common:continue")
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </OnboardingLayout>
  );
}
