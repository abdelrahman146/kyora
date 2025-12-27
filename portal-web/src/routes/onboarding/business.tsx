import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Building2, Globe, DollarSign } from "lucide-react";
import { OnboardingLayout } from "@/components/templates";
import { useOnboarding } from "@/contexts/OnboardingContext";
import { onboardingApi } from "@/api/onboarding";
import { translateErrorAsync } from "@/lib/translateError";

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
 * 2. POST /api/onboarding/business
 * 3. Navigate to /onboarding/payment (paid) or /onboarding/complete (free)
 */

// Common countries for business registration
const COUNTRIES = [
  { code: "US", name: "United States" },
  { code: "GB", name: "United Kingdom" },
  { code: "CA", name: "Canada" },
  { code: "AU", name: "Australia" },
  { code: "DE", name: "Germany" },
  { code: "FR", name: "France" },
  { code: "ES", name: "Spain" },
  { code: "IT", name: "Italy" },
  { code: "NL", name: "Netherlands" },
  { code: "SE", name: "Sweden" },
  { code: "AE", name: "United Arab Emirates" },
  { code: "SA", name: "Saudi Arabia" },
  { code: "EG", name: "Egypt" },
  { code: "JO", name: "Jordan" },
  { code: "LB", name: "Lebanon" },
];

// Common currencies
const CURRENCIES = [
  { code: "USD", name: "US Dollar", symbol: "$" },
  { code: "EUR", name: "Euro", symbol: "€" },
  { code: "GBP", name: "British Pound", symbol: "£" },
  { code: "CAD", name: "Canadian Dollar", symbol: "C$" },
  { code: "AUD", name: "Australian Dollar", symbol: "A$" },
  { code: "AED", name: "UAE Dirham", symbol: "د.إ" },
  { code: "SAR", name: "Saudi Riyal", symbol: "﷼" },
  { code: "EGP", name: "Egyptian Pound", symbol: "E£" },
  { code: "JOD", name: "Jordanian Dinar", symbol: "د.ا" },
  { code: "LBP", name: "Lebanese Pound", symbol: "ل.ل" },
];

export default function BusinessSetupPage() {
  const { t } = useTranslation(["onboarding", "common"]);
  const navigate = useNavigate();
  const {
    sessionToken,
    stage,
    isPaidPlan,
    setBusinessDetails,
  } = useOnboarding();

  const [businessName, setBusinessName] = useState("");
  const [descriptor, setDescriptor] = useState("");
  const [country, setCountry] = useState("");
  const [currency, setCurrency] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [descriptorError, setDescriptorError] = useState("");

  // Redirect if not verified
  useEffect(() => {
    if (!sessionToken || stage !== "identity_verified") {
      navigate("/onboarding/verify", { replace: true });
    }
  }, [sessionToken, stage, navigate]);

  // Auto-generate descriptor from business name
  useEffect(() => {
    if (businessName && !descriptor) {
      const generated = businessName
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, "")
        .replace(/\s+/g, "-")
        .slice(0, 50);
      setDescriptor(generated);
    }
  }, [businessName, descriptor]);

  // Validate descriptor format
  useEffect(() => {
    if (descriptor) {
      if (descriptor.length < 3) {
        setDescriptorError(t("onboarding.business.descriptorTooShort"));
      } else if (!/^[a-z0-9-]+$/.test(descriptor)) {
        setDescriptorError(t("onboarding.business.descriptorInvalidFormat"));
      } else {
        setDescriptorError("");
      }
    } else {
      setDescriptorError("");
    }
  }, [descriptor, t]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!businessName.trim()) {
      setError(t("onboarding.business.nameRequired"));
      return;
    }

    if (!descriptor.trim()) {
      setError(t("onboarding.business.descriptorRequired"));
      return;
    }

    if (descriptorError) {
      setError(descriptorError);
      return;
    }

    if (!country) {
      setError(t("onboarding.business.countryRequired"));
      return;
    }

    if (!currency) {
      setError(t("onboarding.business.currencyRequired"));
      return;
    }

    if (!sessionToken) return;

    try {
      setIsSubmitting(true);

      await onboardingApi.setBusiness({
        sessionToken,
        name: businessName.trim(),
        descriptor: descriptor.trim(),
        country,
        currency,
      });

      setBusinessDetails(businessName.trim(), descriptor.trim(), country, currency);

      // Navigate to payment if paid plan, otherwise complete
      if (isPaidPlan) {
        navigate("/onboarding/payment");
      } else {
        navigate("/onboarding/complete");
      }
    } catch (err) {
      const message = await translateErrorAsync(err, t);
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <OnboardingLayout currentStep={2} totalSteps={5}>
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
                {t("onboarding.business.title")}
              </h2>
              <p className="text-base-content/70 mt-2">
                {t("onboarding.business.subtitle")}
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Business Name */}
              <div className="form-control">
                <label htmlFor="businessName" className="label">
                  <span className="label-text font-medium">
                    {t("onboarding.business.name")}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  id="businessName"
                  value={businessName}
                  onChange={(e) => setBusinessName(e.target.value)}
                  placeholder={t("onboarding.business.namePlaceholder")}
                  className="input input-bordered"
                  required
                  disabled={isSubmitting}
                />
                <label className="label">
                  <span className="label-text-alt">
                    {t("onboarding.business.nameHint")}
                  </span>
                </label>
              </div>

              {/* Business Descriptor */}
              <div className="form-control">
                <label htmlFor="descriptor" className="label">
                  <span className="label-text font-medium">
                    {t("onboarding.business.descriptor")}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  id="descriptor"
                  value={descriptor}
                  onChange={(e) => setDescriptor(e.target.value.toLowerCase())}
                  placeholder={t("onboarding.business.descriptorPlaceholder")}
                  className={`input input-bordered ${
                    descriptorError ? "input-error" : ""
                  }`}
                  pattern="[a-z0-9-]+"
                  minLength={3}
                  maxLength={50}
                  required
                  disabled={isSubmitting}
                />
                <label className="label">
                  <span className={`label-text-alt ${descriptorError ? "text-error" : ""}`}>
                    {descriptorError || t("onboarding.business.descriptorHint")}
                  </span>
                </label>
              </div>

              {/* Country & Currency Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Country */}
                <div className="form-control">
                  <label htmlFor="country" className="label">
                    <span className="label-text font-medium flex items-center gap-2">
                      <Globe className="w-4 h-4" />
                      {t("onboarding.business.country")}
                    </span>
                    <span className="label-text-alt text-error">*</span>
                  </label>
                  <select
                    id="country"
                    value={country}
                    onChange={(e) => setCountry(e.target.value)}
                    className="select select-bordered"
                    required
                    disabled={isSubmitting}
                  >
                    <option value="">
                      {t("onboarding.business.selectCountry")}
                    </option>
                    {COUNTRIES.map((c) => (
                      <option key={c.code} value={c.code}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Currency */}
                <div className="form-control">
                  <label htmlFor="currency" className="label">
                    <span className="label-text font-medium flex items-center gap-2">
                      <DollarSign className="w-4 h-4" />
                      {t("onboarding.business.currency")}
                    </span>
                    <span className="label-text-alt text-error">*</span>
                  </label>
                  <select
                    id="currency"
                    value={currency}
                    onChange={(e) => setCurrency(e.target.value)}
                    className="select select-bordered"
                    required
                    disabled={isSubmitting}
                  >
                    <option value="">
                      {t("onboarding.business.selectCurrency")}
                    </option>
                    {CURRENCIES.map((c) => (
                      <option key={c.code} value={c.code}>
                        {c.symbol} {c.name} ({c.code})
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              {error && (
                <div className="alert alert-error">
                  <span className="text-sm">{error}</span>
                </div>
              )}

              {/* Actions */}
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => navigate("/onboarding/verify")}
                  className="btn btn-ghost"
                  disabled={isSubmitting}
                >
                  {t("common.back")}
                </button>
                <button
                  type="submit"
                  className="btn btn-primary flex-1"
                  disabled={isSubmitting || !!descriptorError}
                >
                  {isSubmitting ? (
                    <>
                      <span className="loading loading-spinner loading-sm"></span>
                      {t("common.loading")}
                    </>
                  ) : (
                    t("common.continue")
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
