import { type ReactNode } from "react";
import { Header } from "../organisms/Header";
import { Sidebar } from "../organisms/Sidebar";
import { BottomNav } from "../organisms/BottomNav";
import { useBusinessStore } from "../../stores/businessStore";
import { useMediaQuery } from "../../hooks/useMediaQuery";
import { cn } from "@/lib/utils";

interface DashboardLayoutProps {
  children: ReactNode;
  title?: string;
}

/**
 * DashboardLayout Component
 *
 * Main application layout wrapping all dashboard pages.
 *
 * Features:
 * - Responsive layout (Desktop: Sidebar + Header, Mobile: Bottom Nav + Header)
 * - Automatic sidebar collapse handling
 * - Safe area padding for mobile devices
 * - RTL support
 * - Business context management
 *
 * Layout Structure:
 * ```
 * Desktop:
 * ┌─────────────┬──────────────────────────┐
 * │   Sidebar   │         Header           │
 * │             ├──────────────────────────┤
 * │             │                          │
 * │  Navigation │        Content           │
 * │             │                          │
 * │             │                          │
 * └─────────────┴──────────────────────────┘
 *
 * Mobile:
 * ┌──────────────────────────────────────┐
 * │             Header                   │
 * ├──────────────────────────────────────┤
 * │                                      │
 * │            Content                   │
 * │                                      │
 * ├──────────────────────────────────────┤
 * │          Bottom Nav                  │
 * └──────────────────────────────────────┘
 * ```
 *
 * @example
 * ```tsx
 * <DashboardLayout title="Inventory">
 *   <InventoryPage />
 * </DashboardLayout>
 * ```
 */
export function DashboardLayout({ children, title }: DashboardLayoutProps) {
  const isDesktop = useMediaQuery("(min-width: 768px)");
  const { isSidebarCollapsed } = useBusinessStore();

  return (
    <div className="min-h-screen bg-base-100">
      {/* Desktop Sidebar */}
      {isDesktop && <Sidebar />}

      {/* Header */}
      <Header title={title} />

      {/* Main Content */}
      <main
        className={cn(
          "pt-16 transition-all duration-300",
          // Desktop: Add sidebar width padding
          isDesktop && !isSidebarCollapsed && "md:ps-64",
          isDesktop && isSidebarCollapsed && "md:ps-20",
          // Mobile: Add bottom nav height padding
          !isDesktop && "pb-16"
        )}
      >
        <div className="container mx-auto p-4 max-w-7xl">{children}</div>
      </main>

      {/* Mobile Bottom Navigation */}
      {!isDesktop && <BottomNav />}
    </div>
  );
}
