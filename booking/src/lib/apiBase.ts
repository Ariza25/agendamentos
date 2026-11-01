const prodOrigin = 'https://erpagendamentos-backend-88726602471.southamerica-east1.run.app'
const devOrigin = 'http://localhost:8080'
const configuredOrigin = import.meta.env.VITE_BACKEND_ORIGIN?.trim() ?? ''

export const backendOrigin = configuredOrigin || (import.meta.env.PROD ? prodOrigin : devOrigin)

export function buildApiUrl(path: string) {
  if (/^https?:\/\//i.test(path)) {
    return path
  }

  if (!path.startsWith('/')) {
    return `${backendOrigin}/${path}`
  }

  return `${backendOrigin}${path}`
}
