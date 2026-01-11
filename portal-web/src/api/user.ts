/**
 * User API Client
 *
 * Handles user profile and management operations.
 * Includes TanStack Query hooks for data fetching and mutations.
 */

import {
  queryOptions,
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { get, patch } from './client'
import type { User } from './types/auth'
import { STALE_TIME, queryKeys } from '@/lib/queryKeys'
import { showErrorToast, showSuccessToast } from '@/lib/toast'

export interface UpdateUserRequest {
  firstName?: string
  lastName?: string
}

/**
 * User API Client
 */
export const userApi = {
  /**
   * Get current authenticated user profile
   * GET /v1/users/me
   */
  async getCurrentUser(): Promise<User> {
    return get<User>('v1/users/me')
  },

  /**
   * Update current user profile
   * PATCH /v1/users/me
   */
  async updateCurrentUser(data: UpdateUserRequest): Promise<User> {
    return patch<User>('v1/users/me', { json: data })
  },
}

/**
 * Query Options Factories
 *
 * Co-locate query configuration (key + fn + staleTime) for type-safe reuse
 * in components, route loaders, and prefetching.
 */
export const userQueries = {
  /**
   * Query options for fetching current user profile
   */
  current: () =>
    queryOptions({
      queryKey: queryKeys.user.current(),
      queryFn: () => userApi.getCurrentUser(),
      staleTime: STALE_TIME.FIVE_MINUTES,
    }),
}

/**
 * Query Hooks
 */

/**
 * Query to fetch current user profile
 */
export function useCurrentUserQuery() {
  return useQuery(userQueries.current())
}

/**
 * Mutation Hooks
 */

/**
 * Mutation to update current user profile
 */
export function useUpdateUserMutation() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('common')

  return useMutation({
    mutationFn: (data: UpdateUserRequest) => userApi.updateCurrentUser(data),
    onSuccess: (updatedUser) => {
      // Update user in cache
      queryClient.setQueryData(queryKeys.user.current(), updatedUser)
      showSuccessToast(t('profile.updateSuccess'))
    },
    onError: () => {
      showErrorToast(t('profile.updateError'))
    },
  })
}
