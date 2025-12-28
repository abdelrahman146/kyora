/**
 * Home Page (Authenticated)
 *
 * Landing page for authenticated users to select a business or access workspace features.
 *
 * Features:
 * - Business selection cards
 * - Quick links to billing, workspace, account
 * - Support and documentation links
 * - Mobile-first design
 * - Fully localized
 */

import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  Building2,
  CreditCard,
  Settings,
  Users,
  HelpCircle,
  BookOpen,
  ChevronRight,
  Plus,
} from "lucide-react";
import { useAuth } from "../hooks/useAuth";
import { useBusinessStore } from "../stores/businessStore";
import { businessApi } from "../api/business";
import type { Business } from "../api/types/business";
import { Logo } from "../components/atoms/Logo";
import { LanguageSwitcher } from "../components/molecules/LanguageSwitcher";

/**
 * Home Page Component
 */
export default function HomePage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { businesses, setBusinesses, setSelectedBusiness } = useBusinessStore();
  const [isLoading, setIsLoading] = useState(true);

  // Fetch businesses on mount
  useEffect(() => {
    const fetchBusinesses = async () => {
      try {
        setIsLoading(true);
        const businesses = await businessApi.listBusinesses();
        setBusinesses(businesses);
      } catch (error) {
        console.error("Failed to fetch businesses:", error);
      } finally {
        setIsLoading(false);
      }
    };

    void fetchBusinesses();
  }, [setBusinesses]);

  const handleBusinessSelect = (business: Business) => {
    setSelectedBusiness(business);
    void navigate(`/businesses/${business.descriptor}/dashboard`);
  };

  const quickLinks = [
    {
      key: "workspace",
      icon: Users,
      label: t("home.workspace_settings"),
      description: t("home.workspace_settings_desc"),
      path: "/workspace",
      disabled: true,
    },
    {
      key: "account",
      icon: Settings,
      label: t("home.account_settings"),
      description: t("home.account_settings_desc"),
      path: "/account",
      disabled: true,
    },
    {
      key: "billing",
      icon: CreditCard,
      label: t("home.billing"),
      description: t("home.billing_desc"),
      path: "/billing",
      disabled: true,
    },
  ];

  const supportLinks = [
    {
      key: "help",
      icon: HelpCircle,
      label: t("home.help_center"),
      path: "https://help.kyora.app",
      external: true,
    },
    {
      key: "docs",
      icon: BookOpen,
      label: t("home.documentation"),
      path: "https://docs.kyora.app",
      external: true,
    },
  ];

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <header className="sticky top-0 z-30 bg-base-100 border-b border-base-300">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Logo size="md" showText />
            <div className="flex items-center gap-2">
              <LanguageSwitcher variant="compact" />
              <div className="divider divider-horizontal mx-0" />
              <span className="text-sm text-base-content/70 hidden sm:inline">
                {user?.firstName} {user?.lastName}
              </span>
              <button
                className="btn btn-ghost btn-sm"
                onClick={() => {
                  void navigate("/auth/login");
                }}
              >
                {t("auth.logout")}
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8 max-w-6xl">
        {/* Welcome Section */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">
            {t("home.welcome")}, {user?.firstName}!
          </h1>
          <p className="text-base-content/70">{t("home.select_business_or_manage")}</p>
        </div>

        {/* Businesses Section */}
        <section className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">{t("home.your_businesses")}</h2>
            <button className="btn btn-primary btn-sm gap-2" disabled>
              <Plus size={18} />
              {t("home.add_business")}
            </button>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={i} className="skeleton h-32 rounded-box" />
              ))}
            </div>
          ) : businesses.length === 0 ? (
            <div className="card bg-base-100 border border-base-300">
              <div className="card-body items-center text-center py-12">
                <Building2 size={48} className="text-base-content/30 mb-4" />
                <h3 className="text-lg font-semibold mb-2">{t("home.no_businesses")}</h3>
                <p className="text-sm text-base-content/60 mb-4">
                  {t("home.no_businesses_desc")}
                </p>
                <button className="btn btn-primary gap-2" disabled>
                  <Plus size={18} />
                  {t("home.create_first_business")}
                </button>
              </div>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {businesses.map((business) => (
                <button
                  key={business.id}
                  onClick={() => {
                    handleBusinessSelect(business);
                  }}
                  className="card bg-base-100 border border-base-300 hover:border-primary hover:shadow-md transition-all text-start"
                >
                  <div className="card-body p-4">
                    <div className="flex items-start gap-3">
                      {business.logoUrl ? (
                        <img
                          src={business.logoUrl}
                          alt={business.name}
                          className="w-12 h-12 rounded-lg object-cover"
                        />
                      ) : (
                        <div className="avatar placeholder">
                          <div className="bg-primary text-primary-content rounded-lg w-12 h-12">
                            <span className="text-lg font-bold">
                              {business.name.charAt(0).toUpperCase()}
                            </span>
                          </div>
                        </div>
                      )}
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold truncate">{business.name}</h3>
                        <p className="text-sm text-base-content/60 truncate">
                          @{business.descriptor}
                        </p>
                      </div>
                      <ChevronRight size={20} className="text-base-content/40" />
                    </div>
                  </div>
                </button>
              ))}
            </div>
          )}
        </section>

        {/* Quick Links Section */}
        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-4">{t("home.quick_links")}</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {quickLinks.map((link) => (
              <button
                key={link.key}
                onClick={() => {
                  if (!link.disabled) void navigate(link.path);
                }}
                disabled={link.disabled}
                className="card bg-base-100 border border-base-300 hover:border-primary hover:shadow-md transition-all text-start disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <div className="card-body p-4">
                  <div className="flex items-center gap-3 mb-2">
                    <link.icon size={24} className="text-primary" />
                    <h3 className="font-semibold">{link.label}</h3>
                  </div>
                  <p className="text-sm text-base-content/60">{link.description}</p>
                  {link.disabled && (
                    <span className="badge badge-sm badge-ghost mt-2">
                      {t("common.coming_soon")}
                    </span>
                  )}
                </div>
              </button>
            ))}
          </div>
        </section>

        {/* Support Links Section */}
        <section>
          <h2 className="text-xl font-semibold mb-4">{t("home.support")}</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {supportLinks.map((link) => (
              <a
                key={link.key}
                href={link.path}
                target={link.external ? "_blank" : undefined}
                rel={link.external ? "noopener noreferrer" : undefined}
                className="card bg-base-100 border border-base-300 hover:border-primary hover:shadow-md transition-all"
              >
                <div className="card-body p-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <link.icon size={20} className="text-primary" />
                      <span className="font-medium">{link.label}</span>
                    </div>
                    <ChevronRight size={20} className="text-base-content/40" />
                  </div>
                </div>
              </a>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
