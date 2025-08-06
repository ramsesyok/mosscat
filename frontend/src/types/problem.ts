import type { Problem as OpenApiProblem } from '@/api'

export interface Problem extends OpenApiProblem {}

export interface ApiError extends Error {
  status?: number
  problem?: Problem
}
