import {
  createContext,
  useContext,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import {
  type SessionStage,
  type Plan,
  type StartSessionResponse,
  type GetSessionResponse,
} from "@/api/types/onboarding";
import { onboardingApi } from "@/api/onboarding";
import { sessionStorage as sessionTokenStorage } from "@/lib/sessionStorage";

/**
 * Onboarding Context State
 *
 * Manages the entire onboarding flow state including:
 * - Session token for tracking progress
 * - Current stage of onboarding
 * - Selected plan information
 * - Email verification status
 * - Business setup status
 * - Payment completion status
 */

interface OnboardingState {
  // Session tracking
  sessionToken: string | null;
  stage: SessionStage | null;
  
  // Plan selection
  selectedPlan: Plan | null;
  isPaidPlan: boolean;
  
  // Email and identity
  email: string | null;
  isEmailVerified: boolean;
  
  // Business details
  businessName: string | null;
  businessDescriptor: string | null;
  businessCountry: string | null;
  businessCurrency: string | null;
  
  // Payment
  isPaymentComplete: boolean;
  checkoutUrl: string | null;
  
  // Loading and error states
  isLoading: boolean;
  error: string | null;
}

interface OnboardingContextValue extends OnboardingState {
  // Session management
  startSession: (response: StartSessionResponse, email: string, plan: Plan) => void;
  loadSession: (sessionToken: string) => Promise<void>;
  loadSessionFromStorage: () => Promise<boolean>;
  clearSession: () => Promise<void>;
  updateStage: (stage: SessionStage) => void;
  
  // Plan selection
  setSelectedPlan: (plan: Plan) => void;
  
  // Payment
  setCheckoutUrl: (url: string | null) => void;
  
  // Reset
  resetOnboarding: () => void;
  
  // Helpers
  canProceedToNextStep: () => boolean;
  getCurrentStepNumber: () => number;
}

const OnboardingContext = createContext<OnboardingContextValue | undefined>(
  undefined
);

const INITIAL_STATE: OnboardingState = {
  sessionToken: null,
  stage: null,
  selectedPlan: null,
  isPaidPlan: false,
  email: null,
  isEmailVerified: false,
  businessName: null,
  businessDescriptor: null,
  businessCountry: null,
  businessCurrency: null,
  isPaymentComplete: false,
  checkoutUrl: null,
  isLoading: false,
  error: null,
};

/**
 * OnboardingProvider Component
 *
 * Provides onboarding state and actions to all child components.
 * Should wrap the onboarding route tree.
 *
 * @example
 * ```tsx
 * <OnboardingProvider>
 *   <OnboardingRoutes />
 * </OnboardingProvider>
 * ```
 */
export function OnboardingProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<OnboardingState>(INITIAL_STATE);

  // Update state without persistence (backend is source of truth)
  const updateState = useCallback((updates: Partial<OnboardingState>) => {
    setState((prev) => ({ ...prev, ...updates }));
  }, []);

  const startSession = useCallback(
    (response: StartSessionResponse, email: string, plan: Plan) => {
      // Save token to localStorage for persistence
      sessionTokenStorage.setToken(response.sessionToken);
      
      updateState({
        sessionToken: response.sessionToken,
        stage: response.stage,
        email,
        selectedPlan: plan,
        isPaidPlan: response.isPaid,
      });
    },
    [updateState]
  );

  const updateStage = useCallback(
    (stage: SessionStage) => {
      updateState({ stage });
    },
    [updateState]
  );

  const setSelectedPlan = useCallback(
    (plan: Plan) => {
      updateState({ selectedPlan: plan });
    },
    [updateState]
  );

  /**
   * Load session from backend by token.
   * Updates all state from backend response.
   */
  const loadSession = useCallback(
    async (sessionToken: string) => {
      updateState({ isLoading: true, error: null });
      
      try {
        const session: GetSessionResponse = await onboardingApi.getSession(sessionToken);
        
        // Update state from backend session
        updateState({
          sessionToken: session.sessionToken,
          stage: session.stage,
          email: session.email,
          isEmailVerified: session.emailVerified,
          selectedPlan: session.planId
            ? {
                id: session.planId,
                descriptor: session.planDescriptor,
                name: session.planDescriptor, // Backend doesn't return name, use descriptor
                price: "0", // Backend doesn't return price in session
                currency: session.businessCurrency || "USD",
                billingCycle: "monthly",
                features: {
                  customerManagement: false,
                  inventoryManagement: false,
                  orderManagement: false,
                  expenseManagement: false,
                  accounting: false,
                  basicAnalytics: false,
                  financialReports: false,
                  dataImport: false,
                  dataExport: false,
                  advancedAnalytics: false,
                  advancedFinancialReports: false,
                  orderPaymentLinks: false,
                  invoiceGeneration: false,
                  exportAnalyticsData: false,
                  aiBusinessAssistant: false,
                },
                limits: {
                  maxTeamMembers: 1,
                  maxBusinesses: 1,
                  maxOrdersPerMonth: 100,
                },
                description: "",
              }
            : null,
          isPaidPlan: session.isPaidPlan,
          businessName: session.businessName,
          businessDescriptor: session.businessDescriptor,
          businessCountry: session.businessCountry,
          businessCurrency: session.businessCurrency,
          isPaymentComplete:
            session.paymentStatus === "succeeded" ||
            session.paymentStatus === "skipped",
          checkoutUrl: null,
          isLoading: false,
          error: null,
        });
      } catch (error) {
        const message = error instanceof Error ? error.message : "Failed to load session";
        updateState({ isLoading: false, error: message });
        throw error;
      }
    },
    [updateState]
  );

  /**
   * Load session from localStorage if it exists.
   * Returns true if session was loaded successfully.
   */
  const loadSessionFromStorage = useCallback(async (): Promise<boolean> => {
    const token = sessionTokenStorage.getToken();
    if (!token) return false;

    try {
      await loadSession(token);
      return true;
    } catch {
      // Token invalid or expired, clear it
      sessionTokenStorage.clearToken();
      return false;
    }
  }, [loadSession]);

  /**
   * Delete session from backend and clear localStorage.
   * Used when user chooses "Start Fresh".
   */
  const clearSession = useCallback(async () => {
    if (state.sessionToken) {
      try {
        await onboardingApi.deleteSession(state.sessionToken);
      } catch (error) {
        console.error("Failed to delete session:", error);
      }
    }
    sessionTokenStorage.clearToken();
    setState(INITIAL_STATE);
  }, [state.sessionToken]);

  const setCheckoutUrl = useCallback(
    (url: string | null) => {
      updateState({ checkoutUrl: url });
    },
    [updateState]
  );

  const resetOnboarding = useCallback(() => {
    sessionTokenStorage.clearToken();
    setState(INITIAL_STATE);
  }, []);

  const canProceedToNextStep = useCallback((): boolean => {
    const { stage } = state;

    switch (stage) {
      case "plan_selected":
        return false; // Need to verify email
      case "identity_pending":
        return false; // Waiting for OTP verification
      case "identity_verified":
        return true; // Can proceed to business setup
      case "business_staged":
        return true; // Can proceed to payment or complete
      case "payment_pending":
        return false; // Waiting for payment
      case "payment_confirmed":
        return true; // Can complete onboarding
      case "ready_to_commit":
        return true; // Can complete onboarding
      default:
        return false;
    }
  }, [state]);

  const getCurrentStepNumber = useCallback((): number => {
    const { stage } = state;

    switch (stage) {
      case "plan_selected":
      case "identity_pending":
        return 1; // Email verification step
      case "identity_verified":
        return 2; // Business setup step
      case "business_staged":
      case "payment_pending":
      case "payment_confirmed":
        return 3; // Payment step
      case "ready_to_commit":
        return 4; // Complete step
      default:
        return 0; // Plan selection step
    }
  }, [state]);

  const value: OnboardingContextValue = {
    ...state,
    startSession,
    loadSession,
    loadSessionFromStorage,
    clearSession,
    updateStage,
    setSelectedPlan,
    setCheckoutUrl,
    resetOnboarding,
    canProceedToNextStep,
    getCurrentStepNumber,
  };

  return (
    <OnboardingContext.Provider value={value}>
      {children}
    </OnboardingContext.Provider>
  );
}

/**
 * useOnboarding Hook
 *
 * Access onboarding state and actions from any component within OnboardingProvider.
 *
 * @throws Error if used outside OnboardingProvider
 *
 * @example
 * ```tsx
 * const { sessionToken, stage, updateStage } = useOnboarding();
 * ```
 */
export function useOnboarding(): OnboardingContextValue {
  const context = useContext(OnboardingContext);

  if (context === undefined) {
    throw new Error("useOnboarding must be used within OnboardingProvider");
  }

  return context;
}
