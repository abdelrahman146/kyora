---
description: Frontend HTTP client - Ky configuration, TanStack Query patterns, error handling, retry/timeout (reusable across portal-web, storefront-web)
applyTo: "portal-web/**,storefront-web/**"
---

# HTTP Client & TanStack Query

Ky-based HTTP client with TanStack Query integration.

**Cross-refs:**

- Architecture: `./architecture.instructions.md`
- Forms: `./forms.instructions.md` (form submission)
- i18n: `./i18n.instructions.md` (error translation)

---

## 1. Ky Client Setup

### Basic Configuration

```typescript
import ky from "ky";

export const apiClient = ky.create({
  prefixUrl: import.meta.env.VITE_API_URL || "http://localhost:8080",
  timeout: 30000, // 30 seconds
  retry: {
    limit: 2,
    methods: ["get", "put", "head", "delete", "options", "trace"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 3000, // Max 3 seconds between retries
  },
  headers: {
    "Content-Type": "application/json",
  },
});
```

### With Auto-Refresh

```typescript
export const apiClient = ky.create({
  hooks: {
    beforeRequest: [
      (request) => {
        const token = getAccessToken();
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        if (response.status === 401) {
          const newToken = await refreshAccessToken();
          if (newToken) {
            request.headers.set("Authorization", `Bearer ${newToken}`);
            return ky(request);
          }
          clearTokens();
          window.location.href = "/auth/login";
        }
      },
    ],
    beforeError: [
      async (error) => {
        const { response } = error;
        if (response) {
          try {
            const body = await response.json();
            error.message = body.detail || body.message || error.message;
          } catch {
            // Response not JSON
          }
        }
        return error;
      },
    ],
  },
});
```

---

## 2. Typed Wrappers

```typescript
// src/api/client.ts
import { HTTPError } from "ky";

export async function get<T>(url: string, options?: Options): Promise<T> {
  return apiClient.get(url, options).json<T>();
}

export async function post<T>(
  url: string,
  json: unknown,
  options?: Options,
): Promise<T> {
  return apiClient.post(url, { ...options, json }).json<T>();
}

export async function postVoid(
  url: string,
  json: unknown,
  options?: Options,
): Promise<void> {
  await apiClient.post(url, { ...options, json });
}

export async function put<T>(
  url: string,
  json: unknown,
  options?: Options,
): Promise<T> {
  return apiClient.put(url, { ...options, json }).json<T>();
}

export async function patch<T>(
  url: string,
  json: unknown,
  options?: Options,
): Promise<T> {
  return apiClient.patch(url, { ...options, json }).json<T>();
}

export async function del<T>(url: string, options?: Options): Promise<T> {
  return apiClient.delete(url, options).json<T>();
}

export async function delVoid(url: string, options?: Options): Promise<void> {
  await apiClient.delete(url, options);
}
```

---

## 3. TanStack Query Integration

### API Module Pattern

```typescript
// src/api/customer.ts
import { get, post, patch, del } from "./client";
import type { Customer, CreateCustomerRequest } from "./types/customer";

export const customerApi = {
  async list(params?: { page?: number }): Promise<Customer[]> {
    return get<Customer[]>("v1/customers", { searchParams: params });
  },

  async get(id: string): Promise<Customer> {
    return get<Customer>(`v1/customers/${id}`);
  },

  async create(data: CreateCustomerRequest): Promise<Customer> {
    return post<Customer>("v1/customers", data);
  },

  async update(id: string, data: Partial<Customer>): Promise<Customer> {
    return patch<Customer>(`v1/customers/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    return del<void>(`v1/customers/${id}`);
  },
};
```

### Query Options Factory

```typescript
// src/api/customer.ts (continued)
export const customerQueries = {
  all: (params?: { page?: number }) => ({
    queryKey: ["customers", params] as const,
    queryFn: () => customerApi.list(params),
    staleTime: 1000 * 60, // 1 minute
  }),

  detail: (id: string) => ({
    queryKey: ["customers", id] as const,
    queryFn: () => customerApi.get(id),
    staleTime: 1000 * 60 * 5, // 5 minutes
  }),
};
```

### Usage in Components

```tsx
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { customerApi, customerQueries } from "@/api/customer";

function CustomersPage() {
  const { data, isLoading, error } = useQuery(customerQueries.all({ page: 1 }));

  if (error) return <ErrorMessage error={error} />;
  if (isLoading) return <Skeleton />;

  return <CustomersList customers={data} />;
}
```

---

## 4. Mutations

### Basic Mutation

```tsx
const mutation = useMutation({
  mutationFn: customerApi.create,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["customers"] });
    toast.success(t("customer.created"));
  },
});
```

### With Error Handling

```tsx
const mutation = useMutation({
  mutationFn: customerApi.create,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["customers"] });
    toast.success(t("customer.created"));
  },
  onError: async (error) => {
    const translated = await translateErrorAsync(error, t);
    toast.error(translated);
  },
});
```

### With Optimistic Update

```tsx
const mutation = useMutation({
  mutationFn: customerApi.update,
  onMutate: async (updatedCustomer) => {
    await queryClient.cancelQueries({
      queryKey: ["customers", updatedCustomer.id],
    });

    const previous = queryClient.getQueryData([
      "customers",
      updatedCustomer.id,
    ]);

    queryClient.setQueryData(
      ["customers", updatedCustomer.id],
      updatedCustomer,
    );

    return { previous };
  },
  onError: (err, variables, context) => {
    if (context?.previous) {
      queryClient.setQueryData(["customers", variables.id], context.previous);
    }
  },
  onSettled: (_, __, variables) => {
    queryClient.invalidateQueries({ queryKey: ["customers", variables.id] });
  },
});
```

---

## 5. Error Handling

### Global Error Handler

```typescript
// src/main.tsx
const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: async (error, query) => {
      if (query.meta?.errorToast === "off") return;

      const translated = await translateErrorAsync(error, t);
      toast.error(translated);
    },
  }),
  mutationCache: new MutationCache({
    onError: async (error, _variables, _context, mutation) => {
      if (mutation.meta?.errorToast === "off") return;

      const translated = await translateErrorAsync(error, t);
      toast.error(translated);
    },
  }),
});
```

### Opt-Out from Global Handler

```tsx
const mutation = useMutation({
  mutationFn: customerApi.create,
  meta: { errorToast: "off" }, // Suppress global toast
  onError: async (error) => {
    const translated = await translateErrorAsync(error, t);
    setFormError(translated); // Show inline instead
  },
});
```

### Error Translation

```typescript
import { translateErrorAsync } from "@/lib/translateError";

try {
  await api.login(credentials);
} catch (error) {
  const translated = await translateErrorAsync(error, t);
  toast.error(translated); // Shows translated message
}
```

---

## 6. Retry Logic

### Custom Retry

```typescript
export const apiClient = ky.create({
  retry: {
    limit: 3,
    methods: ["get", "put", "head", "delete"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],

    delay: (attemptCount) => {
      return Math.min(1000 * 2 ** attemptCount, 10000); // Max 10s
    },

    shouldRetry: ({ error, retryCount }) => {
      if (error instanceof HTTPError) {
        const status = error.response.status;

        if (status === 429 && retryCount <= 2) return true;
        if (status >= 400 && status < 500) return false; // Don't retry 4xx
      }

      return undefined; // Use default
    },
  },
});
```

### TanStack Query Retry

```typescript
const { data } = useQuery({
  queryKey: ["customers"],
  queryFn: customerApi.list,
  retry: 2, // Retry failed queries 2 times
  retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
});
```

---

## 7. Timeout Management

### Per-Request Timeout

```typescript
const data = await apiClient
  .get("v1/large-report", {
    timeout: 120000, // 2 minutes
  })
  .json();
```

### Disable Timeout

```typescript
const stream = await apiClient.get("v1/stream", {
  timeout: false, // No timeout
});
```

---

## 8. Request Cancellation

```typescript
const controller = new AbortController();
const { signal } = controller;

const promise = apiClient.get("v1/data", { signal }).json();

// Cancel after 5 seconds
setTimeout(() => controller.abort(), 5000);

try {
  const data = await promise;
} catch (error) {
  if (error.name === "AbortError") {
    console.log("Request cancelled");
  }
}
```

### With TanStack Query

```tsx
function SearchComponent() {
  const [query, setQuery] = useState("");

  const { data, isLoading } = useQuery({
    queryKey: ["search", query],
    queryFn: ({ signal }) => api.search(query, { signal }),
    enabled: query.length > 2,
  });

  return <input onChange={(e) => setQuery(e.target.value)} />;
}
```

---

## 9. Cache Invalidation

### Invalidate All

```typescript
queryClient.invalidateQueries({ queryKey: ["customers"] });
```

### Invalidate Specific

```typescript
queryClient.invalidateQueries({ queryKey: ["customers", "123"] });
```

### Invalidate Multiple

```typescript
await Promise.all([
  queryClient.invalidateQueries({ queryKey: ["customers"] }),
  queryClient.invalidateQueries({ queryKey: ["orders"] }),
]);
```

### Prefetch

```typescript
await queryClient.prefetchQuery(customerQueries.all({ page: 1 }));
```

### Ensure Data

```typescript
// In route loader
export const Route = createFileRoute("/customers")({
  loader: async ({ context }) => {
    return context.queryClient.ensureQueryData(
      customerQueries.all({ page: 1 }),
    );
  },
});
```

---

## 10. Search Params

```typescript
// Object
const users = await get<User[]>("v1/users", {
  searchParams: {
    page: 1,
    limit: 20,
    status: "active",
  },
});
// => GET /v1/users?page=1&limit=20&status=active

// URLSearchParams
const params = new URLSearchParams();
params.set("q", "search query");
params.set("filter", "active");

const results = await get<SearchResult[]>("v1/search", {
  searchParams: params,
});
```

---

## 11. Form Data

### Sending multipart/form-data

```typescript
const formData = new FormData();
formData.append("name", "Product Name");
formData.append("image", fileInput.files[0]);

const product = await apiClient
  .post("v1/products", {
    body: formData,
    // Content-Type automatically set
  })
  .json<Product>();
```

### URL-Encoded

```typescript
const params = new URLSearchParams();
params.set("username", "user@example.com");
params.set("password", "secret");

const response = await apiClient
  .post("v1/auth/login", {
    body: params,
    // Content-Type automatically set to application/x-www-form-urlencoded
  })
  .json<AuthResponse>();
```

---

## 12. Progress Tracking

### Upload Progress

```typescript
const response = await apiClient.post("v1/upload", {
  body: formData,
  onUploadProgress: (progress, chunk) => {
    const percent = Math.round(progress.percent * 100);
    updateProgressBar(percent);
  },
});
```

### Download Progress

```typescript
const response = await apiClient.get("v1/large-file", {
  onDownloadProgress: (progress, chunk) => {
    const percent = Math.round(progress.percent * 100);
    updateProgressBar(percent);
  },
});
```

---

## 13. Common Patterns

### Pagination

```typescript
export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export async function fetchPaginated<T>(
  endpoint: string,
  page: number = 1,
  limit: number = 20,
): Promise<PaginatedResponse<T>> {
  return get<PaginatedResponse<T>>(endpoint, {
    searchParams: { page, limit },
  });
}
```

### File Download

```typescript
export async function downloadFile(
  url: string,
  filename: string,
): Promise<void> {
  const response = await apiClient.get(url);
  const blob = await response.blob();

  const downloadUrl = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = downloadUrl;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(downloadUrl);
}
```

### Debounced Search

```typescript
import { debounce } from "lodash-es";

let searchController: AbortController | null = null;

export const debouncedSearch = debounce(async (query: string) => {
  if (searchController) {
    searchController.abort();
  }

  searchController = new AbortController();

  try {
    const results = await get<SearchResult[]>("v1/search", {
      searchParams: { q: query },
      signal: searchController.signal,
    });
    return results;
  } catch (error) {
    if (error.name === "AbortError") return [];
    throw error;
  }
}, 300);
```

---

## 14. Performance

### Stale Time

```typescript
const { data } = useQuery({
  queryKey: ["countries"],
  queryFn: getCountries,
  staleTime: 1000 * 60 * 60, // 1 hour (rarely changes)
});
```

### Cache Time

```typescript
const { data } = useQuery({
  queryKey: ["customers", id],
  queryFn: () => customerApi.get(id),
  cacheTime: 1000 * 60 * 10, // Keep in cache for 10 minutes
});
```

### Parallel Queries

```tsx
const { data: customers } = useQuery(customerQueries.all());
const { data: orders } = useQuery(orderQueries.all());
const { data: products } = useQuery(productQueries.all());
```

---

## Agent Validation

Before completing HTTP task:

- ☑ Using typed wrappers (`get<T>`, `post<T>`, etc.)
- ☑ Query keys from centralized factory
- ☑ API modules follow pattern: `{resource}Api` + `{resource}Queries`
- ☑ Mutations invalidate related queries
- ☑ Global error handler configured
- ☑ Retry logic configured (2 retries for GET/PUT/DELETE)
- ☑ Timeout set (30s default, longer for reports)
- ☑ Error translation via `translateErrorAsync`
- ☑ Auth token in `Authorization` header
- ☑ 401 triggers refresh + retry

---

## Resources

- Ky Docs: https://github.com/sindresorhus/ky
- TanStack Query: https://tanstack.com/query/latest
- Implementation: `portal-web/src/api/client.ts`
