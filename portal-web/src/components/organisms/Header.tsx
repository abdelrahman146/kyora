import { useTranslation } from "react-i18next";
import { Menu } from "lucide-react";
import { IconButton } from "../atoms/IconButton";
import { BusinessSwitcher } from "../molecules/BusinessSwitcher";
import { UserMenu } from "../molecules/UserMenu";
import { useBusinessStore } from "../../stores/businessStore";
import { useLanguage } from "../../hooks/useLanguage";
import { useMediaQuery } from "../../hooks/useMediaQuery";
import { cn } from "@/lib/utils";

interface HeaderProps {
  title?: string;
}

/**
 * Header Component
 *
 * Top header bar containing business switcher, page title, and user menu.
 *
 * Features:
 * - Business switcher dropdown
 * - Page title (responsive, hidden on mobile if space is tight)
 * - User profile menu
 * - Mobile: Hamburger menu to open sidebar drawer
 * - Desktop: No hamburger (sidebar is always visible)
 * - RTL support with logical properties
 * - Responsive padding based on sidebar state
 * - Fixed positioning with proper z-index stacking
 */
export function Header({ title }: HeaderProps) {
  const { t } = useTranslation();
  const { toggleLanguage } = useLanguage();
  const isDesktop = useMediaQuery("(min-width: 768px)");
  const { isSidebarCollapsed, openSidebar } = useBusinessStore();

  return (
    <header
      className={cn(
        "fixed top-0 end-0 h-16 bg-base-100 border-b border-base-300 z-30 transition-all duration-300",
        // Desktop: Adjust start position based on sidebar width
        isDesktop && !isSidebarCollapsed && "md:start-64",
        isDesktop && isSidebarCollapsed && "md:start-20",
        // Mobile: Full width
        "start-0"
      )}
    >
      <div className="h-full flex items-center justify-between gap-2 px-4">
        {/* Left Section: Mobile Menu + Title */}
        <div className="flex items-center gap-3 flex-1 min-w-0">
          {/* Mobile Menu Toggle - Opens Sidebar Drawer */}
          {!isDesktop && (
            <IconButton
              icon={Menu}
              size="md"
              variant="ghost"
              onClick={openSidebar}
              aria-label={t("dashboard.open_menu")}
            />
          )}

          {/* Page Title - Hidden on small screens to save space */}
          {title && (
            <h1 className="text-lg font-bold text-base-content truncate">
              {title}
            </h1>
          )}
        </div>

        {/* Right Section: Business Switcher + Language Toggle + User Menu */}
        <div className="flex items-center gap-2 shrink-0">
          {/* Business Switcher - Compact on mobile */}
          <BusinessSwitcher />

          {/* Language Toggle - Show only on desktop or larger mobile */}
          <div className="hidden sm:block">
            <LanguageSwitcher variant="compact" />
          </div>

          {/* User Menu */}
          <UserMenu />
        </div>
      </div>
    </header>
  );
}
