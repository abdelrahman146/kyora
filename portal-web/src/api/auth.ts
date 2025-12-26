import apiClient from "./client";
import type { LoginRequest, LoginResponse } from "./types";

export const authApi = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    return apiClient.post("auth/login", { json: credentials }).json();
  },

  logout: async (): Promise<void> => {
    return apiClient.post("auth/logout").json();
  },

  refresh: async (): Promise<{ accessToken: string }> => {
    return apiClient.post("auth/refresh").json();
  },

  getCurrentUser: async () => {
    return apiClient.get("auth/me").json();
  },
};
