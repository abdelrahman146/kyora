import { useEffect, useId } from "react";
import { useTranslation } from "react-i18next";
import toast from "react-hot-toast";
import { Controller, useForm, useWatch } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { BottomSheet } from "../../molecules/BottomSheet";
import { CountrySelect } from "../../molecules/CountrySelect";
import { PhoneCodeSelect } from "../../molecules/PhoneCodeSelect";
import { SocialMediaInputs } from "../../molecules/SocialMediaInputs";
import { FormInput, FormSelect } from "@/components";
import { createCustomer } from "@/api/customer";
import type {
  CustomerGender,
  Customer,
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
    instagramUsername: z.string().trim().optional(),
    facebookUsername: z.string().trim().optional(),
    tiktokUsername: z.string().trim().optional(),
    snapchatUsername: z.string().trim().optional(),
    xUsername: z.string().trim().optional(),
    whatsappNumber: z.string().trim().optional(),
  })
  .refine(
    (values) => {
      const hasPhoneNumber = values.phoneNumber.trim() !== "";
      return !hasPhoneNumber || values.phoneCode.trim() !== "";
    },
    { message: "validation.required", path: ["phoneCode"] }
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
    instagramUsername: "",
    facebookUsername: "",
    tiktokUsername: "",
    snapchatUsername: "",
    xUsername: "",
    whatsappNumber: "",
  };
}

export function AddCustomerSheet({
  isOpen,
  onClose,
  businessDescriptor,
  businessCountryCode,
  onCreated,
}: AddCustomerSheetProps) {
  const { t } = useTranslation();
  const { t: tErrors } = useTranslation("errors");
  const formId = useId();

  const countries = useMetadataStore((s) => s.countries);
  const countriesStatus = useMetadataStore((s) => s.status);
  const loadCountries = useMetadataStore((s) => s.loadCountries);

  const countriesReady = countries.length > 0 || countriesStatus === "loaded";

  useEffect(() => {
    if (!isOpen) return;
    void loadCountries();
  }, [isOpen, loadCountries]);

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
  
  // Watch social media fields for the inputs component
  const socialMediaValues = {
    instagramUsername: useWatch({ control, name: "instagramUsername" }) ?? "",
    facebookUsername: useWatch({ control, name: "facebookUsername" }) ?? "",
    tiktokUsername: useWatch({ control, name: "tiktokUsername" }) ?? "",
    snapchatUsername: useWatch({ control, name: "snapchatUsername" }) ?? "",
    xUsername: useWatch({ control, name: "xUsername" }) ?? "",
    whatsappNumber: useWatch({ control, name: "whatsappNumber" }) ?? "",
  };

  useEffect(() => {
    if (!isOpen) return;
    if (!countriesReady) return;
    const selected = countries.find((c) => c.code === selectedCountryCode);
    if (!selected?.phonePrefix) return;
    setValue("phoneCode", selected.phonePrefix, { shouldValidate: true });
  }, [isOpen, countriesReady, countries, selectedCountryCode, setValue]);

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
        instagramUsername: values.instagramUsername?.trim() !== "" ? values.instagramUsername?.trim() : undefined,
        facebookUsername: values.facebookUsername?.trim() !== "" ? values.facebookUsername?.trim() : undefined,
        tiktokUsername: values.tiktokUsername?.trim() !== "" ? values.tiktokUsername?.trim() : undefined,
        snapchatUsername: values.snapchatUsername?.trim() !== "" ? values.snapchatUsername?.trim() : undefined,
        xUsername: values.xUsername?.trim() !== "" ? values.xUsername?.trim() : undefined,
        whatsappNumber: values.whatsappNumber?.trim() !== "" ? values.whatsappNumber?.trim() : undefined,
      });

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
              <CountrySelect
                value={field.value}
                onChange={field.onChange}
                error={errors.countryCode?.message ? tErrors(errors.countryCode.message) : undefined}
                disabled={isSubmitting}
                required
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
              <PhoneCodeSelect
                value={field.value}
                onChange={field.onChange}
                error={errors.phoneCode?.message ? tErrors(errors.phoneCode.message) : undefined}
                disabled={isSubmitting}
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

        <SocialMediaInputs
          instagramUsername={socialMediaValues.instagramUsername}
          onInstagramChange={(value) => {
            setValue("instagramUsername", value);
          }}
          instagramError={errors.instagramUsername?.message ? tErrors(errors.instagramUsername.message) : undefined}
          facebookUsername={socialMediaValues.facebookUsername}
          onFacebookChange={(value) => {
            setValue("facebookUsername", value);
          }}
          facebookError={errors.facebookUsername?.message ? tErrors(errors.facebookUsername.message) : undefined}
          tiktokUsername={socialMediaValues.tiktokUsername}
          onTiktokChange={(value) => {
            setValue("tiktokUsername", value);
          }}
          tiktokError={errors.tiktokUsername?.message ? tErrors(errors.tiktokUsername.message) : undefined}
          snapchatUsername={socialMediaValues.snapchatUsername}
          onSnapchatChange={(value) => {
            setValue("snapchatUsername", value);
          }}
          snapchatError={errors.snapchatUsername?.message ? tErrors(errors.snapchatUsername.message) : undefined}
          xUsername={socialMediaValues.xUsername}
          onXChange={(value) => {
            setValue("xUsername", value);
          }}
          xError={errors.xUsername?.message ? tErrors(errors.xUsername.message) : undefined}
          whatsappNumber={socialMediaValues.whatsappNumber}
          onWhatsappChange={(value) => {
            setValue("whatsappNumber", value);
          }}
          whatsappError={errors.whatsappNumber?.message ? tErrors(errors.whatsappNumber.message) : undefined}
          disabled={isSubmitting}
          defaultExpanded={false}
        />
      </form>
    </BottomSheet>
  );
}
