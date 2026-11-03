export const SHORT_DAY_LABELS = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

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

export function formatDateLong(dateKey: string): string {
  if (!dateKey) return ''
  const date = fromDateKey(dateKey)
  if (Number.isNaN(date.getTime())) return ''
  const dayLabel = SHORT_DAY_LABELS[date.getDay()] ?? ''
  const month = date.toLocaleDateString('pt-BR', { month: 'long' })
  return `${dayLabel}, ${date.getDate()} de ${month}`
}

export function formatTime(rfc3339: string): string {
  if (!rfc3339) return ''
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

export function formatDateTime(rfc3339: string): string {
  if (!rfc3339) return ''
  const date = new Date(rfc3339)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const dayLabel = SHORT_DAY_LABELS[date.getDay()] ?? ''
  const month = date.toLocaleDateString('pt-BR', { month: 'long' })
  const time = formatTime(rfc3339)
  return `${dayLabel}, ${date.getDate()} de ${month}${time ? ` às ${time}` : ''}`
}

export function formatMoney(valueInCents: number): string {
  const safe = Number.isFinite(valueInCents) ? valueInCents : 0
  return new Intl.NumberFormat('pt-BR', { currency: 'BRL', style: 'currency' }).format(safe / 100)
}

export function nextDateKeys(limit: number, offsetStartDays = 0): string[] {
  const keys: string[] = []
  const today = new Date()
  for (let index = offsetStartDays; index < offsetStartDays + limit; index += 1) {
    const date = new Date(today)
    date.setDate(today.getDate() + index)
    keys.push(toDateKey(date))
  }
  return keys
}
