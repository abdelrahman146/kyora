import { Navigate, Outlet } from "react-router-dom";
import { OnboardingProvider } from "@/contexts/OnboardingContext";
import { useAuth } from "@/hooks/useAuth";

/**
 * Onboarding Root Layout
 *
 * Wraps all onboarding steps with:
 * - OnboardingProvider for state management
 * - Authentication guard (redirect if already logged in)
 * - Outlet for nested routes
 *
 * Route Structure:
 * /onboarding
 *   /plan - Plan selection
 *   /verify - Email verification
 *   /business - Business setup
 *   /payment - Payment checkout (for paid plans)
 *   /complete - Finalization and welcome
 */
export default function OnboardingRoot() {
  const { isAuthenticated, isLoading } = useAuth();

  // Show loading while checking auth status
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-base-100">
        <div className="flex flex-col items-center gap-4">
          <span className="loading loading-spinner loading-lg text-primary"></span>
        </div>
      </div>
    );
  }

  // Redirect to dashboard if already authenticated
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <OnboardingProvider>
      <Outlet />
    </OnboardingProvider>
  );
}
