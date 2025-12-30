// @vitest-environment jsdom

import { beforeEach, describe, expect, it } from 'vitest'
import { authStore, initializeAuth } from './authStore'

describe('authStore.initializeAuth', () => {
  beforeEach(() => {
    // Ensure a consistent initial state.
    authStore.setState(() => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,
    }))

    // Ensure no refresh token exists.
    document.cookie = 'kyora_refresh_token=; Max-Age=0; path=/'
  })

  it('clears loading when no refresh token is present', async () => {
    await initializeAuth()

    expect(authStore.state.isLoading).toBe(false)
    expect(authStore.state.isAuthenticated).toBe(false)
    expect(authStore.state.user).toBe(null)
  })
})
