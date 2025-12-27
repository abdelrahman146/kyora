import {
  createContext,
  useContext,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import type {
  SessionStage,
  Plan,
  StartSessionResponse,
} from "@/api/types/onboarding";

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
}

interface OnboardingContextValue extends OnboardingState {
  // Session management
  startSession: (response: StartSessionResponse, email: string, plan: Plan) => void;
  updateStage: (stage: SessionStage) => void;
  
  // Plan selection
  setSelectedPlan: (plan: Plan) => void;
  
  // Email verification
  markEmailVerified: () => void;
  
  // Business setup
  setBusinessDetails: (
    name: string,
    descriptor: string,
    country: string,
    currency: string
  ) => void;
  
  // Payment
  setCheckoutUrl: (url: string | null) => void;
  markPaymentComplete: () => void;
  
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
  const [state, setState] = useState<OnboardingState>(() => {
    // Try to restore session from sessionStorage
    const stored = sessionStorage.getItem("kyora_onboarding_state");
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch {
        return INITIAL_STATE;
      }
    }
    return INITIAL_STATE;
  });

  // Persist state to sessionStorage whenever it changes
  const updateState = useCallback((updates: Partial<OnboardingState>) => {
    setState((prev) => {
      const newState = { ...prev, ...updates };
      sessionStorage.setItem("kyora_onboarding_state", JSON.stringify(newState));
      return newState;
    });
  }, []);

  const startSession = useCallback(
    (response: StartSessionResponse, email: string, plan: Plan) => {
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

  const markEmailVerified = useCallback(() => {
    updateState({ isEmailVerified: true, stage: "identity_verified" });
  }, [updateState]);

  const setBusinessDetails = useCallback(
    (
      name: string,
      descriptor: string,
      country: string,
      currency: string
    ) => {
      updateState({
        businessName: name,
        businessDescriptor: descriptor,
        businessCountry: country,
        businessCurrency: currency,
        stage: "business_staged",
      });
    },
    [updateState]
  );

  const setCheckoutUrl = useCallback(
    (url: string | null) => {
      updateState({ checkoutUrl: url });
    },
    [updateState]
  );

  const markPaymentComplete = useCallback(() => {
    updateState({ isPaymentComplete: true, stage: "ready_to_commit" });
  }, [updateState]);

  const resetOnboarding = useCallback(() => {
    setState(INITIAL_STATE);
    sessionStorage.removeItem("kyora_onboarding_state");
  }, []);

  const canProceedToNextStep = useCallback((): boolean => {
    const { stage, isPaidPlan, isPaymentComplete } = state;

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
    updateStage,
    setSelectedPlan,
    markEmailVerified,
    setBusinessDetails,
    setCheckoutUrl,
    markPaymentComplete,
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
