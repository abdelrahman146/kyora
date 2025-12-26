// API Type Definitions
// Generated from backend swagger.json or manually typed

export interface User {
  id: string;
  email: string;
  name: string;
  workspaceId: string;
  role: "admin" | "member";
}

export interface Business {
  id: string;
  name: string;
  slug: string;
  workspaceId: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

// Add more types as needed based on backend API
