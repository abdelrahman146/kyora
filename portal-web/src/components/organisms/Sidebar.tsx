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
  X,
} from "lucide-react";
import { Logo } from "../atoms/Logo";
import { IconButton } from "../atoms/IconButton";
import { useBusinessStore } from "../../stores/businessStore";
import { useLanguage } from "../../hooks/useLanguage";
import { useMediaQuery } from "../../hooks/useMediaQuery";
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
 * Responsive navigation sidebar:
 * - Desktop: Fixed vertical sidebar with collapse functionality
 * - Mobile: Drawer overlay (slides from start)
 *
 * Features:
 * - Collapsible on desktop (full width ↔ icon-only)
 * - Drawer overlay on mobile with slide animation
 * - Active route highlighting
 * - RTL support with proper directional logic
 * - Smooth transitions
 * - Touch-friendly targets (min 44px)
 */
export function Sidebar() {
  const { t } = useTranslation();
  const location = useLocation();
  const { isRTL } = useLanguage();
  const isDesktop = useMediaQuery("(min-width: 768px)");
  const { isSidebarCollapsed, isSidebarOpen, toggleSidebar, closeSidebar, selectedBusiness } =
    useBusinessStore();

  // On mobile, always render (CSS handles visibility with transform)
  // On desktop, always visible

  return (
    <aside
      className={cn(
        "fixed top-0 h-screen bg-base-100 border-base-300 z-50",
        // Desktop: Always visible, collapsible width, smooth transition
        isDesktop && "border-e transition-all duration-300",
        isDesktop && !isSidebarCollapsed && "w-64",
        isDesktop && isSidebarCollapsed && "w-20",
        // Mobile: Drawer from start with slide animation
        !isDesktop && "w-64 shadow-2xl transition-transform duration-300",
        !isDesktop && !isRTL && "start-0",
        !isDesktop && isRTL && "end-0",
        // Mobile animation states
        !isDesktop && !isRTL && !isSidebarOpen && "-translate-x-full",
        !isDesktop && isRTL && !isSidebarOpen && "translate-x-full"
      )}
    >
      {/* Header Section */}
      <div className={cn(
        "h-16 flex items-center border-b border-base-300",
        isDesktop && isSidebarCollapsed ? "justify-center gap-2 px-2" : "justify-between px-4"
      )}>
        {isDesktop && isSidebarCollapsed ? (
          // Collapsed desktop: Center logo and chevron with gap
          <>
            <Logo size="md" showText={false} />
            <IconButton
              icon={isRTL ? ChevronLeft : ChevronRight}
              size="sm"
              variant="ghost"
              onClick={toggleSidebar}
              aria-label={t("dashboard.toggle_sidebar")}
            />
          </>
        ) : (
          // Expanded: Logo on start, toggle/close on end
          <>
            <div className="flex-1">
              <Logo size="md" showText />
            </div>
            {isDesktop ? (
              <IconButton
                icon={isRTL ? ChevronRight : ChevronLeft}
                size="sm"
                variant="ghost"
                onClick={toggleSidebar}
                aria-label={t("dashboard.toggle_sidebar")}
              />
            ) : (
              <IconButton
                icon={X}
                size="sm"
                variant="ghost"
                onClick={closeSidebar}
                aria-label={t("common.close")}
              />
            )}
          </>
        )}
      </div>

      {/* Navigation Links */}
      <nav className="p-2 space-y-1 overflow-y-auto h-[calc(100vh-4rem)]">
        {navItems.map((item) => {
          const itemPath = selectedBusiness && item.key !== "dashboard" 
            ? `/dashboard/${selectedBusiness.descriptor}${item.path}`
            : item.path;
          const isActive = location.pathname.startsWith(itemPath);
          const Icon = item.icon;

          return (
            <Link
              key={item.key}
              to={itemPath}
              onClick={() => {
                // Close drawer on mobile after navigation
                if (!isDesktop) closeSidebar();
              }}
              className={cn(
                "flex items-center gap-3 px-4 py-3 rounded-lg transition-colors",
                "hover:bg-base-200 active:scale-98",
                // Active state styling
                isActive && "bg-primary/10 text-primary font-semibold",
                !isActive && "text-base-content",
                // Desktop collapsed: center icon
                isDesktop && isSidebarCollapsed && "justify-center px-0"
              )}
              title={isSidebarCollapsed ? t(`dashboard.${item.key}`) : undefined}
            >
              <Icon size={20} className="shrink-0" />
              {/* Hide text when collapsed on desktop */}
              {(!isDesktop || !isSidebarCollapsed) && (
                <span className="text-sm truncate">{t(`dashboard.${item.key}`)}</span>
              )}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
