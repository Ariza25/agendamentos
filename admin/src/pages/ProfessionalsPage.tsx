import { Plus, UserCog } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createProfessional, deleteProfessional, fetchProfessionals, setProfessionalActive, updateProfessional } from '../lib/adminApi'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Checkbox, Field, Form, Input } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Professional } from '../types'

const COLOR_OPTIONS = ['#b42318', '#0f766e', '#b45309', '#4338ca', '#be185d', '#15803d', '#c2410c', '#475569']

export function ProfessionalsPage() {
  const { showToast } = useToast()
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingProfessional, setEditingProfessional] = useState<Professional | null>(null)
  const [professionalToDelete, setProfessionalToDelete] = useState<Professional | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadProfessionals() {
    setLoading(true)
    fetchProfessionals()
      .then(setProfessionals)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os profissionais.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadProfessionals, [])

  async function handleToggleActive(professional: Professional) {
    try {
      await setProfessionalActive(professional.id, !professional.active)
      showToast({ message: professional.active ? 'Profissional desativado.' : 'Profissional ativado.', tone: 'success' })
      loadProfessionals()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!professionalToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteProfessional(professionalToDelete.id)
      showToast({ message: 'Profissional removido.', tone: 'success' })
      setProfessionalToDelete(null)
      loadProfessionals()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <Button
            icon={<Plus className="h-4 w-4" />}
            onClick={() => {
              setEditingProfessional(null)
              setIsFormOpen(true)
            }}
          >
            Novo profissional
          </Button>
        }
        subtitle="Profissionais que realizam os atendimentos."
        title="Profissionais"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : professionals.length === 0 ? (
        <Section>
          <EmptyState
            icon={<UserCog className="h-10 w-10" />}
            description="Nenhum profissional cadastrado. Cadastre um profissional para receber agendamentos."
            title="Sem profissionais cadastrados"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {professionals.map((professional) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={professional.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full font-headline text-sm font-black text-white" style={{ backgroundColor: professional.color || '#b42318' }}>
                        {professional.name.trim().charAt(0).toUpperCase() || '?'}
                      </span>
                      <p className="truncate font-bold text-[var(--color-ink)]">{professional.name}</p>
                    </div>
                    <p className="mt-1.5 truncate text-sm text-[var(--color-ink-soft)]">
                      {professional.email || 'Sem email'}
                      {professional.phone ? ` · ${professional.phone}` : ''}
                    </p>
                  </div>
                  <Badge tone={professional.active ? 'success' : 'danger'}>{professional.active ? 'Ativo' : 'Inativo'}</Badge>
                </div>
                <div className="mt-auto flex flex-wrap items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingProfessional(professional)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => void handleToggleActive(professional)} size="sm" variant="secondary">
                    {professional.active ? 'Desativar' : 'Ativar'}
                  </Button>
                  <Button onClick={() => setProfessionalToDelete(professional)} size="sm" variant="danger">
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isFormOpen ? (
        <ProfessionalFormModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadProfessionals()
          }}
          professional={editingProfessional}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O profissional "${professionalToDelete?.name ?? ''}" será removido. Agendamentos futuros ficarão sem vínculo.`}
        onCancel={() => setProfessionalToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(professionalToDelete)}
        title="Remover profissional?"
      />
    </PageContainer>
  )
}

function ProfessionalFormModal({
  professional,
  onClose,
  onSaved,
}: {
  professional: Professional | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(professional?.name ?? '')
  const [email, setEmail] = useState(professional?.email ?? '')
  const [phone, setPhone] = useState(professional?.phone ?? '')
  const [color, setColor] = useState(professional?.color ?? COLOR_OPTIONS[0])
  const [active, setActive] = useState(professional?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome do profissional.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = { name: name.trim(), email: email.trim(), phone: phone.trim(), color, active }
      if (professional) {
        await updateProfessional(professional.id, input)
        showToast({ message: 'Profissional atualizado.', tone: 'success' })
      } else {
        await createProfessional(input)
        showToast({ message: 'Profissional criado.', tone: 'success' })
      }
      onSaved()
    } catch (submitError) {
      setError((submitError as Error).message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      footer={
        <>
          <Button disabled={isSubmitting} onClick={onClose} variant="secondary">
            Cancelar
          </Button>
          <Button loading={isSubmitting} type="submit" form="professional-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={professional ? 'Atualize os dados do profissional.' : 'Cadastre um novo profissional.'}
      title={professional ? 'Editar profissional' : 'Novo profissional'}
    >
      <Form id="professional-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Email">
            <Input className="mt-2" onChange={(event) => setEmail(event.target.value)} type="email" value={email} />
          </Field>
          <Field label="Telefone">
            <Input className="mt-2" onChange={(event) => setPhone(event.target.value)} value={phone} />
          </Field>
        </div>
        <Field label="Cor">
          <div className="mt-2 flex flex-wrap gap-2">
            {COLOR_OPTIONS.map((option) => (
              <button
                aria-label={`Cor ${option}`}
                className={`h-8 w-8 rounded-full border-2 transition ${color === option ? 'border-[var(--color-ink)] scale-110' : 'border-transparent'}`}
                key={option}
                onClick={() => setColor(option)}
                style={{ backgroundColor: option }}
                type="button"
              />
            ))}
          </div>
        </Field>
        <Field label=" ">
          <Checkbox checked={active} onChange={(event) => setActive(event.target.checked)} />
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Profissional ativo (disponível para agendamento online)</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
