/**
 * Accounting API Client
 *
 * Business-scoped API for managing accounting data:
 * - Expenses (one-time and recurring)
 * - Investments (owner capital in)
 * - Withdrawals (owner capital out)
 * - Assets (fixed assets)
 * - Summary (financial overview including "Safe to Draw")
 *
 * Based on backend routes at: /v1/businesses/:businessDescriptor/accounting/...
 * See: .github/instructions/accounting.instructions.md
 */
import {
  keepPreviousData,
  queryOptions,
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'
import { del, get, patch, post } from './client'
import type { UseMutationOptions } from '@tanstack/react-query'

import type {
  AccountingSummary,
  Asset,
  AssetType,
  CreateAssetRequest,
  CreateExpenseRequest,
  CreateInvestmentRequest,
  CreateRecurringExpenseRequest,
  CreateWithdrawalRequest,
  Expense,
  ExpenseCategory,
  Investment,
  ListResponse,
  RecentActivitiesResponse,
  RecentActivity,
  RecurringExpense,
  UpdateAssetRequest,
  UpdateExpenseRequest,
  UpdateInvestmentRequest,
  UpdateRecurringExpenseRequest,
  UpdateRecurringExpenseStatusRequest,
  UpdateWithdrawalRequest,
  Withdrawal,
} from './types/accounting'
import { STALE_TIME } from '@/lib/queryKeys'

// Re-export types for convenience
export type {
  Expense,
  RecurringExpense,
  Investment,
  Withdrawal,
  Asset,
  AssetType,
  AccountingSummary,
  ExpenseCategory,
  RecentActivity,
  RecentActivitiesResponse,
}
export {
  expenseCategoryEnum,
  expenseTypeEnum,
  recurringExpenseFrequencyEnum,
  recurringExpenseStatusEnum,
  recentActivityTypeEnum,
  assetTypeEnum,
} from './types/accounting'

// =============================================================================
// List Response Types
// =============================================================================

export type ListExpensesResponse = ListResponse<Expense>
export type ListRecurringExpensesResponse = ListResponse<RecurringExpense>
export type ListInvestmentsResponse = ListResponse<Investment>
export type ListWithdrawalsResponse = ListResponse<Withdrawal>
export type ListAssetsResponse = ListResponse<Asset>

// =============================================================================
// List Params
// =============================================================================

export interface ListExpensesParams {
  page?: number
  pageSize?: number
  orderBy?: Array<string>
  from?: string // YYYY-MM-DD
  to?: string // YYYY-MM-DD
  category?: ExpenseCategory
}

export interface ListRecurringExpensesParams {
  page?: number
  pageSize?: number
  orderBy?: Array<string>
}

export interface ListInvestmentsParams {
  page?: number
  pageSize?: number
  orderBy?: Array<string>
}

export interface ListWithdrawalsParams {
  page?: number
  pageSize?: number
  orderBy?: Array<string>
}

export interface ListAssetsParams {
  page?: number
  pageSize?: number
  orderBy?: Array<string>
}

export interface SummaryParams {
  from?: string // YYYY-MM-DD
  to?: string // YYYY-MM-DD
}

export interface RecentActivitiesParams {
  limit?: number // Default: 10, Max: 50
}

// =============================================================================
// URL Builder Helpers
// =============================================================================

function buildSearchParams(
  params: Record<string, string | number | Array<string> | undefined | null>,
): string {
  const searchParams = new URLSearchParams()

  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null) continue

    if (Array.isArray(value)) {
      value.forEach((v) => searchParams.append(key, v))
    } else {
      searchParams.set(key, String(value))
    }
  }

  const query = searchParams.toString()
  return query ? `?${query}` : ''
}

// =============================================================================
// Accounting API Client
// =============================================================================

export const accountingApi = {
  // ---------------------------------------------------------------------------
  // Summary
  // ---------------------------------------------------------------------------

  /**
   * Get accounting summary with "Safe to Draw" amount
   * Optional date range filters for expense totals
   */
  async getSummary(
    businessDescriptor: string,
    params?: SummaryParams,
  ): Promise<AccountingSummary> {
    const query = buildSearchParams({
      from: params?.from,
      to: params?.to,
    })
    return get<AccountingSummary>(
      `v1/businesses/${businessDescriptor}/accounting/summary${query}`,
    )
  },

  // ---------------------------------------------------------------------------
  // Recent Activities
  // ---------------------------------------------------------------------------

  /**
   * Get recent accounting activities (expenses, investments, withdrawals)
   * Returns a unified list sorted by date descending
   */
  async getRecentActivities(
    businessDescriptor: string,
    params?: RecentActivitiesParams,
  ): Promise<RecentActivitiesResponse> {
    const query = buildSearchParams({
      limit: params?.limit,
    })
    return get<RecentActivitiesResponse>(
      `v1/businesses/${businessDescriptor}/accounting/recent-activities${query}`,
    )
  },

  // ---------------------------------------------------------------------------
  // Expenses
  // ---------------------------------------------------------------------------

  /**
   * List expenses for a business
   * Default sort: -occurredOn (most recent first)
   */
  async listExpenses(
    businessDescriptor: string,
    params?: ListExpensesParams,
  ): Promise<ListExpensesResponse> {
    const query = buildSearchParams({
      page: params?.page,
      pageSize: params?.pageSize,
      orderBy: params?.orderBy,
      from: params?.from,
      to: params?.to,
      category: params?.category,
    })
    return get<ListExpensesResponse>(
      `v1/businesses/${businessDescriptor}/accounting/expenses${query}`,
    )
  },

  /**
   * Get expense by ID
   */
  async getExpense(
    businessDescriptor: string,
    expenseId: string,
  ): Promise<Expense> {
    return get<Expense>(
      `v1/businesses/${businessDescriptor}/accounting/expenses/${expenseId}`,
    )
  },

  /**
   * Create a new expense
   */
  async createExpense(
    businessDescriptor: string,
    data: CreateExpenseRequest,
  ): Promise<Expense> {
    return post<Expense>(
      `v1/businesses/${businessDescriptor}/accounting/expenses`,
      { json: data },
    )
  },

  /**
   * Update an expense
   */
  async updateExpense(
    businessDescriptor: string,
    expenseId: string,
    data: UpdateExpenseRequest,
  ): Promise<Expense> {
    return patch<Expense>(
      `v1/businesses/${businessDescriptor}/accounting/expenses/${expenseId}`,
      { json: data },
    )
  },

  /**
   * Delete an expense
   */
  async deleteExpense(
    businessDescriptor: string,
    expenseId: string,
  ): Promise<void> {
    return del(
      `v1/businesses/${businessDescriptor}/accounting/expenses/${expenseId}`,
    )
  },

  // ---------------------------------------------------------------------------
  // Recurring Expenses
  // ---------------------------------------------------------------------------

  /**
   * List recurring expense templates
   */
  async listRecurringExpenses(
    businessDescriptor: string,
    params?: ListRecurringExpensesParams,
  ): Promise<ListRecurringExpensesResponse> {
    const query = buildSearchParams({
      page: params?.page,
      pageSize: params?.pageSize,
      orderBy: params?.orderBy,
    })
    return get<ListRecurringExpensesResponse>(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses${query}`,
    )
  },

  /**
   * Get recurring expense by ID (includes expenses array)
   */
  async getRecurringExpense(
    businessDescriptor: string,
    recurringExpenseId: string,
  ): Promise<RecurringExpense> {
    return get<RecurringExpense>(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses/${recurringExpenseId}`,
    )
  },

  /**
   * Create a new recurring expense template
   * Optionally backfills historical expenses if autoCreateHistoricalExpenses=true
   */
  async createRecurringExpense(
    businessDescriptor: string,
    data: CreateRecurringExpenseRequest,
  ): Promise<RecurringExpense> {
    return post<RecurringExpense>(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses`,
      { json: data },
    )
  },

  /**
   * Update a recurring expense template
   */
  async updateRecurringExpense(
    businessDescriptor: string,
    recurringExpenseId: string,
    data: UpdateRecurringExpenseRequest,
  ): Promise<RecurringExpense> {
    return patch<RecurringExpense>(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses/${recurringExpenseId}`,
      { json: data },
    )
  },

  /**
   * Update recurring expense status (active/paused/ended/canceled)
   */
  async updateRecurringExpenseStatus(
    businessDescriptor: string,
    recurringExpenseId: string,
    data: UpdateRecurringExpenseStatusRequest,
  ): Promise<RecurringExpense> {
    return patch<RecurringExpense>(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses/${recurringExpenseId}/status`,
      { json: data },
    )
  },

  /**
   * Delete a recurring expense template
   */
  async deleteRecurringExpense(
    businessDescriptor: string,
    recurringExpenseId: string,
  ): Promise<void> {
    return del(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses/${recurringExpenseId}`,
    )
  },

  /**
   * Get occurrences (expense instances) for a recurring expense
   */
  async getRecurringExpenseOccurrences(
    businessDescriptor: string,
    recurringExpenseId: string,
  ): Promise<Array<Expense>> {
    return get<Array<Expense>>(
      `v1/businesses/${businessDescriptor}/accounting/recurring-expenses/${recurringExpenseId}/occurrences`,
    )
  },

  // ---------------------------------------------------------------------------
  // Investments
  // ---------------------------------------------------------------------------

  /**
   * List investments (owner capital injections)
   */
  async listInvestments(
    businessDescriptor: string,
    params?: ListInvestmentsParams,
  ): Promise<ListInvestmentsResponse> {
    const query = buildSearchParams({
      page: params?.page,
      pageSize: params?.pageSize,
      orderBy: params?.orderBy,
    })
    return get<ListInvestmentsResponse>(
      `v1/businesses/${businessDescriptor}/accounting/investments${query}`,
    )
  },

  /**
   * Get investment by ID
   */
  async getInvestment(
    businessDescriptor: string,
    investmentId: string,
  ): Promise<Investment> {
    return get<Investment>(
      `v1/businesses/${businessDescriptor}/accounting/investments/${investmentId}`,
    )
  },

  /**
   * Create a new investment
   */
  async createInvestment(
    businessDescriptor: string,
    data: CreateInvestmentRequest,
  ): Promise<Investment> {
    return post<Investment>(
      `v1/businesses/${businessDescriptor}/accounting/investments`,
      { json: data },
    )
  },

  /**
   * Update an investment
   */
  async updateInvestment(
    businessDescriptor: string,
    investmentId: string,
    data: UpdateInvestmentRequest,
  ): Promise<Investment> {
    return patch<Investment>(
      `v1/businesses/${businessDescriptor}/accounting/investments/${investmentId}`,
      { json: data },
    )
  },

  /**
   * Delete an investment
   */
  async deleteInvestment(
    businessDescriptor: string,
    investmentId: string,
  ): Promise<void> {
    return del(
      `v1/businesses/${businessDescriptor}/accounting/investments/${investmentId}`,
    )
  },

  // ---------------------------------------------------------------------------
  // Withdrawals
  // ---------------------------------------------------------------------------

  /**
   * List withdrawals (owner capital draws)
   */
  async listWithdrawals(
    businessDescriptor: string,
    params?: ListWithdrawalsParams,
  ): Promise<ListWithdrawalsResponse> {
    const query = buildSearchParams({
      page: params?.page,
      pageSize: params?.pageSize,
      orderBy: params?.orderBy,
    })
    return get<ListWithdrawalsResponse>(
      `v1/businesses/${businessDescriptor}/accounting/withdrawals${query}`,
    )
  },

  /**
   * Get withdrawal by ID
   */
  async getWithdrawal(
    businessDescriptor: string,
    withdrawalId: string,
  ): Promise<Withdrawal> {
    return get<Withdrawal>(
      `v1/businesses/${businessDescriptor}/accounting/withdrawals/${withdrawalId}`,
    )
  },

  /**
   * Create a new withdrawal
   */
  async createWithdrawal(
    businessDescriptor: string,
    data: CreateWithdrawalRequest,
  ): Promise<Withdrawal> {
    return post<Withdrawal>(
      `v1/businesses/${businessDescriptor}/accounting/withdrawals`,
      { json: data },
    )
  },

  /**
   * Update a withdrawal
   */
  async updateWithdrawal(
    businessDescriptor: string,
    withdrawalId: string,
    data: UpdateWithdrawalRequest,
  ): Promise<Withdrawal> {
    return patch<Withdrawal>(
      `v1/businesses/${businessDescriptor}/accounting/withdrawals/${withdrawalId}`,
      { json: data },
    )
  },

  /**
   * Delete a withdrawal
   */
  async deleteWithdrawal(
    businessDescriptor: string,
    withdrawalId: string,
  ): Promise<void> {
    return del(
      `v1/businesses/${businessDescriptor}/accounting/withdrawals/${withdrawalId}`,
    )
  },

  // ---------------------------------------------------------------------------
  // Assets
  // ---------------------------------------------------------------------------

  /**
   * List fixed assets
   */
  async listAssets(
    businessDescriptor: string,
    params?: ListAssetsParams,
  ): Promise<ListAssetsResponse> {
    const query = buildSearchParams({
      page: params?.page,
      pageSize: params?.pageSize,
      orderBy: params?.orderBy,
    })
    return get<ListAssetsResponse>(
      `v1/businesses/${businessDescriptor}/accounting/assets${query}`,
    )
  },

  /**
   * Get asset by ID
   */
  async getAsset(businessDescriptor: string, assetId: string): Promise<Asset> {
    return get<Asset>(
      `v1/businesses/${businessDescriptor}/accounting/assets/${assetId}`,
    )
  },

  /**
   * Create a new asset
   */
  async createAsset(
    businessDescriptor: string,
    data: CreateAssetRequest,
  ): Promise<Asset> {
    return post<Asset>(
      `v1/businesses/${businessDescriptor}/accounting/assets`,
      { json: data },
    )
  },

  /**
   * Update an asset
   */
  async updateAsset(
    businessDescriptor: string,
    assetId: string,
    data: UpdateAssetRequest,
  ): Promise<Asset> {
    return patch<Asset>(
      `v1/businesses/${businessDescriptor}/accounting/assets/${assetId}`,
      { json: data },
    )
  },

  /**
   * Delete an asset
   */
  async deleteAsset(
    businessDescriptor: string,
    assetId: string,
  ): Promise<void> {
    return del(
      `v1/businesses/${businessDescriptor}/accounting/assets/${assetId}`,
    )
  },
}

// =============================================================================
// Query Options Factory (TanStack Query)
// =============================================================================

export const accountingQueries = {
  all: ['accounting'] as const,

  // ---------------------------------------------------------------------------
  // Summary
  // ---------------------------------------------------------------------------
  summaries: () => [...accountingQueries.all, 'summary'] as const,
  summary: (businessDescriptor: string, params?: SummaryParams) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.summaries(),
        businessDescriptor,
        params,
      ] as const,
      queryFn: () => accountingApi.getSummary(businessDescriptor, params),
      staleTime: STALE_TIME.ONE_MINUTE,
    }),

  // ---------------------------------------------------------------------------
  // Recent Activities
  // ---------------------------------------------------------------------------
  recentActivitiesKey: () =>
    [...accountingQueries.all, 'recent-activities'] as const,
  recentActivities: (
    businessDescriptor: string,
    params?: RecentActivitiesParams,
  ) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.recentActivitiesKey(),
        businessDescriptor,
        params,
      ] as const,
      queryFn: () =>
        accountingApi.getRecentActivities(businessDescriptor, params),
      staleTime: STALE_TIME.THIRTY_SECONDS,
    }),

  // ---------------------------------------------------------------------------
  // Expenses
  // ---------------------------------------------------------------------------
  expenses: () => [...accountingQueries.all, 'expenses'] as const,
  expenseList: (businessDescriptor: string, params?: ListExpensesParams) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.expenses(),
        'list',
        businessDescriptor,
        params,
      ] as const,
      queryFn: () => accountingApi.listExpenses(businessDescriptor, params),
      staleTime: STALE_TIME.THIRTY_SECONDS,
      placeholderData: keepPreviousData,
    }),
  expenseDetails: () => [...accountingQueries.expenses(), 'detail'] as const,
  expenseDetail: (businessDescriptor: string, expenseId: string) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.expenseDetails(),
        businessDescriptor,
        expenseId,
      ] as const,
      queryFn: () => accountingApi.getExpense(businessDescriptor, expenseId),
      staleTime: STALE_TIME.THIRTY_SECONDS,
    }),

  // ---------------------------------------------------------------------------
  // Recurring Expenses
  // ---------------------------------------------------------------------------
  recurringExpenses: () =>
    [...accountingQueries.all, 'recurring-expenses'] as const,
  recurringExpenseList: (
    businessDescriptor: string,
    params?: ListRecurringExpensesParams,
  ) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.recurringExpenses(),
        'list',
        businessDescriptor,
        params,
      ] as const,
      queryFn: () =>
        accountingApi.listRecurringExpenses(businessDescriptor, params),
      staleTime: STALE_TIME.ONE_MINUTE,
      placeholderData: keepPreviousData,
    }),
  recurringExpenseDetails: () =>
    [...accountingQueries.recurringExpenses(), 'detail'] as const,
  recurringExpenseDetail: (
    businessDescriptor: string,
    recurringExpenseId: string,
  ) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.recurringExpenseDetails(),
        businessDescriptor,
        recurringExpenseId,
      ] as const,
      queryFn: () =>
        accountingApi.getRecurringExpense(
          businessDescriptor,
          recurringExpenseId,
        ),
      staleTime: STALE_TIME.ONE_MINUTE,
    }),
  recurringExpenseOccurrences: (
    businessDescriptor: string,
    recurringExpenseId: string,
  ) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.recurringExpenses(),
        'occurrences',
        businessDescriptor,
        recurringExpenseId,
      ] as const,
      queryFn: () =>
        accountingApi.getRecurringExpenseOccurrences(
          businessDescriptor,
          recurringExpenseId,
        ),
      staleTime: STALE_TIME.ONE_MINUTE,
    }),

  // ---------------------------------------------------------------------------
  // Investments
  // ---------------------------------------------------------------------------
  investments: () => [...accountingQueries.all, 'investments'] as const,
  investmentList: (
    businessDescriptor: string,
    params?: ListInvestmentsParams,
  ) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.investments(),
        'list',
        businessDescriptor,
        params,
      ] as const,
      queryFn: () => accountingApi.listInvestments(businessDescriptor, params),
      staleTime: STALE_TIME.ONE_MINUTE,
      placeholderData: keepPreviousData,
    }),
  investmentDetails: () =>
    [...accountingQueries.investments(), 'detail'] as const,
  investmentDetail: (businessDescriptor: string, investmentId: string) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.investmentDetails(),
        businessDescriptor,
        investmentId,
      ] as const,
      queryFn: () =>
        accountingApi.getInvestment(businessDescriptor, investmentId),
      staleTime: STALE_TIME.ONE_MINUTE,
    }),

  // ---------------------------------------------------------------------------
  // Withdrawals
  // ---------------------------------------------------------------------------
  withdrawals: () => [...accountingQueries.all, 'withdrawals'] as const,
  withdrawalList: (
    businessDescriptor: string,
    params?: ListWithdrawalsParams,
  ) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.withdrawals(),
        'list',
        businessDescriptor,
        params,
      ] as const,
      queryFn: () => accountingApi.listWithdrawals(businessDescriptor, params),
      staleTime: STALE_TIME.ONE_MINUTE,
      placeholderData: keepPreviousData,
    }),
  withdrawalDetails: () =>
    [...accountingQueries.withdrawals(), 'detail'] as const,
  withdrawalDetail: (businessDescriptor: string, withdrawalId: string) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.withdrawalDetails(),
        businessDescriptor,
        withdrawalId,
      ] as const,
      queryFn: () =>
        accountingApi.getWithdrawal(businessDescriptor, withdrawalId),
      staleTime: STALE_TIME.ONE_MINUTE,
    }),

  // ---------------------------------------------------------------------------
  // Assets
  // ---------------------------------------------------------------------------
  assets: () => [...accountingQueries.all, 'assets'] as const,
  assetList: (businessDescriptor: string, params?: ListAssetsParams) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.assets(),
        'list',
        businessDescriptor,
        params,
      ] as const,
      queryFn: () => accountingApi.listAssets(businessDescriptor, params),
      staleTime: STALE_TIME.ONE_MINUTE,
      placeholderData: keepPreviousData,
    }),
  assetDetails: () => [...accountingQueries.assets(), 'detail'] as const,
  assetDetail: (businessDescriptor: string, assetId: string) =>
    queryOptions({
      queryKey: [
        ...accountingQueries.assetDetails(),
        businessDescriptor,
        assetId,
      ] as const,
      queryFn: () => accountingApi.getAsset(businessDescriptor, assetId),
      staleTime: STALE_TIME.ONE_MINUTE,
    }),
}

// =============================================================================
// React Query Hooks
// =============================================================================

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

export function useAccountingSummaryQuery(
  businessDescriptor: string,
  params?: SummaryParams,
) {
  return useQuery(accountingQueries.summary(businessDescriptor, params))
}

// ---------------------------------------------------------------------------
// Recent Activities
// ---------------------------------------------------------------------------

export function useRecentActivitiesQuery(
  businessDescriptor: string,
  params?: RecentActivitiesParams,
) {
  return useQuery(
    accountingQueries.recentActivities(businessDescriptor, params),
  )
}

// ---------------------------------------------------------------------------
// Expenses
// ---------------------------------------------------------------------------

export function useExpensesQuery(
  businessDescriptor: string,
  params?: ListExpensesParams,
) {
  return useQuery(accountingQueries.expenseList(businessDescriptor, params))
}

export function useExpenseQuery(businessDescriptor: string, expenseId: string) {
  return useQuery(
    accountingQueries.expenseDetail(businessDescriptor, expenseId),
  )
}

export function useCreateExpenseMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Expense, Error, CreateExpenseRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateExpenseRequest) =>
      accountingApi.createExpense(businessDescriptor, data),
    ...options,
  })
}

export function useUpdateExpenseMutation(
  businessDescriptor: string,
  expenseId: string,
  options?: UseMutationOptions<Expense, Error, UpdateExpenseRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateExpenseRequest) =>
      accountingApi.updateExpense(businessDescriptor, expenseId, data),
    ...options,
  })
}

export function useDeleteExpenseMutation(
  businessDescriptor: string,
  expenseId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () =>
      accountingApi.deleteExpense(businessDescriptor, expenseId),
    ...options,
  })
}

// ---------------------------------------------------------------------------
// Recurring Expenses
// ---------------------------------------------------------------------------

export function useRecurringExpensesQuery(
  businessDescriptor: string,
  params?: ListRecurringExpensesParams,
) {
  return useQuery(
    accountingQueries.recurringExpenseList(businessDescriptor, params),
  )
}

export function useRecurringExpenseQuery(
  businessDescriptor: string,
  recurringExpenseId: string,
  options?: { enabled?: boolean },
) {
  return useQuery({
    ...accountingQueries.recurringExpenseDetail(
      businessDescriptor,
      recurringExpenseId,
    ),
    ...options,
  })
}

export function useCreateRecurringExpenseMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<
    RecurringExpense,
    Error,
    CreateRecurringExpenseRequest
  >,
) {
  return useMutation({
    mutationFn: (data: CreateRecurringExpenseRequest) =>
      accountingApi.createRecurringExpense(businessDescriptor, data),
    ...options,
  })
}

export function useUpdateRecurringExpenseMutation(
  businessDescriptor: string,
  recurringExpenseId: string,
  options?: UseMutationOptions<
    RecurringExpense,
    Error,
    UpdateRecurringExpenseRequest
  >,
) {
  return useMutation({
    mutationFn: (data: UpdateRecurringExpenseRequest) =>
      accountingApi.updateRecurringExpense(
        businessDescriptor,
        recurringExpenseId,
        data,
      ),
    ...options,
  })
}

export function useUpdateRecurringExpenseStatusMutation(
  businessDescriptor: string,
  recurringExpenseId: string,
  options?: UseMutationOptions<
    RecurringExpense,
    Error,
    UpdateRecurringExpenseStatusRequest
  >,
) {
  return useMutation({
    mutationFn: (data: UpdateRecurringExpenseStatusRequest) =>
      accountingApi.updateRecurringExpenseStatus(
        businessDescriptor,
        recurringExpenseId,
        data,
      ),
    ...options,
  })
}

export function useDeleteRecurringExpenseMutation(
  businessDescriptor: string,
  recurringExpenseId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () =>
      accountingApi.deleteRecurringExpense(
        businessDescriptor,
        recurringExpenseId,
      ),
    ...options,
  })
}

// ---------------------------------------------------------------------------
// Investments
// ---------------------------------------------------------------------------

export function useInvestmentsQuery(
  businessDescriptor: string,
  params?: ListInvestmentsParams,
) {
  return useQuery(accountingQueries.investmentList(businessDescriptor, params))
}

export function useInvestmentQuery(
  businessDescriptor: string,
  investmentId: string,
) {
  return useQuery(
    accountingQueries.investmentDetail(businessDescriptor, investmentId),
  )
}

export function useCreateInvestmentMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Investment, Error, CreateInvestmentRequest>,
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateInvestmentRequest) =>
      accountingApi.createInvestment(businessDescriptor, data),
    onSuccess: (data, variables, onMutateResult, context) => {
      // Invalidate investments list and summary
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.investments(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.summaries(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.recentActivitiesKey(),
      })
      options?.onSuccess?.(data, variables, onMutateResult, context)
    },
    ...options,
  })
}

export function useUpdateInvestmentMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<
    Investment,
    Error,
    { id: string; data: UpdateInvestmentRequest }
  >,
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateInvestmentRequest }) =>
      accountingApi.updateInvestment(businessDescriptor, id, data),
    onSuccess: (data, variables, onMutateResult, context) => {
      // Invalidate investments list and summary
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.investments(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.summaries(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.recentActivitiesKey(),
      })
      options?.onSuccess?.(data, variables, onMutateResult, context)
    },
    ...options,
  })
}

export function useDeleteInvestmentMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<void, Error, string>,
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) =>
      accountingApi.deleteInvestment(businessDescriptor, id),
    onSuccess: (data, variables, onMutateResult, context) => {
      // Invalidate investments list and summary
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.investments(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.summaries(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.recentActivitiesKey(),
      })
      options?.onSuccess?.(data, variables, onMutateResult, context)
    },
    ...options,
  })
}

// ---------------------------------------------------------------------------
// Withdrawals
// ---------------------------------------------------------------------------

export function useWithdrawalsQuery(
  businessDescriptor: string,
  params?: ListWithdrawalsParams,
) {
  return useQuery(accountingQueries.withdrawalList(businessDescriptor, params))
}

export function useWithdrawalQuery(
  businessDescriptor: string,
  withdrawalId: string,
) {
  return useQuery(
    accountingQueries.withdrawalDetail(businessDescriptor, withdrawalId),
  )
}

export function useCreateWithdrawalMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Withdrawal, Error, CreateWithdrawalRequest>,
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: CreateWithdrawalRequest) =>
      accountingApi.createWithdrawal(businessDescriptor, data),
    onSuccess: (data, variables, onMutateResult, context) => {
      // Invalidate withdrawals list and summary
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.withdrawals(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.summaries(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.recentActivitiesKey(),
      })
      options?.onSuccess?.(data, variables, onMutateResult, context)
    },
    ...options,
  })
}

export function useUpdateWithdrawalMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<
    Withdrawal,
    Error,
    { id: string; data: UpdateWithdrawalRequest }
  >,
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateWithdrawalRequest }) =>
      accountingApi.updateWithdrawal(businessDescriptor, id, data),
    onSuccess: (data, variables, onMutateResult, context) => {
      // Invalidate withdrawals list and summary
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.withdrawals(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.summaries(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.recentActivitiesKey(),
      })
      options?.onSuccess?.(data, variables, onMutateResult, context)
    },
    ...options,
  })
}

export function useDeleteWithdrawalMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<void, Error, string>,
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) =>
      accountingApi.deleteWithdrawal(businessDescriptor, id),
    onSuccess: (data, variables, onMutateResult, context) => {
      // Invalidate withdrawals list and summary
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.withdrawals(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.summaries(),
      })
      void queryClient.invalidateQueries({
        queryKey: accountingQueries.recentActivitiesKey(),
      })
      options?.onSuccess?.(data, variables, onMutateResult, context)
    },
    ...options,
  })
}

// ---------------------------------------------------------------------------
// Assets
// ---------------------------------------------------------------------------

export function useAssetsQuery(
  businessDescriptor: string,
  params?: ListAssetsParams,
) {
  return useQuery(accountingQueries.assetList(businessDescriptor, params))
}

export function useAssetQuery(businessDescriptor: string, assetId: string) {
  return useQuery(accountingQueries.assetDetail(businessDescriptor, assetId))
}

export function useCreateAssetMutation(
  businessDescriptor: string,
  options?: UseMutationOptions<Asset, Error, CreateAssetRequest>,
) {
  return useMutation({
    mutationFn: (data: CreateAssetRequest) =>
      accountingApi.createAsset(businessDescriptor, data),
    ...options,
  })
}

export function useUpdateAssetMutation(
  businessDescriptor: string,
  assetId: string,
  options?: UseMutationOptions<Asset, Error, UpdateAssetRequest>,
) {
  return useMutation({
    mutationFn: (data: UpdateAssetRequest) =>
      accountingApi.updateAsset(businessDescriptor, assetId, data),
    ...options,
  })
}

export function useDeleteAssetMutation(
  businessDescriptor: string,
  assetId: string,
  options?: UseMutationOptions<void, Error, void>,
) {
  return useMutation({
    mutationFn: () => accountingApi.deleteAsset(businessDescriptor, assetId),
    ...options,
  })
}
