import { useEffect, useId, useMemo } from "react";
import { useTranslation } from "react-i18next";
import toast from "react-hot-toast";
import { Controller, useForm, useWatch } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { BottomSheet } from "../../molecules/BottomSheet";
import { FormInput, FormSelect, FormTextarea } from "@/components";
import { createCustomer, createCustomerAddress } from "@/api/customer";
import type {
  CreateCustomerAddressRequest,
  CustomerGender,
  Customer,
  CustomerAddress,
} from "@/api/types/customer";
import { translateErrorAsync } from "@/lib/translateError";
import { buildE164Phone } from "@/lib/phone";
import { useMetadataStore } from "@/stores/metadataStore";

export interface AddCustomerSheetProps {
  isOpen: boolean;
  onClose: () => void;
  businessDescriptor: string;
  businessCountryCode: string;
  onCreated?: (customer: Customer) => void | Promise<void>;
}

const addCustomerSchema = z
  .object({
    name: z.string().trim().min(1, "validation.required"),
    email: z
      .string()
      .trim()
      .min(1, "validation.required")
      .pipe(z.email("validation.invalid_email")),
    gender: z.enum(["male", "female", "other"], { message: "validation.required" }),
    countryCode: z
      .string()
      .trim()
      .min(1, "validation.required")
      .refine((v) => /^[A-Za-z]{2}$/.test(v), "validation.invalid_country"),
    phoneCode: z
      .string()
      .trim()
      .refine((v) => v === "" || /^\+?\d{1,4}$/.test(v), "validation.invalid_phone_code"),
    phoneNumber: z
      .string()
      .trim()
      .refine((v) => v === "" || /^[0-9\-\s()]{6,20}$/.test(v), "validation.invalid_phone"),
    street: z.string().trim(),
    city: z.string().trim(),
    state: z.string().trim(),
    zipCode: z.string().trim(),
  })
  .refine(
    (values) => {
      const hasPhoneNumber = values.phoneNumber.trim() !== "";
      return !hasPhoneNumber || values.phoneCode.trim() !== "";
    },
    { message: "validation.required", path: ["phoneCode"] }
  )
  .refine(
    (values) => {
      const hasAddress =
        values.street.trim() !== "" ||
        values.city.trim() !== "" ||
        values.state.trim() !== "" ||
        values.zipCode.trim() !== "";
      return !hasAddress || values.city.trim() !== "";
    },
    { message: "validation.required", path: ["city"] }
  )
  .refine(
    (values) => {
      const hasAddress =
        values.street.trim() !== "" ||
        values.city.trim() !== "" ||
        values.state.trim() !== "" ||
        values.zipCode.trim() !== "";
      return !hasAddress || values.state.trim() !== "";
    },
    { message: "validation.required", path: ["state"] }
  )
  .refine(
    (values) => {
      const hasAddress =
        values.street.trim() !== "" ||
        values.city.trim() !== "" ||
        values.state.trim() !== "" ||
        values.zipCode.trim() !== "";

      const hasPhoneNumber = values.phoneNumber.trim() !== "";
      const hasPhoneCode = values.phoneCode.trim() !== "";
      return !hasAddress || (hasPhoneNumber && hasPhoneCode);
    },
    { message: "validation.address_requires_phone", path: ["phoneNumber"] }
  );

export type AddCustomerFormValues = z.infer<typeof addCustomerSchema>;

function getDefaultValues(businessCountryCode: string): AddCustomerFormValues {
  return {
    name: "",
    email: "",
    gender: "other",
    countryCode: businessCountryCode,
    phoneCode: "",
    phoneNumber: "",
    street: "",
    city: "",
    state: "",
    zipCode: "",
  };
}

export function AddCustomerSheet({
  isOpen,
  onClose,
  businessDescriptor,
  businessCountryCode,
  onCreated,
}: AddCustomerSheetProps) {
  const { t, i18n } = useTranslation();
  const { t: tErrors } = useTranslation("errors");
  const formId = useId();

  const countries = useMetadataStore((s) => s.countries);
  const countriesStatus = useMetadataStore((s) => s.status);
  const loadCountries = useMetadataStore((s) => s.loadCountries);

  const countriesReady = countries.length > 0 || countriesStatus === "loaded";

  const isArabic = i18n.language.toLowerCase().startsWith("ar");

  useEffect(() => {
    if (!isOpen) return;
    void loadCountries();
  }, [isOpen, loadCountries]);

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
      return { value: c.code, label };
    });
  }, [countries, isArabic]);

  const phoneCodeOptions = useMemo(() => {
    const seen = new Set<string>();
    const options: { value: string; label: string }[] = [];

    for (const c of countries) {
      if (!c.phonePrefix) continue;
      if (seen.has(c.phonePrefix)) continue;
      seen.add(c.phonePrefix);

      const countryLabel = isArabic ? c.nameAr : c.name;
      const label = `\u200E${c.phonePrefix} — ${countryLabel}`;
      options.push({ value: c.phonePrefix, label });
    }

    return options;
  }, [countries, isArabic]);

  const {
    register,
    control,
    handleSubmit,
    reset,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<AddCustomerFormValues>({
    resolver: zodResolver(addCustomerSchema),
    defaultValues: getDefaultValues(businessCountryCode),
    shouldFocusError: true,
    mode: "onBlur",
  });

  const selectedCountryCode = useWatch({ control, name: "countryCode" });

  useEffect(() => {
    if (!isOpen) return;
    if (!countriesReady) return;
    const selected = countryByCode.get(selectedCountryCode);
    if (!selected?.phonePrefix) return;
    setValue("phoneCode", selected.phonePrefix, { shouldValidate: true });
  }, [isOpen, countriesReady, countryByCode, selectedCountryCode, setValue]);

  useEffect(() => {
    if (!isOpen) {
      reset(getDefaultValues(businessCountryCode));
    }
  }, [isOpen, reset, businessCountryCode]);

  const safeClose = () => {
    if (isSubmitting) return;
    onClose();
  };

  const onSubmit = handleSubmit(async (values) => {
    try {
      const trimOrUndefined = (value: string) => {
        const trimmed = value.trim();
        return trimmed === "" ? undefined : trimmed;
      };

      const phoneCode = values.phoneCode.trim();
      const phoneNumber = values.phoneNumber.trim();

      const normalizedPhone =
        phoneNumber !== "" && phoneCode !== ""
          ? buildE164Phone(phoneCode, phoneNumber)
          : undefined;

      const created = await createCustomer(businessDescriptor, {
        name: values.name.trim(),
        email: values.email.trim(),
        gender: values.gender as CustomerGender,
        countryCode: values.countryCode.trim().toUpperCase(),
        phoneCode: normalizedPhone ? normalizedPhone.phoneCode : undefined,
        phoneNumber: normalizedPhone ? normalizedPhone.phoneNumber : undefined,
      });

      const hasAddress =
        Boolean(values.street.trim()) ||
        Boolean(values.city.trim()) ||
        Boolean(values.state.trim()) ||
        Boolean(values.zipCode.trim());
      if (hasAddress && normalizedPhone) {
        await (
          createCustomerAddress as unknown as (
            businessDescriptor: string,
            customerId: string,
            data: CreateCustomerAddressRequest
          ) => Promise<CustomerAddress>
        )(businessDescriptor, created.id, {
          street: trimOrUndefined(values.street),
          city: values.city.trim(),
          state: values.state.trim(),
          zipCode: trimOrUndefined(values.zipCode),
          countryCode: values.countryCode.trim().toUpperCase(),
          phoneCode: normalizedPhone.phoneCode,
          phone: normalizedPhone.e164,
        });
      }

      toast.success(t("customers.create_success"));

      if (onCreated) {
        await onCreated(created);
      }

      onClose();
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    }
  });

  const footer = (
    <div className="flex gap-2">
      <button
        type="button"
        className="btn btn-ghost flex-1"
        onClick={safeClose}
        disabled={isSubmitting}
        aria-disabled={isSubmitting}
      >
        {t("common.cancel")}
      </button>
      <button
        type="submit"
        form={`add-customer-form-${formId}`}
        className="btn btn-primary flex-1"
        disabled={isSubmitting}
        aria-disabled={isSubmitting}
      >
        {isSubmitting ? t("customers.create_submitting") : t("customers.create_submit")}
      </button>
    </div>
  );

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={safeClose}
      title={t("customers.create_title")}
      footer={footer}
      side="end"
      size="md"
      closeOnOverlayClick={!isSubmitting}
      closeOnEscape={!isSubmitting}
      contentClassName="space-y-4"
      ariaLabel={t("customers.create_title")}
    >
      <form
        id={`add-customer-form-${formId}`}
        onSubmit={(e) => {
          void onSubmit(e);
        }}
        className="space-y-4"
        aria-busy={isSubmitting}
      >
        <FormInput
          label={t("customers.form.name")}
          placeholder={t("customers.form.name_placeholder")}
          autoComplete="name"
          required
          error={errors.name?.message ? tErrors(errors.name.message) : undefined}
          {...register("name")}
        />

        <FormInput
          label={t("customers.form.email")}
          type="email"
          placeholder={t("customers.form.email_placeholder")}
          autoComplete="email"
          inputMode="email"
          required
          error={errors.email?.message ? tErrors(errors.email.message) : undefined}
          {...register("email")}
        />

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <Controller
            control={control}
            name="countryCode"
            render={({ field }) => (
              <FormSelect<string>
                label={t("customers.form.country")}
                options={countryOptions}
                value={field.value}
                onChange={(value) => {
                  field.onChange(value as string);
                }}
                required
                disabled={isSubmitting || !countriesReady}
                placeholder={t("customers.form.select_country")}
                searchable
                error={errors.countryCode?.message ? tErrors(errors.countryCode.message) : undefined}
              />
            )}
          />

          <Controller
            control={control}
            name="gender"
            render={({ field }) => (
              <FormSelect<string>
                label={t("customers.form.gender")}
                options={[
                  { value: "male", label: t("customers.form.gender_male") },
                  { value: "female", label: t("customers.form.gender_female") },
                  { value: "other", label: t("customers.form.gender_other") },
                ]}
                value={field.value}
                onChange={(value) => {
                  field.onChange(value as string);
                }}
                required
                disabled={isSubmitting}
                placeholder={t("customers.form.select_gender")}
                error={errors.gender?.message ? tErrors(errors.gender.message) : undefined}
              />
            )}
          />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <Controller
            control={control}
            name="phoneCode"
            render={({ field }) => (
              <FormSelect<string>
                label={t("customers.form.phone_code")}
                options={phoneCodeOptions}
                value={field.value}
                onChange={(value) => {
                  field.onChange(value as string);
                }}
                disabled={isSubmitting || !countriesReady}
                placeholder={t("customers.form.select_phone_code")}
                searchable
                error={errors.phoneCode?.message ? tErrors(errors.phoneCode.message) : undefined}
              />
            )}
          />

          <div className="sm:col-span-2">
            <FormInput
              label={t("customers.form.phone_number")}
              placeholder={t("customers.form.phone_placeholder")}
              autoComplete="tel"
              inputMode="tel"
              dir="ltr"
              error={errors.phoneNumber?.message ? tErrors(errors.phoneNumber.message) : undefined}
              {...register("phoneNumber")}
            />
          </div>
        </div>

        <div className="divider my-2">{t("customers.form.address_section")}</div>

        <FormTextarea
          label={t("customers.form.street")}
          placeholder={t("customers.form.street_placeholder")}
          rows={2}
          autoComplete="street-address"
          error={errors.street?.message ? tErrors(errors.street.message) : undefined}
          {...register("street")}
        />

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <FormInput
            label={t("customers.form.city")}
            autoComplete="address-level2"
            error={errors.city?.message ? tErrors(errors.city.message) : undefined}
            {...register("city")}
          />

          <FormInput
            label={t("customers.form.state")}
            autoComplete="address-level1"
            error={errors.state?.message ? tErrors(errors.state.message) : undefined}
            {...register("state")}
          />
        </div>

        <FormInput
          label={t("customers.form.zip")}
          autoComplete="postal-code"
          error={errors.zipCode?.message ? tErrors(errors.zipCode.message) : undefined}
          {...register("zipCode")}
        />
      </form>
    </BottomSheet>
  );
}
