export const DAY_LABELS = ['Domingo', 'Segunda-feira', 'Terça-feira', 'Quarta-feira', 'Quinta-feira', 'Sexta-feira', 'Sábado']

export const SHORT_DAY_LABELS = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

export function todayKey(offsetDays = 0): string {
  const date = new Date()
  date.setDate(date.getDate() + offsetDays)
  return toDateKey(date)
}

export function toDateKey(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function fromDateKey(dateKey: string): Date {
  const [year, month, day] = dateKey.split('-').map(Number)
  return new Date(year, month - 1, day)
}

export function formatDate(dateKey: string): string {
  const date = fromDateKey(dateKey)
  return `${DAY_LABELS[date.getDay()]}, ${String(date.getDate()).padStart(2, '0')}/${String(date.getMonth() + 1).padStart(2, '0')}/${date.getFullYear()}`
}

export function formatTime(rfc3339: string): string {
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

export function formatDateTime(rfc3339: string): string {
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const key = toDateKey(date)
  return `${formatDate(key)} às ${formatTime(rfc3339)}`
}

export function formatMoney(valueInCents: number): string {
  return new Intl.NumberFormat('pt-BR', { currency: 'BRL', style: 'currency' }).format((valueInCents ?? 0) / 100)
}

export function localDateTimeValue(rfc3339: string): string {
  const date = new Date(rfc3339)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

export function toRFC3339(localValue: string): string {
  if (!localValue) {
    return ''
  }
  const date = new Date(localValue)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return date.toISOString()
}

export function addMinutes(rfc3339: string, minutes: number): string {
  const date = new Date(rfc3339)
  date.setMinutes(date.getMinutes() + minutes)
  return date.toISOString()
}
