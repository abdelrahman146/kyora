// Global TypeScript type definitions

export type Role = 'admin' | 'member'

export interface PaginatedResponse<T> {
  data: Array<T>
  total: number
  page: number
  perPage: number
  totalPages: number
}

export interface ProblemDetails {
  type: string
  title: string
  status: number
  detail: string
  instance: string
}
