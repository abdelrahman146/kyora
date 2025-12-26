import { type ReactNode, useEffect } from "react";
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
 * Mobile-first application layout with responsive behavior.
 *
 * Features:
 * - Mobile: Fixed header + Bottom nav + Drawer sidebar (overlay)
 * - Desktop: Fixed sidebar + Header (no bottom nav)
 * - Collapsible sidebar on desktop
 * - Overlay with backdrop on mobile when sidebar opens
 * - Safe area padding for mobile devices
 * - RTL support with logical properties
 * - Smooth transitions and animations
 * - Thumb-friendly touch targets (minimum 44px)
 *
 * Layout Structure:
 * ```
 * Desktop (≥768px):
 * ┌─────────────┬──────────────────────────┐
 * │   Sidebar   │         Header           │
 * │  (Fixed)    ├──────────────────────────┤
 * │             │                          │
 * │  Collapsible│        Content           │
 * │  64→20px    │      (Responsive)        │
 * │             │                          │
 * └─────────────┴──────────────────────────┘
 *
 * Mobile (<768px):
 * ┌──────────────────────────────────────┐
 * │         Header (Fixed)               │
 * ├──────────────────────────────────────┤
 * │                                      │
 * │          Content Area                │
 * │      (Scrollable, Safe)              │
 * │                                      │
 * ├──────────────────────────────────────┤
 * │       Bottom Nav (Fixed)             │
 * └──────────────────────────────────────┘
 *
 * Mobile Sidebar (Drawer Overlay):
 * ┌──────────────┬───────────────────────┐
 * │   Sidebar    │   Backdrop (Blur)     │
 * │   (Drawer)   │                       │
 * │   Slide-in   │   Tap to Close        │
 * │              │                       │
 * └──────────────┴───────────────────────┘
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
  const { isSidebarCollapsed, isSidebarOpen, closeSidebar } = useBusinessStore();

  // Close sidebar when switching to desktop
  useEffect(() => {
    if (isDesktop && isSidebarOpen) {
      closeSidebar();
    }
  }, [isDesktop, isSidebarOpen, closeSidebar]);

  // Prevent body scroll when mobile sidebar is open
  useEffect(() => {
    if (!isDesktop && isSidebarOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }

    return () => {
      document.body.style.overflow = "";
    };
  }, [isDesktop, isSidebarOpen]);

  return (
    <div className="min-h-screen bg-base-100">
      {/* Sidebar - Desktop: Fixed, Mobile: Drawer Overlay */}
      <Sidebar />

      {/* Mobile Sidebar Backdrop */}
      {!isDesktop && isSidebarOpen && (
        <div
          className="fixed inset-0 bg-base-content/40 backdrop-blur-sm z-40 md:hidden transition-opacity duration-300"
          onClick={closeSidebar}
          onKeyDown={(e) => {
            if (e.key === "Escape") closeSidebar();
          }}
          role="button"
          tabIndex={0}
          aria-label="Close sidebar"
        />
      )}

      {/* Header - Fixed at top, adjusts based on sidebar state */}
      <Header title={title} />

      {/* Main Content Area */}
      <main
        className={cn(
          // Base: Content below header
          "min-h-screen pt-16 transition-all duration-300",
          // Desktop: Adjust for sidebar width
          isDesktop && !isSidebarCollapsed && "md:ms-64",
          isDesktop && isSidebarCollapsed && "md:ms-20",
          // Mobile: Add bottom nav padding
          !isDesktop && "pb-20"
        )}
      >
        {/* Content Container with max-width and padding */}
        <div className="container mx-auto px-4 py-6 max-w-7xl">
          {children}
        </div>
      </main>

      {/* Mobile Bottom Navigation - Only visible on mobile */}
      {!isDesktop && <BottomNav />}
    </div>
  );
}
