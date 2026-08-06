import { buildAdminAuthHeaders, logoutAdmin, type AdminAccount } from './adminAuth'
import { buildApiUrl } from './apiBase'
import { parseApiResponse } from './apiError'

const ADMIN_GET_CACHE_TTL_MS = 5 * 60 * 1000
const adminGetCache = new Map<string, { expiresAt: number; promise: Promise<unknown> }>()

export async function parseJson<T>(response: Response): Promise<T> {
  try {
    return await parseApiResponse<T>(response, 'Falha na requisicao administrativa.')
  } catch (error) {
    if (response.status === 401) {
      logoutAdmin()
    }
    throw error
  }
}

export async function adminFetch(input: RequestInfo | URL, init: RequestInit = {}) {
  const headers = new Headers(init.headers)
  const authHeaders = buildAdminAuthHeaders()
  Object.entries(authHeaders).forEach(([key, value]) => headers.set(key, value))
  const resolvedInput =
    typeof input === 'string' ? buildApiUrl(input) : input instanceof URL ? buildApiUrl(input.toString()) : input

  return fetch(resolvedInput, {
    ...init,
    credentials: 'include',
    headers,
  })
}

export function cachedAdminJson<T>(path: string, ttlMs = ADMIN_GET_CACHE_TTL_MS): Promise<T> {
  const now = Date.now()
  const cached = adminGetCache.get(path)
  if (cached && cached.expiresAt > now) {
    return cached.promise as Promise<T>
  }

  const promise = adminFetch(path)
    .then((response) => parseJson<T>(response))
    .catch((error) => {
      adminGetCache.delete(path)
      throw error
    })
  adminGetCache.set(path, { expiresAt: now + ttlMs, promise })
  return promise
}

export function invalidateAdminGetCache(...pathPrefixes: string[]) {
  for (const key of adminGetCache.keys()) {
    if (pathPrefixes.some((prefix) => key.startsWith(prefix))) {
      adminGetCache.delete(key)
    }
  }
}

export type AdminMeResponse = AdminAccount
