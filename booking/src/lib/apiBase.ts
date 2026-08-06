const fallbackOrigin = 'http://localhost:8080'
const configuredOrigin = import.meta.env.VITE_BACKEND_ORIGIN?.trim() ?? ''

export const backendOrigin = configuredOrigin || fallbackOrigin

export function buildApiUrl(path: string) {
  if (/^https?:\/\//i.test(path)) {
    return path
  }

  if (!path.startsWith('/')) {
    return `${backendOrigin}/${path}`
  }

  return `${backendOrigin}${path}`
}
