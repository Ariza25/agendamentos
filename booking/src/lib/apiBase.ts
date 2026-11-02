const configuredOrigin = import.meta.env.VITE_BACKEND_ORIGIN?.trim() ?? ''
const configuredPort = import.meta.env.VITE_BACKEND_PORT?.trim() ?? ''
const hostFallback = configuredPort ? `${window.location.protocol}//${window.location.hostname}:${configuredPort}` : ''

export const backendOrigin = configuredOrigin || hostFallback

export function buildApiUrl(path: string) {
  if (/^https?:\/\//i.test(path)) {
    return path
  }

  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  if (!backendOrigin) {
    return normalizedPath
  }

  return `${backendOrigin.replace(/\/+$/, '')}${normalizedPath}`
}
