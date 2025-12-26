import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  User,
  Settings,
  LogOut,
  ChevronDown,
  HelpCircle,
} from "lucide-react";
import { Avatar } from "../atoms/Avatar";
import { useAuth } from "../../hooks/useAuth";
import { cn } from "@/lib/utils";

/**
 * UserMenu Component
 *
 * Dropdown menu in the header showing user info and actions.
 *
 * Features:
 * - User avatar and name
 * - Profile, Settings, Help links
 * - Logout action
 * - RTL support
 *
 * @example
 * ```tsx
 * <UserMenu />
 * ```
 */
export function UserMenu() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [isOpen, setIsOpen] = useState(false);

  const handleLogout = () => {
    void logout()
      .then(() => {
        setIsOpen(false);
        void navigate("/login", { replace: true });
      })
      .catch((error: unknown) => {
        console.error("Logout failed:", error);
        setIsOpen(false);
        void navigate("/login", { replace: true });
      });
  };

  if (!user) return null;

  const userFullName = `${user.firstName} ${user.lastName}`;
  const userInitials = `${user.firstName[0]}${user.lastName[0]}`;

  return (
    <div className="dropdown dropdown-end">
      {/* Trigger Button */}
      <button
        type="button"
        tabIndex={0}
        onClick={() => { setIsOpen(!isOpen); }}
        className="btn btn-ghost gap-2 h-auto min-h-0 py-2 px-2 hover:bg-base-200"
      >
        <Avatar
          src={undefined}
          fallback={userInitials}
          size="sm"
          shape="circle"
        />
        <ChevronDown
          size={16}
          className={cn(
            "transition-transform text-base-content/60 hidden md:block",
            isOpen && "rotate-180"
          )}
        />
      </button>

      {/* Dropdown Menu */}
      {isOpen && (
        <ul
          tabIndex={0}
          className="dropdown-content menu p-2 shadow-lg bg-base-100 rounded-lg w-56 mt-2 border border-base-300"
        >
          {/* User Info */}
          <li className="px-4 py-3 border-b border-base-300">
            <div className="flex flex-col gap-1">
              <span className="font-semibold text-sm text-base-content">
                {userFullName}
              </span>
              <span className="text-xs text-base-content/60">
                {user.email}
              </span>
            </div>
          </li>

          {/* Menu Items */}
          <li>
            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                void navigate("/profile");
              }}
              className="flex items-center gap-3 px-3 py-2 text-start"
            >
              <User size={16} />
              <span className="text-sm">{t("dashboard.profile")}</span>
            </button>
          </li>

          <li>
            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                void navigate("/settings");
              }}
              className="flex items-center gap-3 px-3 py-2 text-start"
            >
              <Settings size={16} />
              <span className="text-sm">{t("dashboard.settings")}</span>
            </button>
          </li>

          <li>
            <button
              type="button"
              onClick={() => {
                setIsOpen(false);
                void navigate("/help");
              }}
              className="flex items-center gap-3 px-3 py-2 text-start"
            >
              <HelpCircle size={16} />
              <span className="text-sm">{t("dashboard.help")}</span>
            </button>
          </li>

          <div className="divider my-1"></div>

          {/* Logout */}
          <li>
            <button
              type="button"
              onClick={handleLogout}
              className="flex items-center gap-3 px-3 py-2 text-start text-error hover:bg-error/10"
            >
              <LogOut size={16} />
              <span className="text-sm">{t("auth.logout")}</span>
            </button>
          </li>
        </ul>
      )}
    </div>
  );
}
