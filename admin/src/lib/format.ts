export const DAY_LABELS = ['Domingo', 'Segunda-feira', 'Terça-feira', 'Quarta-feira', 'Quinta-feira', 'Sexta-feira', 'Sábado']

export const SHORT_DAY_LABELS = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

export function todayKey(offsetDays = 0): string {
  const date = new Date()
  date.setDate(date.getDate() + offsetDays)
  return toDateKey(date)
}

export function toDateKey(date: Date): string {
  if (!date || Number.isNaN(date.getTime())) {
    return ''
  }
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function fromDateKey(dateKey: string): Date {
  if (!dateKey) {
    return new Date()
  }
  if (dateKey.includes('T')) {
    const parsed = new Date(dateKey)
    if (!Number.isNaN(parsed.getTime())) {
      return parsed
    }
  }
  const clean = dateKey.slice(0, 10)
  const [year, month, day] = clean.split('-').map(Number)
  if (!year || !month || !day || Number.isNaN(year) || Number.isNaN(month) || Number.isNaN(day)) {
    const fallback = new Date(dateKey)
    return Number.isNaN(fallback.getTime()) ? new Date() : fallback
  }
  return new Date(year, month - 1, day)
}

export function formatDate(dateKey: string): string {
  if (!dateKey) {
    return ''
  }
  const date = fromDateKey(dateKey)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const dayName = DAY_LABELS[date.getDay()] ?? ''
  const dd = String(date.getDate()).padStart(2, '0')
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const yyyy = date.getFullYear()
  return dayName ? `${dayName}, ${dd}/${mm}/${yyyy}` : `${dd}/${mm}/${yyyy}`
}

export function formatTime(rfc3339: string): string {
  if (!rfc3339) {
    return ''
  }
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

export function formatDateTime(rfc3339: string): string {
  if (!rfc3339) {
    return ''
  }
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const key = toDateKey(date)
  const time = formatTime(rfc3339)
  return time ? `${formatDate(key)} às ${time}` : formatDate(key)
}

export function formatMoney(valueInCents: number): string {
  const safe = Number.isFinite(valueInCents) ? valueInCents : 0
  return new Intl.NumberFormat('pt-BR', { currency: 'BRL', style: 'currency' }).format(safe / 100)
}

export function localDateTimeValue(rfc3339: string): string {
  if (!rfc3339) {
    return ''
  }
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
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
  if (!rfc3339) {
    return ''
  }
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  date.setMinutes(date.getMinutes() + (Number.isFinite(minutes) ? minutes : 0))
  return date.toISOString()
}
