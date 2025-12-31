import { Store } from '@tanstack/react-store'
import type { GetSessionResponse } from '@/api/types/onboarding'
import { onboardingApi } from '@/api/onboarding'
import { createPersistencePlugin } from '@/lib/storePersistence'

/**
 * Onboarding Session Stage
 *
 * Matches backend session stages for consistency
 */
export type SessionStage =
  | 'plan_selected'
  | 'identity_pending'
  | 'identity_verified'
  | 'business_staged'
  | 'payment_pending'
  | 'payment_confirmed'
  | 'ready_to_commit'
  | 'committed'

/**
 * Business Data (staged during onboarding)
 */
export interface BusinessData {
  name: string
  descriptor: string
  country: string
  currency: string
}

/**
 * Onboarding Store State
 *
 * Tracks the entire onboarding flow including:
 * - Session token for backend persistence
 * - Current stage in the onboarding process
 * - User's email and selected plan
 * - Business details (staged, not committed)
 * - Payment completion status
 */
interface OnboardingState {
  sessionToken: string | null
  stage: SessionStage | null
  email: string | null
  planId: string | null
  planDescriptor: string | null
  isPaidPlan: boolean
  businessData: BusinessData | null
  paymentCompleted: boolean
  checkoutUrl: string | null
}

const initialState: OnboardingState = {
  sessionToken: null,
  stage: null,
  email: null,
  planId: null,
  planDescriptor: null,
  isPaidPlan: false,
  businessData: null,
  paymentCompleted: false,
  checkoutUrl: null,
}

/**
 * Create onboarding store with persistence
 *
 * Persists entire state to localStorage with no TTL.
 * Session is cleared manually after completion or abandoned.
 */
export const onboardingStore = new Store<OnboardingState>(initialState)

// Set up persistence plugin
const persistencePlugin = createPersistencePlugin({
  key: 'kyora_onboarding_session',
  store: onboardingStore,
  select: (state: OnboardingState) => state,
  // No TTL - sessions stay until explicitly cleared
})

// Initialize store from localStorage on app load
const persistedState = persistencePlugin.loadState()
if (persistedState) {
  onboardingStore.setState(() => persistedState)
}

/**
 * Start new onboarding session
 *
 * Called after user selects plan and provides email.
 * Stores session token from backend response.
 */
export function startSession(
  sessionToken: string,
  stage: SessionStage,
  email: string,
  planId: string,
  planDescriptor: string,
  isPaidPlan: boolean,
): void {
  onboardingStore.setState(() => ({
    sessionToken,
    stage,
    email,
    planId,
    planDescriptor,
    isPaidPlan,
    businessData: null,
    paymentCompleted: false,
    checkoutUrl: null,
  }))
}

/**
 * Update session stage
 *
 * Called when backend advances the onboarding flow.
 */
export function updateStage(stage: SessionStage): void {
  onboardingStore.setState((state) => ({
    ...state,
    stage,
  }))
}

/**
 * Set email
 *
 * Called during email input step.
 */
export function setEmail(email: string): void {
  onboardingStore.setState((state) => ({
    ...state,
    email,
  }))
}

/**
 * Set plan selection
 *
 * Called during plan selection step.
 */
export function setPlan(
  planId: string,
  planDescriptor: string,
  isPaidPlan: boolean,
): void {
  onboardingStore.setState((state) => ({
    ...state,
    planId,
    planDescriptor,
    isPaidPlan,
  }))
}

/**
 * Set business data
 *
 * Called during business setup step.
 * Business is staged locally but not committed until complete step.
 */
export function setBusiness(businessData: BusinessData): void {
  onboardingStore.setState((state) => ({
    ...state,
    businessData,
  }))
}

/**
 * Set payment completion
 *
 * Called after successful Stripe payment.
 */
export function setPaymentCompleted(completed: boolean): void {
  onboardingStore.setState((state) => ({
    ...state,
    paymentCompleted: completed,
  }))
}

/**
 * Set checkout URL
 *
 * Stores Stripe checkout URL for payment step.
 */
export function setCheckoutUrl(checkoutUrl: string | null): void {
  onboardingStore.setState((state) => ({
    ...state,
    checkoutUrl,
  }))
}

function setStateFromSession(session: GetSessionResponse): void {
  onboardingStore.setState((prev) => {
    const businessData: BusinessData | null =
      session.businessName &&
      session.businessDescriptor &&
      session.businessCountry &&
      session.businessCurrency
        ? {
            name: session.businessName,
            descriptor: session.businessDescriptor,
            country: session.businessCountry,
            currency: session.businessCurrency,
          }
        : prev.businessData

    const paymentCompleted =
      session.paymentStatus === 'succeeded' ||
      session.paymentStatus === 'skipped'

    return {
      ...prev,
      sessionToken: session.sessionToken,
      stage: session.stage,
      email: session.email,
      planId: session.planId,
      planDescriptor: session.planDescriptor,
      isPaidPlan: session.isPaidPlan,
      businessData,
      paymentCompleted,
    }
  })
}

/**
 * Load session from backend by token and hydrate store.
 *
 * This is the source of truth for resume flows.
 */
export async function loadSession(sessionToken: string): Promise<void> {
  const session = await onboardingApi.getSession(sessionToken)
  setStateFromSession(session)
}

/**
 * Load session using persisted token.
 * Returns true if a valid session was loaded.
 */
export async function loadSessionFromStorage(): Promise<boolean> {
  const token = onboardingStore.state.sessionToken
  if (!token) return false

  try {
    await loadSession(token)
    return true
  } catch {
    // Token invalid/expired -> clear persisted state.
    await clearSession()
    return false
  }
}

/**
 * Clear onboarding session
 *
 * Called after completion or when user abandons onboarding.
 * Clears both store state and localStorage.
 */
export async function clearSession(): Promise<void> {
  const token = onboardingStore.state.sessionToken
  if (token) {
    try {
      await onboardingApi.deleteSession(token)
    } catch {
      // Silent fail: session will expire server-side eventually.
    }
  }

  onboardingStore.setState(() => initialState)
  persistencePlugin.clearState()
}

/**
 * Get current stage number for progress indicator
 *
 * Returns 1-6 representing progress through onboarding.
 */
export function getCurrentStageNumber(): number {
  const stage = onboardingStore.state.stage
  if (!stage) return 0

  const stageMap: Record<SessionStage, number> = {
    plan_selected: 1,
    identity_pending: 2,
    identity_verified: 3,
    business_staged: 4,
    payment_pending: 5,
    payment_confirmed: 5,
    ready_to_commit: 6,
    committed: 6,
  }

  return stageMap[stage] || 0
}

/**
 * Initialize TanStack Store Devtools (dev-only)
 *
 * Conditionally loads devtools in development mode only.
 * Production builds will exclude this code via tree-shaking.
 */
if (import.meta.env.DEV) {
  console.log(
    '[onboardingStore] TanStack Store devtools enabled in development mode',
  )
}
