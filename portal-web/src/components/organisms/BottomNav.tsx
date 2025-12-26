import { useTranslation } from "react-i18next";
import { Link, useLocation } from "react-router-dom";
import {
  LayoutDashboard,
  Package,
  ShoppingCart,
  Users,
  MoreHorizontal,
} from "lucide-react";
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
  { key: "more", icon: MoreHorizontal, path: "/more" },
];

/**
 * BottomNav Component
 *
 * Fixed bottom navigation bar for mobile devices.
 *
 * Features:
 * - Shows 5 core navigation items
 * - Active route highlighting
 * - Touch-friendly (44px+ target)
 * - RTL support
 */
export function BottomNav() {
  const { t } = useTranslation();
  const location = useLocation();

  return (
    <nav className="fixed bottom-0 start-0 end-0 h-16 bg-base-200 border-t border-base-300 z-50 md:hidden">
      <div className="h-full flex items-center justify-around px-2">
        {navItems.map((item) => {
          const isActive =
            item.path === "/more"
              ? false // "More" is not a direct route
              : location.pathname.startsWith(item.path);
          const Icon = item.icon;

          return (
            <Link
              key={item.key}
              to={item.path}
              className={cn(
                "flex flex-col items-center justify-center gap-1 flex-1 h-full min-w-0",
                "transition-colors active:scale-95",
                isActive && "text-primary-600",
                !isActive && "text-base-content/60"
              )}
            >
              <Icon size={20} className="flex-shrink-0" />
              <span className="text-xs font-medium truncate">
                {t(`dashboard.${item.key}`)}
              </span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
