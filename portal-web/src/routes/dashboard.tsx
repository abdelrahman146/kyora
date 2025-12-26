import { useAuth } from "../hooks/useAuth";
import { useTranslation } from "react-i18next";
import { DashboardLayout } from "../components/templates";
import { useBusinessStore } from "../stores/businessStore";

/**
 * Dashboard Page
 *
 * Main dashboard showing business overview and key metrics.
 */
export default function DashboardPage() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const { selectedBusiness } = useBusinessStore();

  return (
    <DashboardLayout title={t("dashboard.title")}>
      <div className="space-y-6">
        {/* Welcome Section */}
        <div className="card bg-base-200 shadow-sm">
          <div className="card-body">
            <h2 className="card-title text-2xl">
              {t("dashboard.welcome")}, {user?.firstName}!
            </h2>
            <p className="text-base-content/70">
              {selectedBusiness
                ? `${t("dashboard.managing")}: ${selectedBusiness.name}`
                : t("dashboard.select_business_to_start")}
            </p>
          </div>
        </div>

        {/* Placeholder Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div className="card bg-base-100 shadow-sm border border-base-300">
            <div className="card-body">
              <h3 className="card-title text-lg">{t("dashboard.revenue")}</h3>
              <p className="text-3xl font-bold text-primary">AED 0</p>
              <p className="text-sm text-base-content/60">
                {t("dashboard.this_month")}
              </p>
            </div>
          </div>

          <div className="card bg-base-100 shadow-sm border border-base-300">
            <div className="card-body">
              <h3 className="card-title text-lg">{t("dashboard.orders")}</h3>
              <p className="text-3xl font-bold text-success">0</p>
              <p className="text-sm text-base-content/60">
                {t("dashboard.pending")}
              </p>
            </div>
          </div>

          <div className="card bg-base-100 shadow-sm border border-base-300">
            <div className="card-body">
              <h3 className="card-title text-lg">{t("dashboard.inventory")}</h3>
              <p className="text-3xl font-bold text-warning">0</p>
              <p className="text-sm text-base-content/60">
                {t("dashboard.low_stock")}
              </p>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
