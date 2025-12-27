import { get, patch } from "./client";
import type { User } from "./types/auth";

/**
 * User API Client
 * Handles user profile and management operations
 */
export const userApi = {
  /**
   * Get current authenticated user profile
   * GET /v1/users/me
   */
  async getCurrentUser(): Promise<User> {
    return get<User>("v1/users/me");
  },

  /**
   * Update current user profile
   * PATCH /v1/users/me
   */
  async updateCurrentUser(data: {
    firstName?: string;
    lastName?: string;
  }): Promise<User> {
    return patch<User>("v1/users/me", { json: data });
  },
};
