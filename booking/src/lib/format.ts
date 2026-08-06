export const SHORT_DAY_LABELS = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

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

export function formatDateLong(dateKey: string): string {
  const date = fromDateKey(dateKey)
  const month = date.toLocaleDateString('pt-BR', { month: 'long' })
  return `${SHORT_DAY_LABELS[date.getDay()]}, ${date.getDate()} de ${month}`
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
  const month = date.toLocaleDateString('pt-BR', { month: 'long' })
  return `${SHORT_DAY_LABELS[date.getDay()]}, ${date.getDate()} de ${month} às ${formatTime(rfc3339)}`
}

export function formatMoney(valueInCents: number): string {
  return new Intl.NumberFormat('pt-BR', { currency: 'BRL', style: 'currency' }).format((valueInCents ?? 0) / 100)
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
