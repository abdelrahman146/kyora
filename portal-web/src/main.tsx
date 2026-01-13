import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import {
  MutationCache,
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import { RouterProvider } from '@tanstack/react-router'
import i18n from 'i18next'
import './i18n/init'
import './styles.css'

// Register Chart.js components for tree-shaking
import {
  ArcElement,
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  TimeScale,
  Title,
  Tooltip,
} from 'chart.js'
import { getRouter } from './router'
import { showErrorFromException } from '@/lib/toast'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
  TimeScale,
)

function shouldIgnoreGlobalError(error: unknown): boolean {
  if (!error) return true

  if (error instanceof DOMException && error.name === 'AbortError') {
    return true
  }

  if (error instanceof Error) {
    const name = error.name.toLowerCase()
    if (name.includes('abort') || name.includes('cancel')) {
      return true
    }
  }

  return false
}

const queryErrorToastDeduper = new Map<string, number>()
const QUERY_ERROR_TOAST_DEDUPE_MS = 30_000

// Create QueryClient instance
const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, query) => {
      if (shouldIgnoreGlobalError(error)) return

      // Global is the default source of truth for backend errors.
      // To avoid toast spam for background refetches, dedupe by queryHash.
      // Allow opt-out via meta.
      const meta = query.meta as
        | undefined
        | {
            errorToast?: 'global' | 'off'
          }

      if (meta?.errorToast === 'off') return

      const now = Date.now()
      const last = queryErrorToastDeduper.get(query.queryHash) ?? 0
      if (now - last < QUERY_ERROR_TOAST_DEDUPE_MS) return
      queryErrorToastDeduper.set(query.queryHash, now)

      void showErrorFromException(error, i18n.t)
    },
  }),
  mutationCache: new MutationCache({
    onError: (error, _variables, _context, mutation) => {
      if (shouldIgnoreGlobalError(error)) return

      const meta = mutation.meta as
        | undefined
        | {
            errorToast?: 'global' | 'off'
          }

      if (meta?.errorToast === 'off') return

      void showErrorFromException(error, i18n.t)
    },
  }),
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60, // 1 minute default
      gcTime: 1000 * 60 * 5, // 5 minutes garbage collection
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

// Create router instance
const router = getRouter(queryClient)

// Register router for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const rootElement = document.getElementById('root')

if (!rootElement) {
  throw new Error('Root element not found')
}

createRoot(rootElement).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
