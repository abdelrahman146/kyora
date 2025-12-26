import { useTranslation } from "react-i18next";
import { Link, useLocation } from "react-router-dom";
import {
  LayoutDashboard,
  Package,
  ShoppingCart,
  Users,
  BarChart3,
  Calculator,
  CreditCard,
  UsersRound,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { Logo } from "../atoms/Logo";
import { IconButton } from "../atoms/IconButton";
import { useBusinessStore } from "../../stores/businessStore";
import { useLanguage } from "../../hooks/useLanguage";
import { cn } from "@/lib/utils";

interface NavItem {
  key: string;
  icon: typeof LayoutDashboard;
  path: string;
}

const navItems: NavItem[] = [
  { key: "dashboard", icon: LayoutDashboard, path: "/dashboard" },
  { key: "inventory", icon: Package, path: "/inventory" },
  { key: "orders", icon: ShoppingCart, path: "/orders" },
  { key: "customers", icon: Users, path: "/customers" },
  { key: "analytics", icon: BarChart3, path: "/analytics" },
  { key: "accounting", icon: Calculator, path: "/accounting" },
  { key: "billing", icon: CreditCard, path: "/billing" },
  { key: "team", icon: UsersRound, path: "/team" },
];

/**
 * Sidebar Component
 *
 * Desktop vertical navigation sidebar with collapse functionality.
 *
 * Features:
 * - Collapsible (icon-only mode)
 * - Active route highlighting
 * - RTL support
 * - Smooth transitions
 */
export function Sidebar() {
  const { t } = useTranslation();
  const location = useLocation();
  const { isRTL } = useLanguage();
  const { isSidebarCollapsed, toggleSidebar } = useBusinessStore();

  return (
    <aside
      className={cn(
        "fixed start-0 top-0 h-screen bg-base-200 border-e border-base-300 transition-all duration-300 z-40",
        isSidebarCollapsed ? "w-20" : "w-64"
      )}
    >
      {/* Header */}
      <div className="h-16 flex items-center justify-between px-4 border-b border-base-300">
        {!isSidebarCollapsed && <Logo size="md" showText />}
        {isSidebarCollapsed && <Logo size="md" showText={false} />}
        <IconButton
          icon={isRTL ? ChevronRight : ChevronLeft}
          size="sm"
          variant="ghost"
          onClick={toggleSidebar}
          aria-label={t("dashboard.toggle_sidebar")}
          className={cn(isSidebarCollapsed && (isRTL ? "rotate-180" : "rotate-180"))}
        />
      </div>

      {/* Navigation */}
      <nav className="p-2 space-y-1">
        {navItems.map((item) => {
          const isActive = location.pathname.startsWith(item.path);
          const Icon = item.icon;

          return (
            <Link
              key={item.key}
              to={item.path}
              className={cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors",
                "hover:bg-base-300",
                isActive && "bg-primary-50 text-primary-700 font-semibold",
                !isActive && "text-base-content"
              )}
              title={t(`dashboard.${item.key}`)}
            >
              <Icon size={20} className="flex-shrink-0" />
              {!isSidebarCollapsed && (
                <span className="text-sm">{t(`dashboard.${item.key}`)}</span>
              )}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
