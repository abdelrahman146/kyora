import { useCallback, useEffect, useRef, useState } from 'react'
import type { ReactNode } from 'react'

import type {
  CreateOrderRequest,
  OrderPaymentMethod,
  OrderPaymentStatus,
  OrderPreview,
  OrderStatus,
} from '@/api/order'
import { usePreviewOrderMutation } from '@/api/order'

// Constants for preview behavior
const PREVIEW_DEBOUNCE_MS = 800 // Wait 800ms after last change before preview
const PREVIEW_STALE_MS = 30000 // Preview becomes stale after 30 seconds
const MIN_REQUEST_INTERVAL_MS = 2000 // Minimum 2s between requests to avoid rate limiting

export type PreviewState = 'idle' | 'loading' | 'success' | 'error' | 'stale'

export interface PreviewResult {
  previewData: OrderPreview | null
  isLoading: boolean
  isStale: boolean
  errorMessage: string | null
  previewState: PreviewState
  lastPreviewAt: Date | null
  canSubmit: boolean
  triggerPreview: () => void
}

export interface OrderFormValues {
  customerId: string
  shippingAddressId: string
  channel: string
  shippingZoneId?: string
  discountType?: string
  discountValue?: string
  paymentMethod?: string
  paymentReference?: string
  status?: string
  paymentStatus?: string
  note?: string
  items: Array<{
    variantId: string
    quantity: number
    unitPrice: string
    unitCost?: string
  }>
}

interface OrderPreviewManagerProps {
  businessDescriptor: string
  isOpen: boolean
  formValues: OrderFormValues | null
  children: (props: PreviewResult) => ReactNode
}

/**
 * Builds a preview payload from form values.
 * Returns null if form values are not eligible for preview.
 */
function buildPreviewPayload(
  values: OrderFormValues | null,
): CreateOrderRequest | null {
  if (!values) return null

  // Filter valid items (must have variantId, quantity > 0, and a price)
  const validItems = values.items
    .filter(
      (item) =>
        item.variantId &&
        item.quantity &&
        item.quantity > 0 &&
        item.unitPrice !== '',
    )
    .map((item) => ({
      variantId: item.variantId,
      quantity: Number(item.quantity) || 0,
      unitPrice: item.unitPrice,
      unitCost: item.unitCost || undefined,
    }))

  // Must have customer, address, and at least one item
  if (
    !values.customerId ||
    !values.shippingAddressId ||
    validItems.length === 0
  ) {
    return null
  }

  return {
    customerId: values.customerId,
    shippingAddressId: values.shippingAddressId,
    channel: values.channel || 'instagram',
    shippingZoneId: values.shippingZoneId || undefined,
    discountType:
      values.discountValue && values.discountValue.trim() !== ''
        ? (values.discountType as 'amount' | 'percent')
        : undefined,
    discountValue:
      values.discountValue && values.discountValue.trim() !== ''
        ? values.discountValue
        : undefined,
    paymentMethod: (values.paymentMethod || undefined) as
      | OrderPaymentMethod
      | undefined,
    paymentReference: values.paymentReference || undefined,
    status: (values.status || undefined) as OrderStatus | undefined,
    paymentStatus: (values.paymentStatus || undefined) as
      | OrderPaymentStatus
      | undefined,
    note: values.note || undefined,
    items: validItems,
  }
}

/**
 * Creates a stable key from the payload for comparison.
 * Only includes fields that affect pricing/totals.
 */
function getPayloadKey(payload: CreateOrderRequest | null): string {
  if (!payload) return ''

  // Only include fields that affect the preview calculation
  const keyData = {
    customerId: payload.customerId,
    shippingAddressId: payload.shippingAddressId,
    shippingZoneId: payload.shippingZoneId,
    discountType: payload.discountType,
    discountValue: payload.discountValue,
    items: payload.items.map((i) => ({
      variantId: i.variantId,
      quantity: i.quantity,
      unitPrice: i.unitPrice,
      unitCost: i.unitCost,
    })),
  }

  return JSON.stringify(keyData)
}

/**
 * OrderPreviewManager - Manages dry-run preview API calls with proper debouncing.
 *
 * Key behaviors:
 * - Debounces preview calls (waits 800ms after last form change)
 * - Tracks stale state (preview older than 30s)
 * - Silently ignores rate limit errors (429)
 * - Provides canSubmit based on successful preview
 * - Exposes manual triggerPreview for retry
 */
export function OrderPreviewManager({
  businessDescriptor,
  isOpen,
  formValues,
  children,
}: OrderPreviewManagerProps) {
  // Preview state
  const [previewData, setPreviewData] = useState<OrderPreview | null>(null)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [lastPreviewAt, setLastPreviewAt] = useState<Date | null>(null)
  const [isStale, setIsStale] = useState(false)

  // Refs for tracking state across renders
  const lastSuccessfulPayloadKeyRef = useRef<string>('')
  const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const staleTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const lastRequestTimeRef = useRef<number>(0)
  const pendingPayloadRef = useRef<CreateOrderRequest | null>(null)

  // Mutation hook - let global error handler show translated toasts (except 429)
  const previewMutation = usePreviewOrderMutation(businessDescriptor)

  // Build current payload
  const currentPayload = buildPreviewPayload(formValues)
  const currentPayloadKey = getPayloadKey(currentPayload)
  const isEligible = currentPayload !== null

  // Execute preview request
  const executePreview = useCallback(
    async (payload: CreateOrderRequest) => {
      // Check minimum request interval to avoid rate limiting
      const now = Date.now()
      const timeSinceLastRequest = now - lastRequestTimeRef.current
      if (timeSinceLastRequest < MIN_REQUEST_INTERVAL_MS) {
        // Schedule for later instead of firing immediately
        const delay = MIN_REQUEST_INTERVAL_MS - timeSinceLastRequest + 100
        pendingPayloadRef.current = payload
        debounceTimerRef.current = setTimeout(() => {
          const pending = pendingPayloadRef.current
          if (pending) {
            pendingPayloadRef.current = null
            void executePreview(pending)
          }
        }, delay)
        return
      }

      lastRequestTimeRef.current = now
      pendingPayloadRef.current = null

      try {
        const result = await previewMutation.mutateAsync(payload)

        // Update state on success
        setPreviewData(result)
        setErrorMessage(null)
        setLastPreviewAt(new Date())
        setIsStale(false)
        lastSuccessfulPayloadKeyRef.current = getPayloadKey(payload)

        // Start stale timer
        if (staleTimerRef.current) {
          clearTimeout(staleTimerRef.current)
        }
        staleTimerRef.current = setTimeout(() => {
          setIsStale(true)
        }, PREVIEW_STALE_MS)
      } catch (error: unknown) {
        // Type-safe error handling for HTTP errors
        const err = error as
          | { response?: { status?: number }; status?: number }
          | undefined

        // Silently ignore rate limit errors (429) and retry
        const status = err?.response?.status ?? err?.status
        if (status === 429) {
          // Keep previous state, don't show error
          // Schedule a retry after minimum interval
          debounceTimerRef.current = setTimeout(() => {
            void executePreview(payload)
          }, MIN_REQUEST_INTERVAL_MS)
          return
        }

        // For other errors: set UI error state but DON'T suppress the error
        // Remove meta: { errorToast: 'off' } and let error propagate to global handler
        setErrorMessage('error')
        // Re-throw so global mutation error handler can show translated toast
        throw error
      }
    },
    [previewMutation],
  )

  // Manual trigger for retry button
  const triggerPreview = useCallback(() => {
    if (!currentPayload) return

    // Clear any pending debounce
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current)
      debounceTimerRef.current = null
    }

    // Execute with a slight delay to avoid double-taps
    lastRequestTimeRef.current = 0 // Allow immediate execution
    void executePreview(currentPayload)
  }, [currentPayload, executePreview])

  // Effect: Handle form value changes with debouncing
  useEffect(() => {
    // Skip if not open or not eligible for preview
    if (!isOpen || !isEligible) {
      return
    }

    // Skip if payload hasn't changed from last successful preview
    if (currentPayloadKey === lastSuccessfulPayloadKeyRef.current) {
      return
    }

    // Mark as stale when payload changes
    setIsStale(true)

    // Clear previous debounce timer
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current)
    }

    // Set new debounce timer
    debounceTimerRef.current = setTimeout(() => {
      void executePreview(currentPayload)
    }, PREVIEW_DEBOUNCE_MS)

    // Cleanup on unmount or dependency change
    return () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
        debounceTimerRef.current = null
      }
    }
    // Note: previewData intentionally excluded to avoid infinite loops
  }, [isOpen, isEligible, currentPayloadKey, currentPayload, executePreview])

  // Effect: Reset state when sheet closes
  useEffect(() => {
    if (!isOpen) {
      setPreviewData(null)
      setErrorMessage(null)
      setLastPreviewAt(null)
      setIsStale(false)
      lastSuccessfulPayloadKeyRef.current = ''
      pendingPayloadRef.current = null

      // Clear timers
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
        debounceTimerRef.current = null
      }
      if (staleTimerRef.current) {
        clearTimeout(staleTimerRef.current)
        staleTimerRef.current = null
      }
    }
  }, [isOpen])

  // Effect: Reset preview when form becomes ineligible
  useEffect(() => {
    if (!isEligible && isOpen) {
      setPreviewData(null)
      setErrorMessage(null)
      setLastPreviewAt(null)
      setIsStale(false)
    }
  }, [isEligible, isOpen])

  // Compute preview state
  const previewState: PreviewState = previewMutation.isPending
    ? 'loading'
    : errorMessage
      ? 'error'
      : isStale
        ? 'stale'
        : previewData
          ? 'success'
          : 'idle'

  // Can submit only if we have successful, non-stale preview
  const canSubmit = Boolean(
    previewData && isEligible && !isStale && !previewMutation.isPending,
  )

  return children({
    previewData,
    isLoading: previewMutation.isPending,
    isStale,
    errorMessage,
    previewState,
    lastPreviewAt,
    canSubmit,
    triggerPreview,
  })
}
