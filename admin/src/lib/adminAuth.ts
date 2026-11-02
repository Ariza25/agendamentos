import { buildApiError } from './apiError'
import { buildApiUrl } from './apiBase'

export type AdminRole = 'owner' | 'admin' | 'receptionist' | 'professional'

export type AdminAccount = {
  id: string
  name: string
  email: string
  phone: string
  role?: AdminRole
  tenantId: string
  mfaEnabled: boolean
  createdAt: string
}

export type AdminSession = {
  user: AdminAccount
}

export type AdminLoginInput = {
  email: string
  password: string
  mfaCode?: string
}

export type AdminAuthError = {
  message: string
  code: string
  status: number
  fieldErrors?: Record<string, string>
}

export const ADMIN_ROLES = new Set<AdminRole>(['owner', 'admin', 'receptionist', 'professional'])

const ADMIN_USER_STORAGE_KEY = 'erpagendamentos:admin:user'
const ADMIN_CSRF_STORAGE_KEY = 'erpagendamentos:admin:csrf'
const ADMIN_CLIENT_ID_STORAGE_KEY = 'erpagendamentos:admin:client-id'
const ADMIN_AUTH_EVENT = 'erpagendamentos:admin:auth:changed'

export async function loginAdmin(input: AdminLoginInput): Promise<AdminAccount> {
  const response = await fetch(buildApiUrl('/api/v1/auth/login'), {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...buildAdminAuthHeaders() },
    body: JSON.stringify({ ...input, audience: 'admin' }),
  })

  const payload = (await response.json().catch(() => null)) as
    | ({ error?: string; fieldErrors?: Record<string, string>; code?: string; csrfToken?: string } & Partial<AdminSession>)
    | null

  if (!response.ok || !payload?.csrfToken || !payload.user) {
    throw buildApiError(payload ?? null, 'Falha ao autenticar.', response.status)
  }

  if (!isAdminRole(payload.user.role)) {
    throw buildApiError(
      { error: 'Este usuario nao tem acesso administrativo.', code: 'forbidden' },
      'Este usuario nao tem acesso administrativo.',
    )
  }

  persistAdminSession(payload.user, payload.csrfToken)
  return payload.user
}

export type SignupInput = {
  businessName?: string
  slug?: string
  adminName?: string
  ownerName: string
  email: string
  password: string
  phone: string
  segment?: string
}

export async function signupAdmin(input: SignupInput): Promise<AdminAccount> {
  const requestBody = {
    ...input,
    adminName: input.adminName || input.ownerName,
    ownerName: input.ownerName || input.adminName,
  }
  const response = await fetch(buildApiUrl('/api/v1/auth/signup'), {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...buildAdminAuthHeaders() },
    body: JSON.stringify(requestBody),
  })

  const resData = (await response.json().catch(() => null)) as
    | ({
        error?: string
        fieldErrors?: Record<string, string>
        code?: string
        csrfToken?: string
        user?: Partial<AdminAccount>
        session?: { user?: Partial<AdminAccount> }
      })
    | null

  const userPayload = resData?.user || resData?.session?.user

  if (!response.ok || !resData?.csrfToken || !userPayload) {
    throw buildApiError(resData ?? null, 'Falha ao cadastrar estabelecimento.', response.status)
  }

  const user = normalizeAdminAccount(userPayload)
  persistAdminSession(user, resData.csrfToken)
  return user
}

export async function refreshAdminAccount() {
  const response = await fetch(buildApiUrl('/api/v1/auth/me'), {
    credentials: 'include',
    headers: buildAdminAuthHeaders(),
  })

  if (!response.ok) {
    logoutAdmin()
    return null
  }

  const user = normalizeAdminAccount((await response.json()) as Partial<AdminAccount>)
  if (!isAdminRole(user.role)) {
    logoutAdmin()
    return null
  }

  persistAdminUser(user)
  return user
}

export function logoutAdmin() {
  if (typeof window === 'undefined') {
    return
  }

  void fetch(buildApiUrl('/api/v1/auth/logout'), {
    method: 'POST',
    credentials: 'include',
    headers: buildAdminAuthHeaders(),
  }).catch(() => {})
  window.localStorage.removeItem(ADMIN_CSRF_STORAGE_KEY)
  window.localStorage.removeItem(ADMIN_USER_STORAGE_KEY)
  dispatchAdminAuthChange()
}

export function getAdminSession(): AdminSession | null {
  const user = getAdminAccount()
  if (!user || !isAdminRole(user.role)) {
    return null
  }

  return { user }
}

export function getAdminAccount(): AdminAccount | null {
  if (typeof window === 'undefined') {
    return null
  }

  const raw = window.localStorage.getItem(ADMIN_USER_STORAGE_KEY)
  if (!raw) {
    return null
  }

  try {
    return normalizeAdminAccount(JSON.parse(raw) as Partial<AdminAccount>)
  } catch {
    return null
  }
}

export function buildAdminAuthHeaders(): Record<string, string> {
  const csrfToken = getAdminCSRFToken()
  return {
    ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
    'X-Client-ID': getOrCreateAdminClientID(),
  }
}

export function subscribeToAdminAuthChanges(callback: () => void) {
  const handler = () => callback()
  window.addEventListener(ADMIN_AUTH_EVENT, handler)
  window.addEventListener('storage', handler)
  return () => {
    window.removeEventListener(ADMIN_AUTH_EVENT, handler)
    window.removeEventListener('storage', handler)
  }
}

function persistAdminSession(user: AdminAccount, csrfToken: string) {
  if (typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(ADMIN_CSRF_STORAGE_KEY, csrfToken)
  persistAdminUser(user)
}

function getAdminCSRFToken() {
  if (typeof window === 'undefined') {
    return ''
  }

  return window.localStorage.getItem(ADMIN_CSRF_STORAGE_KEY)?.trim() ?? ''
}

function getOrCreateAdminClientID() {
  if (typeof window === 'undefined') return 'server-render-admin'
  const current = window.localStorage.getItem(ADMIN_CLIENT_ID_STORAGE_KEY)?.trim()
  if (current) return current
  const clientID =
    typeof crypto.randomUUID === 'function'
      ? crypto.randomUUID()
      : `admin-${Date.now()}-${Math.random().toString(16).slice(2)}`
  window.localStorage.setItem(ADMIN_CLIENT_ID_STORAGE_KEY, clientID)
  return clientID
}

function persistAdminUser(user: AdminAccount) {
  if (typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(ADMIN_USER_STORAGE_KEY, JSON.stringify(normalizeAdminAccount(user)))
  dispatchAdminAuthChange()
}

function dispatchAdminAuthChange() {
  if (typeof window === 'undefined') {
    return
  }

  window.dispatchEvent(new Event(ADMIN_AUTH_EVENT))
}

export async function requestAdminPasswordReset(email: string): Promise<void> {
  const response = await fetch(buildApiUrl('/api/v1/auth/password-reset/request'), {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...buildAdminAuthHeaders() },
    body: JSON.stringify({ email, audience: 'admin' }),
  })
  const payload = (await response.json().catch(() => null)) as ({ error?: string; code?: string } | null) | null

  if (!response.ok) {
    throw buildApiError(payload, 'Nao foi possivel solicitar a recuperacao de senha.', response.status)
  }
}

export async function confirmAdminPasswordReset(input: {
  token: string
  newPassword: string
}): Promise<void> {
  const response = await fetch(buildApiUrl('/api/v1/auth/password-reset/confirm'), {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...buildAdminAuthHeaders() },
    body: JSON.stringify(input),
  })
  const payload = (await response.json().catch(() => null)) as ({ error?: string; code?: string } | null) | null

  if (!response.ok) {
    throw buildApiError(payload, 'Nao foi possivel definir a nova senha.', response.status)
  }
}

export function isAdminRole(role?: string): role is AdminRole {
  return role ? ADMIN_ROLES.has(role as AdminRole) : false
}

function normalizeAdminAccount(input: Partial<AdminAccount>): AdminAccount {
  return {
    id: input.id ?? '',
    name: input.name ?? '',
    email: input.email ?? '',
    phone: input.phone ?? '',
    role: input.role,
    tenantId: input.tenantId ?? '',
    mfaEnabled: input.mfaEnabled ?? false,
    createdAt: input.createdAt ?? '',
  }
}

export function roleLabel(role?: string): string {
  switch (role) {
    case 'owner':
      return 'Proprietario'
    case 'admin':
      return 'Administrador'
    case 'receptionist':
      return 'Recepcionista'
    case 'professional':
      return 'Profissional'
    default:
      return role ?? 'Sem perfil'
  }
}

export type GoogleAuthStatus = {
  enabled: boolean
  clientId?: string
}

export async function getGoogleAuthStatus(): Promise<GoogleAuthStatus> {
  try {
    const response = await fetch(buildApiUrl('/api/v1/auth/google/status'))
    if (!response.ok) {
      return { enabled: false }
    }
    const data = await response.json()
    return {
      enabled: Boolean(data?.enabled),
      clientId: data?.clientId,
    }
  } catch {
    return { enabled: false }
  }
}

export function initiateGoogleLogin(returnUrl?: string) {
  if (typeof window === 'undefined') return
  const currentUrl = returnUrl || window.location.href.split('#')[0]
  const targetUrl = new URL(buildApiUrl('/api/v1/auth/google/login'), window.location.origin)
  targetUrl.searchParams.set('audience', 'admin')
  targetUrl.searchParams.set('return_url', currentUrl)
  window.location.href = targetUrl.toString()
}

export function handleGoogleAuthRedirect(): {
  user: AdminAccount | null
  error: string | null
} {
  if (typeof window === 'undefined') {
    return { user: null, error: null }
  }

  // Check URL fragment: #google_auth=success&csrfToken=...&user=...
  const hash = window.location.hash
  if (hash && hash.includes('google_auth=success')) {
    const params = new URLSearchParams(hash.replace(/^#/, ''))
    const csrfToken = params.get('csrfToken') || ''
    const userRaw = params.get('user')
    if (csrfToken && userRaw) {
      try {
        const parsed = JSON.parse(decodeURIComponent(userRaw)) as Partial<AdminAccount>
        const user = normalizeAdminAccount(parsed)
        persistAdminSession(user, csrfToken)
        window.history.replaceState(null, '', window.location.pathname + window.location.search)
        return { user, error: null }
      } catch {
        // ignore parse error
      }
    }
  }

  // Check query parameter error: ?google_error=...
  const searchParams = new URLSearchParams(window.location.search)
  const googleError = searchParams.get('google_error')
  if (googleError) {
    searchParams.delete('google_error')
    const newSearch = searchParams.toString()
    const newUrl = window.location.pathname + (newSearch ? `?${newSearch}` : '') + window.location.hash
    window.history.replaceState(null, '', newUrl)

    if (googleError === 'not_configured') {
      return {
        user: null,
        error:
          'A integração com o Google ainda não foi configurada. As chaves GOOGLE_CLIENT_ID e GOOGLE_CLIENT_SECRET devem ser adicionadas nas variáveis de ambiente.',
      }
    }
    if (googleError === 'cancelled') {
      return { user: null, error: 'O login com o Google foi cancelado.' }
    }
    return { user: null, error: 'Não foi possível concluir o login com o Google. Tente novamente.' }
  }

  return { user: null, error: null }
}
