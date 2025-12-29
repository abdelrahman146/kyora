/**
 * Customer Details Page
 *
 * Displays complete customer information including:
 * - Basic details (name, email, phone, gender, country)
 * - Addresses
 * - Notes
 * - Order statistics
 *
 * Actions:
 * - Edit customer
 * - Delete customer
 * - Add/edit addresses
 * - Add/edit notes
 *
 * Mobile-first with responsive design and RTL support
 */

import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  ArrowLeft,
  Edit,
  Trash2,
  Mail,
  Phone,
  MapPin,
  Globe,
  User,
  ShoppingBag,
} from "lucide-react";
import { DashboardLayout } from "../../components/templates";
import { Avatar } from "../../components/atoms/Avatar";
import { Dialog } from "../../components/atoms/Dialog";
import { EditCustomerSheet } from "../../components/organisms/customers/EditCustomerSheet";
import { useBusinessStore } from "../../stores/businessStore";
import { useMetadataStore } from "../../stores/metadataStore";
import { getCustomer, deleteCustomer } from "../../api/customer";
import type { Customer } from "../../api/types/customer";
import toast from "react-hot-toast";
import { translateErrorAsync } from "@/lib/translateError";

export default function CustomerDetailsPage() {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const { businessDescriptor: urlBusinessDescriptor, customerId } = useParams<{
    businessDescriptor: string;
    customerId: string;
  }>();
  const { selectedBusiness, setSelectedBusinessId, businesses } = useBusinessStore();
  const { countries } = useMetadataStore();

  const isArabic = i18n.language.toLowerCase().startsWith("ar");

  // Sync URL business descriptor with state on mount
  useEffect(() => {
    if (urlBusinessDescriptor && businesses.length > 0) {
      const business = businesses.find((b) => b.descriptor === urlBusinessDescriptor);
      if (business && business.id !== selectedBusiness?.id) {
        setSelectedBusinessId(business.id);
      }
    }
  }, [urlBusinessDescriptor, businesses, selectedBusiness, setSelectedBusinessId]);

  const businessDescriptor = selectedBusiness?.descriptor ?? urlBusinessDescriptor;

  // State
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  // Fetch customer details
  useEffect(() => {
    const fetchCustomer = async () => {
      if (!businessDescriptor || !customerId) return;

      try {
        setIsLoading(true);
        const data = await getCustomer(businessDescriptor, customerId);
        setCustomer(data);
      } catch (error) {
        console.error("Failed to fetch customer:", error);
        const message = await translateErrorAsync(error, t);
        toast.error(message);
        void navigate(`/businesses/${businessDescriptor}/customers`);
      } finally {
        setIsLoading(false);
      }
    };

    void fetchCustomer();
  }, [businessDescriptor, customerId, t, navigate]);

  // Handlers
  const handleBack = () => {
    if (!businessDescriptor) return;
    void navigate(`/businesses/${businessDescriptor}/customers`);
  };

  const handleDelete = async () => {
    if (!businessDescriptor || !customerId) return;

    try {
      setIsDeleting(true);
      await deleteCustomer(businessDescriptor, customerId);
      toast.success(t("customers.delete_success"));
      void navigate(`/businesses/${businessDescriptor}/customers`);
    } catch (error) {
      const message = await translateErrorAsync(error, t);
      toast.error(message);
    } finally {
      setIsDeleting(false);
      setIsDeleteDialogOpen(false);
    }
  };

  const handleCustomerUpdated = (updated: Customer) => {
    setCustomer(updated);
  };

  const getInitials = (name: string): string => {
    return name
      .split(" ")
      .map((w) => w[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  const formatPhone = (): string | null => {
    if (customer?.phoneCode && customer.phoneNumber) {
      return `${customer.phoneCode} ${customer.phoneNumber}`;
    }
    return null;
  };

  const getCountryInfo = (countryCode: string) => {
    const country = countries.find((c) => c.code === countryCode);
    return {
      name: isArabic ? (country?.nameAr ?? countryCode) : (country?.name ?? countryCode),
      flag: country?.flag,
    };
  };

  const getGenderLabel = (gender: string): string => {
    const genderMap: Record<string, string> = {
      male: t("customers.form.gender_male"),
      female: t("customers.form.gender_female"),
      other: t("customers.form.gender_other"),
    };
    return genderMap[gender] || gender;
  };

  // Guard: No business descriptor in URL
  if (!businessDescriptor || !customerId) {
    return (
      <DashboardLayout title={t("customers.details_title")}>
        <div className="alert alert-warning">
          <span>{t("dashboard.select_business_to_start")}</span>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout title={customer?.name ?? t("customers.details_title")}>
      <div className="space-y-4">
        {/* Header */}
        <div className="flex items-center gap-4">
          <button
            type="button"
            className="btn btn-ghost btn-sm gap-2"
            onClick={handleBack}
            aria-label={t("common.back")}
          >
            <ArrowLeft size={18} className={isArabic ? "rotate-180" : ""} />
            <span className="hidden sm:inline">{t("common.back")}</span>
          </button>
          <h1 className="text-2xl font-bold flex-1 truncate">
            {isLoading ? (
              <div className="skeleton h-8 w-48"></div>
            ) : (
              customer?.name ?? t("customers.details_title")
            )}
          </h1>
          {!isLoading && customer && (
            <div className="flex gap-2">
              <button
                type="button"
                className="btn btn-outline btn-sm gap-2"
                onClick={() => {
                  setIsEditOpen(true);
                }}
              >
                <Edit size={16} />
                <span className="hidden sm:inline">{t("common.edit")}</span>
              </button>
              <button
                type="button"
                className="btn btn-error btn-outline btn-sm gap-2"
                onClick={() => {
                  setIsDeleteDialogOpen(true);
                }}
              >
                <Trash2 size={16} />
                <span className="hidden sm:inline">{t("common.delete")}</span>
              </button>
            </div>
          )}
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="space-y-4">
            <div className="card bg-base-100 border border-base-300 shadow-sm">
              <div className="card-body">
                <div className="flex items-center gap-4">
                  <div className="skeleton w-24 h-24 rounded-full"></div>
                  <div className="flex-1 space-y-2">
                    <div className="skeleton h-6 w-48"></div>
                    <div className="skeleton h-4 w-32"></div>
                  </div>
                </div>
              </div>
            </div>
            <div className="skeleton h-64 rounded-box"></div>
          </div>
        )}

        {/* Customer Details */}
        {!isLoading && customer && (
          <>
            {/* Header Card with Avatar */}
            <div className="card bg-base-100 border border-base-300 shadow-sm">
              <div className="card-body">
                <div className="flex flex-col sm:flex-row items-center sm:items-start gap-4">
                  <Avatar
                    src={customer.avatarUrl}
                    alt={customer.name}
                    fallback={getInitials(customer.name)}
                    size="xl"
                  />
                  <div className="flex-1 text-center sm:text-start">
                    <h2 className="text-2xl font-bold">{customer.name}</h2>
                    {customer.email && (
                      <div className="flex items-center gap-2 justify-center sm:justify-start mt-2 text-base-content/70">
                        <Mail size={16} />
                        <span>{customer.email}</span>
                      </div>
                    )}
                    {formatPhone() && (
                      <div className="flex items-center gap-2 justify-center sm:justify-start mt-1 text-base-content/70">
                        <Phone size={16} />
                        <span dir="ltr">{formatPhone()}</span>
                      </div>
                    )}
                  </div>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-2 gap-4 mt-6">
                  <div className="stat bg-base-200 rounded-box p-4">
                    <div className="stat-figure text-success">
                      <ShoppingBag size={28} />
                    </div>
                    <div className="stat-title text-sm">{t("customers.orders_count")}</div>
                    <div className="text-2xl font-bold text-success">{customer.ordersCount ?? 0}</div>
                  </div>
                  <div className="stat bg-base-200 rounded-box p-4">
                    <div className="stat-title text-sm">{t("customers.total_spent")}</div>
                    <div className="text-2xl font-bold text-primary">
                      {selectedBusiness?.currency ?? "AED"}{" "}
                      {(customer.totalSpent ?? 0).toLocaleString(isArabic ? "ar-AE" : "en-US", {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Details Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Basic Information */}
              <div className="card bg-base-100 border border-base-300 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title text-lg">{t("customers.details.basic_info")}</h3>
                  <div className="space-y-3 mt-4">
                    <div className="flex items-start gap-3">
                      <User size={18} className="text-base-content/60 mt-0.5" />
                      <div className="flex-1">
                        <div className="text-xs text-base-content/60">
                          {t("customers.form.gender")}
                        </div>
                        <div className="font-medium">{getGenderLabel(customer.gender)}</div>
                      </div>
                    </div>
                    <div className="flex items-start gap-3">
                      <Globe size={18} className="text-base-content/60 mt-0.5" />
                      <div className="flex-1">
                        <div className="text-xs text-base-content/60">
                          {t("customers.form.country")}
                        </div>
                        <div className="flex items-center gap-2 font-medium">
                          {getCountryInfo(customer.countryCode).flag && (
                            <span className="text-lg">{getCountryInfo(customer.countryCode).flag}</span>
                          )}
                          <span>{getCountryInfo(customer.countryCode).name}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Addresses */}
              <div className="card bg-base-100 border border-base-300 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title text-lg">{t("customers.details.addresses")}</h3>
                  <div className="space-y-3 mt-4">
                    {customer.addresses && customer.addresses.length > 0 ? (
                      customer.addresses.map((address) => (
                        <div key={address.id} className="flex items-start gap-3 p-3 bg-base-200 rounded-lg">
                          <MapPin size={18} className="text-base-content/60 mt-0.5" />
                          <div className="flex-1">
                            <div className="font-medium">
                              {address.city}, {address.state}
                            </div>
                            {address.street && (
                              <div className="text-sm text-base-content/70">{address.street}</div>
                            )}
                            {address.zipCode && (
                              <div className="text-sm text-base-content/70">
                                {address.zipCode}
                              </div>
                            )}
                            <div className="text-sm text-base-content/70 mt-1" dir="ltr">
                              {address.phoneCode} {address.phoneNumber}
                            </div>
                          </div>
                        </div>
                      ))
                    ) : (
                      <div className="text-center py-8 text-base-content/60">
                        <MapPin size={32} className="mx-auto mb-2 opacity-40" />
                        <p>{t("customers.details.no_addresses")}</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>

            {/* Notes Section */}
            {customer.notes && customer.notes.length > 0 && (
              <div className="card bg-base-100 border border-base-300 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title text-lg">{t("customers.details.notes")}</h3>
                  <div className="space-y-2 mt-4">
                    {customer.notes.map((note) => (
                      <div key={note.id} className="p-3 bg-base-200 rounded-lg">
                        <p className="text-sm">{note.content}</p>
                        <div className="text-xs text-base-content/60 mt-2">
                          {new Date(note.createdAt).toLocaleDateString(
                            isArabic ? "ar-AE" : "en-US",
                            {
                              year: "numeric",
                              month: "long",
                              day: "numeric",
                            }
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Edit Customer Sheet */}
      {customer && (
        <EditCustomerSheet
          isOpen={isEditOpen}
          onClose={() => {
            setIsEditOpen(false);
          }}
          businessDescriptor={businessDescriptor}
          customer={customer}
          onUpdated={handleCustomerUpdated}
        />
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={isDeleteDialogOpen}
        onClose={() => {
          setIsDeleteDialogOpen(false);
        }}
        title={t("customers.delete_confirm_title")}
        size="sm"
        footer={
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn btn-ghost"
              onClick={() => {
                setIsDeleteDialogOpen(false);
              }}
              disabled={isDeleting}
            >
              {t("common.cancel")}
            </button>
            <button
              type="button"
              className="btn btn-error"
              onClick={() => {
                void handleDelete();
              }}
              disabled={isDeleting}
            >
              {isDeleting && <span className="loading loading-spinner loading-sm" />}
              {t("common.delete")}
            </button>
          </div>
        }
      >
        <p>{t("customers.delete_confirm_message", { name: customer?.name })}</p>
      </Dialog>
    </DashboardLayout>
  );
}
