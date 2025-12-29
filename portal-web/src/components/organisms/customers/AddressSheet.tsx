/**
 * AddressSheet Component
 *
 * Reusable bottom sheet for adding/editing customer addresses.
 * Handles form validation, submission, and RTL support.
 *
 * Features:
 * - Mobile-first responsive design
 * - Country and phone code selection from metadata
 * - Bilingual support (Arabic/English)
 * - Form validation with Zod
 * - Optimistic UI updates
 */

import { useEffect, useMemo } from "react";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import { BottomSheet } from "../../molecules/BottomSheet";
import { CountrySelect } from "../../molecules/CountrySelect";
import { PhoneCodeSelect } from "../../molecules/PhoneCodeSelect";
import { useMetadataStore } from "../../../stores/metadataStore";
import type { CustomerAddress } from "../../../api/types/customer";
import type { CreateAddressRequest, UpdateAddressRequest } from "../../../api/address";
import { buildE164Phone, parseE164Phone } from "../../../lib/phone";
import toast from "react-hot-toast";
import { translateErrorAsync } from "@/lib/translateError";

interface AddressSheetProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateAddressRequest | UpdateAddressRequest) => Promise<CustomerAddress>;
  address?: CustomerAddress; // If provided, we're editing
  submitLabel?: string;
}

// Zod schema
const addressSchema = z.object({
  countryCode: z.string().length(2, "Country is required"),
  state: z.string().min(1, "State is required"),
  city: z.string().min(1, "City is required"),
  phoneCode: z.string().min(1, "Phone code is required"),
  phoneNumber: z.string().min(1, "Phone number is required"),
  street: z.string().optional(),
  zipCode: z.string().optional(),
});

type FormData = z.infer<typeof addressSchema>;

export function AddressSheet({
  isOpen,
  onClose,
  onSubmit,
  address,
  submitLabel,
}: AddressSheetProps) {
  const { t } = useTranslation();
  const { countries, status, loadCountries } = useMetadataStore();

  // Check if countries are ready (loaded and have data)
  const countriesReady = useMemo(
    () => countries.length > 0 && status === "loaded",
    [countries.length, status]
  );

  // Fetch countries if not loaded
  useEffect(() => {
    if (status === "idle") {
      void loadCountries();
    }
  }, [status, loadCountries]);

  // Parse address phone if editing
  const initialPhoneData = useMemo(() => {
    if (address) {
      return parseE164Phone(address.phoneCode, address.phoneNumber);
    }
    return { phoneCode: "", phoneNumber: "" };
  }, [address]);

  // Form setup
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting, isDirty },
    reset,
    watch,
    setValue,
    control,
  } = useForm<FormData>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      countryCode: address?.countryCode ?? "",
      state: address?.state ?? "",
      city: address?.city ?? "",
      phoneCode: initialPhoneData.phoneCode,
      phoneNumber: initialPhoneData.phoneNumber,
      street: address?.street ?? "",
      zipCode: address?.zipCode ?? "",
    },
  });

  // Reset form when address changes or sheet opens
  useEffect(() => {
    if (isOpen) {
      reset({
        countryCode: address?.countryCode ?? "",
        state: address?.state ?? "",
        city: address?.city ?? "",
        phoneCode: initialPhoneData.phoneCode,
        phoneNumber: initialPhoneData.phoneNumber,
        street: address?.street ?? "",
        zipCode: address?.zipCode ?? "",
      });
    }
  }, [isOpen, address, initialPhoneData, reset]);

  // Watch country code to auto-set phone code (always auto-update phone code)
  const selectedCountryCode = watch("countryCode");

  useEffect(() => {
    if (selectedCountryCode && countriesReady) {
      const country = countries.find((c) => c.code === selectedCountryCode);
      if (country?.phonePrefix) {
        setValue("phoneCode", country.phonePrefix, { shouldValidate: false });
      }
    }
  }, [selectedCountryCode, countries, countriesReady, setValue]);

  // Handle form submission
  const handleFormSubmit = async (data: FormData) => {
    try {
      // Build E.164 phone
      const phoneData = buildE164Phone(data.phoneCode, data.phoneNumber);

      if (address) {
        // Update
        const updateData: UpdateAddressRequest = {
          countryCode: data.countryCode,
          state: data.state,
          city: data.city,
          phoneCode: data.phoneCode,
          phoneNumber: data.phoneNumber,
          street: data.street,
          zipCode: data.zipCode,
        };
        await onSubmit(updateData);
        toast.success(t("customers.address.update_success"));
      } else {
        // Create
        const createData: CreateAddressRequest = {
          countryCode: data.countryCode,
          state: data.state,
          city: data.city,
          phoneCode: data.phoneCode,
          phone: phoneData.e164, // Backend expects 'phone' field with E.164 format
          street: data.street,
          zipCode: data.zipCode,
        };
        await onSubmit(createData);
        toast.success(t("customers.address.create_success"));
      }
      onClose();
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    }
  };

  return (
    <BottomSheet
      isOpen={isOpen}
      onClose={onClose}
      title={address ? t("customers.address.edit_title") : t("customers.address.add_title")}
    >
      <form
        onSubmit={(e) => {
          void handleSubmit(handleFormSubmit)(e);
        }}
        className="space-y-4"
      >
        {/* Country */}
        <Controller
          name="countryCode"
          control={control}
          render={({ field }) => (
            <CountrySelect
              value={field.value}
              onChange={field.onChange}
              error={errors.countryCode?.message}
              required
            />
          )}
        />

        {/* State */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">
              {t("customers.form.state")} <span className="text-error">*</span>
            </span>
          </label>
          <input
            {...register("state")}
            type="text"
            placeholder={t("customers.form.state_placeholder")}
            className={`input input-bordered w-full ${errors.state ? "input-error" : ""}`}
          />
          {errors.state && (
            <label className="label">
              <span className="label-text-alt text-error">{errors.state.message}</span>
            </label>
          )}
        </div>

        {/* City */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">
              {t("customers.form.city")} <span className="text-error">*</span>
            </span>
          </label>
          <input
            {...register("city")}
            type="text"
            placeholder={t("customers.form.city_placeholder")}
            className={`input input-bordered w-full ${errors.city ? "input-error" : ""}`}
          />
          {errors.city && (
            <label className="label">
              <span className="label-text-alt text-error">{errors.city.message}</span>
            </label>
          )}
        </div>

        {/* Street (Optional) */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t("customers.form.street")}</span>
          </label>
          <input
            {...register("street")}
            type="text"
            placeholder={t("customers.form.street_placeholder")}
            className="input input-bordered w-full"
          />
        </div>

        {/* Zip Code (Optional) */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t("customers.form.zip_code")}</span>
          </label>
          <input
            {...register("zipCode")}
            type="text"
            placeholder={t("customers.form.zip_placeholder")}
            className="input input-bordered w-full"
          />
        </div>

        {/* Phone Code - Auto-updated from country, disabled */}
        <Controller
          name="phoneCode"
          control={control}
          render={({ field }) => (
            <PhoneCodeSelect
              value={field.value}
              onChange={field.onChange}
              error={errors.phoneCode?.message}
              disabled
              required
            />
          )}
        />

        {/* Phone Number */}
        <div className="form-control">
          <label className="label">
            <span className="label-text">
              {t("customers.form.phone_number")} <span className="text-error">*</span>
            </span>
          </label>
          <input
            {...register("phoneNumber")}
            type="tel"
            placeholder={t("customers.form.phone_placeholder")}
            className={`input input-bordered w-full ${errors.phoneNumber ? "input-error" : ""}`}
          />
          {errors.phoneNumber && (
            <label className="label">
              <span className="label-text-alt text-error">{errors.phoneNumber.message}</span>
            </label>
          )}
        </div>

        {/* Footer Actions */}
        <div className="flex gap-2 pt-4">
          <button
            type="button"
            className="btn btn-ghost flex-1"
            onClick={onClose}
            disabled={isSubmitting}
          >
            {t("common.cancel")}
          </button>
          <button
            type="submit"
            className="btn btn-primary flex-1"
            disabled={isSubmitting || (address ? !isDirty : false)}
          >
            {isSubmitting && <span className="loading loading-spinner loading-sm" />}
            {submitLabel ?? (address ? t("common.update") : t("common.add"))}
          </button>
        </div>
      </form>
    </BottomSheet>
  );
}
