export function Spinner({ label = 'Carregando...' }: { label?: string }) {
  return (
    <div className="flex items-center justify-center gap-3 px-6 py-12">
      <span className="h-5 w-5 animate-spin rounded-full border-2 border-[var(--color-line)] border-t-[var(--color-primary)]" />
      <span className="text-sm font-bold text-[var(--color-ink-soft)]">{label}</span>
    </div>
  )
}
