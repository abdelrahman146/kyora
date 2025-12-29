import { useEffect, useId } from "react";
import { useTranslation } from "react-i18next";
import toast from "react-hot-toast";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { BottomSheet } from "../../molecules/BottomSheet";
import { CountrySelect } from "../../molecules/CountrySelect";
import { PhoneCodeSelect } from "../../molecules/PhoneCodeSelect";
import { FormInput, FormSelect } from "@/components";
import { updateCustomer } from "@/api/customer";
import type { CustomerGender, Customer } from "@/api/types/customer";
import { translateErrorAsync } from "@/lib/translateError";
import { buildE164Phone } from "@/lib/phone";

export interface EditCustomerSheetProps {
  isOpen: boolean;
  onClose: () => void;
  businessDescriptor: string;
  customer: Customer;
  onUpdated?: (customer: Customer) => void | Promise<void>;
}

const editCustomerSchema = z
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
  })
  .refine(
    (values) => {
      const hasPhoneNumber = values.phoneNumber.trim() !== "";
      return !hasPhoneNumber || values.phoneCode.trim() !== "";
    },
    { message: "validation.required", path: ["phoneCode"] }
  );

export type EditCustomerFormValues = z.infer<typeof editCustomerSchema>;

function getDefaultValues(customer: Customer): EditCustomerFormValues {
  return {
    name: customer.name,
    email: customer.email ?? "",
    gender: customer.gender,
    countryCode: customer.countryCode,
    phoneCode: customer.phoneCode ?? "",
    phoneNumber: customer.phoneNumber ?? "",
  };
}

export function EditCustomerSheet({
  isOpen,
  onClose,
  businessDescriptor,
  customer,
  onUpdated,
}: EditCustomerSheetProps) {
  const { t } = useTranslation();
  const { t: tErrors } = useTranslation("errors");
  const formId = useId();

  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting, isDirty },
  } = useForm<EditCustomerFormValues>({
    resolver: zodResolver(editCustomerSchema),
    defaultValues: getDefaultValues(customer),
    shouldFocusError: true,
    mode: "onBlur",
  });

  useEffect(() => {
    if (isOpen) {
      reset(getDefaultValues(customer));
    }
  }, [isOpen, reset, customer]);

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

      const updated = await updateCustomer(businessDescriptor, customer.id, {
        name: values.name.trim(),
        email: values.email.trim(),
        gender: values.gender as CustomerGender,
        countryCode: values.countryCode.trim().toUpperCase(),
        phoneCode: normalizedPhone ? normalizedPhone.phoneCode : undefined,
        phoneNumber: normalizedPhone ? normalizedPhone.phoneNumber : undefined,
      });

      toast.success(t("customers.update_success"));

      if (onUpdated) {
        await onUpdated(updated);
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
        form={`edit-customer-form-${formId}`}
        className="btn btn-primary flex-1"
        disabled={isSubmitting || !isDirty}
        aria-disabled={isSubmitting || !isDirty}
      >
        {isSubmitting ? t("customers.update_submitting") : t("customers.update_submit")}
      </button>
    </div>
  );

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={safeClose}
      title={t("customers.edit_title")}
      footer={footer}
      side="end"
      size="md"
      closeOnOverlayClick={!isSubmitting}
      closeOnEscape={!isSubmitting}
      contentClassName="space-y-4"
      ariaLabel={t("customers.edit_title")}
    >
      <form
        id={`edit-customer-form-${formId}`}
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
      </form>
    </BottomSheet>
  );
}
