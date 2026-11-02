const configuredAdmin = import.meta.env.VITE_ADMIN_URL?.trim()
const configuredSignup = import.meta.env.VITE_SIGNUP_URL?.trim()
const configuredBooking = import.meta.env.VITE_BOOKING_URL?.trim()
const configuredWidget = import.meta.env.VITE_WIDGET_SCRIPT_URL?.trim()

export const APP_LINKS = {
  admin: configuredAdmin || './admin/',
  signup: configuredSignup || './admin/#/signup',
  booking: configuredBooking || './booking/',
}

export const WIDGET_SCRIPT_URL = configuredWidget || './widget.js'

export const SECTION_LINKS = {
  features: '#funcionalidades',
  howItWorks: '#como-funciona',
  pricing: '#precos',
}
