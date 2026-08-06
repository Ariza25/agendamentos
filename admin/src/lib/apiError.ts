import { getApiErrorMessage, translateFieldErrors } from './apiMessages'

export type ApiErrorPayload = {
  error?: string
  code?: string
  fieldErrors?: Record<string, string>
}

export class ApiError extends Error {
  code: string
  status: number
  fieldErrors: Record<string, string>

  constructor(message: string, code = 'unknown_error', status = 0, fieldErrors: Record<string, string> = {}) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
    this.fieldErrors = fieldErrors
  }
}

export function buildApiError(
  payload: ApiErrorPayload | null | undefined,
  fallbackMessage: string,
  status = 0,
): ApiError {
  const code = payload?.code || 'unknown_error'
  return new ApiError(
    getApiErrorMessage(code, payload?.error || fallbackMessage),
    code,
    status,
    translateFieldErrors(payload?.fieldErrors),
  )
}

export async function parseApiResponse<T>(response: Response, fallbackMessage: string): Promise<T> {
  const payload = (await response.json().catch(() => null)) as (T & ApiErrorPayload) | null

  if (!response.ok) {
    throw buildApiError(payload ?? null, fallbackMessage, response.status)
  }

  return payload as T
}
