/**
 * Re-export authentication utilities for convenience
 * This maintains backward compatibility with existing imports
 */
export { AuthProvider } from "./contexts/AuthContext";
export { useAuth } from "./hooks/useAuth";
export type { AuthContextType } from "./contexts/auth/AuthContext";
