---
description: Frontend testing - Vitest setup, testing-library patterns, E2E patterns, mock patterns (reusable across portal-web, storefront-web)
applyTo: "portal-web/**,storefront-web/**"
---

# Frontend Testing

Vitest + React Testing Library + Playwright patterns.

**Cross-refs:**

- Architecture: `./architecture.instructions.md`
- HTTP client: `./http-client.instructions.md` (mocking API calls)

---

## 1. Vitest Setup

### Configuration

```typescript
// vite.config.ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: ["./src/test/setup.ts"],
    globals: true,
    css: true,
  },
});
```

### Setup File

```typescript
// src/test/setup.ts
import "@testing-library/jest-dom";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";

afterEach(() => {
  cleanup();
});
```

---

## 2. Unit Tests

### Component Tests

```typescript
// src/components/atoms/Button.test.tsx
// @vitest-environment jsdom

import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from './Button';

describe('Button', () => {
  it('should render with label', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('should call onClick when clicked', async () => {
    const onClick = vi.fn();
    render(<Button onClick={onClick}>Click me</Button>);

    await userEvent.click(screen.getByText('Click me'));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('should be disabled when disabled prop is true', () => {
    render(<Button disabled>Click me</Button>);
    expect(screen.getByText('Click me')).toBeDisabled();
  });
});
```

### Store Tests

```typescript
// src/stores/authStore.test.ts
// @vitest-environment jsdom

import { describe, it, expect, beforeEach } from "vitest";
import { authStore, loginUser, logoutUser } from "./authStore";

describe("authStore", () => {
  beforeEach(() => {
    authStore.setState({ user: null, isAuthenticated: false });
  });

  it("should set user on login", () => {
    const mockUser = { id: "1", email: "test@example.com" };
    loginUser(mockUser);

    expect(authStore.state.user).toEqual(mockUser);
    expect(authStore.state.isAuthenticated).toBe(true);
  });

  it("should clear user on logout", () => {
    const mockUser = { id: "1", email: "test@example.com" };
    authStore.setState({ user: mockUser, isAuthenticated: true });

    logoutUser();

    expect(authStore.state.user).toBeNull();
    expect(authStore.state.isAuthenticated).toBe(false);
  });
});
```

### Utility Tests

```typescript
// src/lib/formatCurrency.test.ts
import { describe, it, expect } from "vitest";
import { formatCurrency } from "./formatCurrency";

describe("formatCurrency", () => {
  it("should format USD correctly", () => {
    expect(formatCurrency(1000, "USD")).toBe("$1,000.00");
  });

  it("should format AED correctly", () => {
    expect(formatCurrency(1000, "AED")).toBe("AED 1,000.00");
  });

  it("should handle zero", () => {
    expect(formatCurrency(0, "USD")).toBe("$0.00");
  });

  it("should handle negative values", () => {
    expect(formatCurrency(-500, "USD")).toBe("-$500.00");
  });
});
```

---

## 3. Integration Tests

### Form Tests

```typescript
// src/features/auth/components/LoginForm.test.tsx
// @vitest-environment jsdom

import { describe, it, expect, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LoginForm } from './LoginForm';

describe('LoginForm', () => {
  it('should show validation errors on invalid input', async () => {
    render(<LoginForm onSuccess={vi.fn()} />);

    await userEvent.click(screen.getByRole('button', { name: /login/i }));

    expect(await screen.findByText(/invalid email/i)).toBeInTheDocument();
  });

  it('should submit form with valid data', async () => {
    const onSuccess = vi.fn();
    render(<LoginForm onSuccess={onSuccess} />);

    await userEvent.type(screen.getByLabelText(/email/i), 'test@example.com');
    await userEvent.type(screen.getByLabelText(/password/i), 'password123');
    await userEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
    });
  });
});
```

### Query Tests

```typescript
// src/api/customer.test.tsx
// @vitest-environment jsdom

import { describe, it, expect } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useQuery } from '@tanstack/react-query';
import { customerQueries } from './customer';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return ({ children }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('customerQueries', () => {
  it('should fetch customers', async () => {
    const { result } = renderHook(
      () => useQuery(customerQueries.all()),
      { wrapper: createWrapper() }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toHaveLength(3);
  });
});
```

---

## 4. Mocking

### Mocking API Calls

```typescript
import { vi } from "vitest";
import * as api from "@/api/customer";

// Mock entire module
vi.mock("@/api/customer", () => ({
  customerApi: {
    list: vi.fn(() =>
      Promise.resolve([
        { id: "1", name: "Customer 1" },
        { id: "2", name: "Customer 2" },
      ]),
    ),
    create: vi.fn((data) => Promise.resolve({ id: "3", ...data })),
  },
}));

// Mock specific function
vi.spyOn(api, "list").mockResolvedValue([{ id: "1", name: "Customer 1" }]);
```

### Mocking Stores

```typescript
import { vi } from "vitest";
import * as authStore from "@/stores/authStore";

vi.spyOn(authStore, "authStore", "get").mockReturnValue({
  state: {
    user: { id: "1", email: "test@example.com" },
    isAuthenticated: true,
  },
  setState: vi.fn(),
  subscribe: vi.fn(),
});
```

### Mocking Router

```typescript
import { vi } from "vitest";

const mockNavigate = vi.fn();

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
  useParams: () => ({ businessDescriptor: "test-business" }),
  useSearch: () => ({ page: 1 }),
}));
```

---

## 5. Testing Patterns

### Async Behavior

```typescript
it('should show loading state', async () => {
  render(<CustomersList />);

  expect(screen.getByText(/loading/i)).toBeInTheDocument();

  await waitFor(() => {
    expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
  });
});
```

### User Interactions

```typescript
it('should toggle dropdown', async () => {
  render(<Dropdown />);

  const button = screen.getByRole('button');

  await userEvent.click(button);
  expect(screen.getByRole('menu')).toBeVisible();

  await userEvent.click(button);
  expect(screen.queryByRole('menu')).not.toBeInTheDocument();
});
```

### Error Handling

```typescript
it('should show error message on failure', async () => {
  vi.spyOn(api, 'list').mockRejectedValue(new Error('Network error'));

  render(<CustomersList />);

  await waitFor(() => {
    expect(screen.getByText(/failed to load/i)).toBeInTheDocument();
  });
});
```

---

## 6. E2E Tests (Playwright)

### Configuration

```typescript
// playwright.config.ts
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests/e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    baseURL: "http://localhost:3000",
    trace: "on-first-retry",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: "Mobile Safari",
      use: { ...devices["iPhone 12"] },
    },
  ],
  webServer: {
    command: "npm run dev",
    url: "http://localhost:3000",
    reuseExistingServer: !process.env.CI,
  },
});
```

### E2E Test Example

```typescript
// tests/e2e/login.spec.ts
import { test, expect } from "@playwright/test";

test.describe("Login Flow", () => {
  test("should login successfully", async ({ page }) => {
    await page.goto("/auth/login");

    await page.fill('input[name="email"]', "test@example.com");
    await page.fill('input[name="password"]', "password123");
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL("/");
    await expect(page.locator("text=Dashboard")).toBeVisible();
  });

  test("should show error on invalid credentials", async ({ page }) => {
    await page.goto("/auth/login");

    await page.fill('input[name="email"]', "wrong@example.com");
    await page.fill('input[name="password"]', "wrongpassword");
    await page.click('button[type="submit"]');

    await expect(page.locator("text=Invalid credentials")).toBeVisible();
  });
});
```

### Page Object Pattern

```typescript
// tests/e2e/pages/LoginPage.ts
export class LoginPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/auth/login");
  }

  async login(email: string, password: string) {
    await this.page.fill('input[name="email"]', email);
    await this.page.fill('input[name="password"]', password);
    await this.page.click('button[type="submit"]');
  }

  async getErrorMessage() {
    return this.page.locator('[role="alert"]').textContent();
  }
}

// Usage
test("should login", async ({ page }) => {
  const loginPage = new LoginPage(page);
  await loginPage.goto();
  await loginPage.login("test@example.com", "password123");

  await expect(page).toHaveURL("/");
});
```

---

## 7. Coverage

### Run with Coverage

```bash
npx vitest --coverage
```

### Coverage Config

```typescript
// vite.config.ts
export default defineConfig({
  test: {
    coverage: {
      provider: "v8",
      reporter: ["text", "html", "json"],
      exclude: [
        "node_modules/",
        "src/test/",
        "**/*.test.{ts,tsx}",
        "**/*.spec.{ts,tsx}",
      ],
      thresholds: {
        lines: 80,
        functions: 80,
        branches: 80,
        statements: 80,
      },
    },
  },
});
```

---

## 8. Best Practices

### Do's

- ✅ Test user behavior, not implementation
- ✅ Use `screen.getByRole()` over `getByTestId()`
- ✅ Mock external dependencies (API, stores)
- ✅ Test error states and edge cases
- ✅ Keep tests independent (no shared state)
- ✅ Use descriptive test names
- ✅ Arrange-Act-Assert pattern

### Don'ts

- ❌ Don't test implementation details
- ❌ Don't test third-party libraries
- ❌ Don't share state between tests
- ❌ Don't over-mock (keep tests realistic)
- ❌ Don't skip async cleanup
- ❌ Don't test styling (use E2E instead)

---

## Agent Validation

Before completing testing task:

- ☑ Test files co-located with source (`*.test.ts`)
- ☑ `// @vitest-environment jsdom` comment at top
- ☑ `beforeEach` used for test isolation
- ☑ API calls mocked (no real HTTP)
- ☑ Loading/error states tested
- ☑ User interactions tested with `userEvent`
- ☑ Async behavior uses `waitFor`
- ☑ E2E tests use Page Object pattern
- ☑ Test names descriptive (what it should do)

---

## Resources

- Vitest Docs: https://vitest.dev
- Testing Library: https://testing-library.com
- Playwright: https://playwright.dev
