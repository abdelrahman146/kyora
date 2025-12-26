import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useLanguage } from "../hooks/useLanguage";

/**
 * Dashboard Page (Placeholder)
 *
 * This is a temporary dashboard for testing authentication flow.
 * Will be replaced with actual dashboard implementation.
 */
export default function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { language, toggleLanguage } = useLanguage();

  const handleLogout = () => {
    void logout()
      .then(() => {
        void navigate("/login");
      })
      .catch((error: unknown) => {
        console.error("Logout failed:", error);
        void navigate("/login");
      });
  };

  return (
    <div className="min-h-screen bg-base-100">
      <div className="container mx-auto p-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-4xl font-bold text-primary">
            {t("dashboard.title")}
          </h1>
          <div className="flex gap-4">
            <button onClick={toggleLanguage} className="btn btn-ghost">
              {language === "ar" ? "English" : "العربية"}
            </button>
            <button onClick={handleLogout} className="btn btn-error">
              Logout
            </button>
          </div>
        </div>

        <div className="card bg-base-200 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t("dashboard.welcome")}!</h2>
            <div className="space-y-2">
              <p>
                <strong>Name:</strong> {user?.firstName} {user?.lastName}
              </p>
              <p>
                <strong>Email:</strong> {user?.email}
              </p>
              <p>
                <strong>Role:</strong>{" "}
                <span className="badge badge-primary">{user?.role}</span>
              </p>
              <p>
                <strong>Workspace ID:</strong> {user?.workspaceId}
              </p>
            </div>
          </div>
        </div>

        <div className="alert alert-success mt-8">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>Authentication successful! You are now logged in.</span>
        </div>
      </div>
    </div>
  );
}
