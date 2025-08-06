export interface PagingQuery {
  page: number
  size: number
}

export interface PagingResult<T> {
  items: T[]
  page: number
  size: number
  total: number
}

export type { ApiError, Problem } from '@/types/problem'
