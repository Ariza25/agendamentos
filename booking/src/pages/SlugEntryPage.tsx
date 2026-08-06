import { CalendarDays, ChevronRight } from 'lucide-react'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'

import { Button } from '../shared/components/Button'
import { Field, Form, Input } from '../shared/components/FormControls'

export function SlugEntryPage() {
  const navigate = useNavigate()
  const [slug, setSlug] = useState('')
  const [error, setError] = useState('')

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const trimmed = slug.trim()
    if (!trimmed) {
      setError('Informe o nome do estabelecimento para continuar.')
      return
    }
    navigate(`/booking/${encodeURIComponent(trimmed)}`)
  }

  return (
    <main className="flex min-h-dvh items-center justify-center px-4 py-10">
      <section className="w-full max-w-md rounded-3xl border border-[var(--color-line)] bg-white p-8 shadow-sm">
        <div className="flex items-center gap-3">
          <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
            <CalendarDays className="h-6 w-6" />
          </span>
          <div>
            <h1 className="font-headline text-2xl font-extrabold text-[var(--color-ink)]">Agende seu horário</h1>
            <p className="text-sm font-semibold text-[var(--color-ink-soft)]">Escolha um serviço e reserve em instantes.</p>
          </div>
        </div>

        <Form onSubmit={handleSubmit}>
          <Field error={error} hint="Ex.: meubarbearia" label="Nome do estabelecimento">
            <Input
              autoFocus
              className="mt-2"
              onChange={(event) => {
                setSlug(event.target.value)
                setError('')
              }}
              placeholder="Digite o nome do estabelecimento"
              value={slug}
            />
          </Field>
          <Button icon={<ChevronRight className="h-4 w-4" />} size="lg" type="submit">
            Continuar
          </Button>
        </Form>
      </section>
    </main>
  )
}
