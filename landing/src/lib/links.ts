const isProd = import.meta.env.PROD

export const APP_LINKS = {
  admin: isProd ? './admin/' : 'http://localhost:5173',
  signup: isProd ? './admin/#/signup' : 'http://localhost:5173/#/signup',
  booking: isProd ? './booking/' : 'http://localhost:5175',
}

export const WIDGET_SCRIPT_URL = isProd
  ? 'https://erpagendamentos-backend-88726602471.southamerica-east1.run.app/widget.js'
  : 'http://localhost:8080/widget.js'

export const SECTION_LINKS = {
  features: '#funcionalidades',
  howItWorks: '#como-funciona',
  pricing: '#precos',
}
