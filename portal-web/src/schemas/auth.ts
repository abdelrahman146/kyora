import { z } from "zod";

/**
 * Login Form Validation Schema
 * Validates email and password fields for the login form
 */
export const loginSchema = z.object({
  email: z
    .string()
    .min(1, "validation.required")
    .pipe(z.email("validation.invalid_email")),
  password: z.string().min(1, "validation.required"),
});

export type LoginFormData = z.infer<typeof loginSchema>;
