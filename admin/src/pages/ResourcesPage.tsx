import { DoorOpen, Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createResource,
  deleteResource,
  fetchResources,
  setResourceActive,
  updateResource,
} from '../lib/adminApi'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Checkbox, Field, Form, Input, Textarea } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Resource } from '../types'

const TYPE_OPTIONS = ['room', 'equipment', 'vehicle', 'other']

const TYPE_LABELS: Record<string, string> = {
  room: 'Sala',
  equipment: 'Equipamento',
  vehicle: 'Veículo',
  other: 'Outro',
}

export function ResourcesPage() {
  const { showToast } = useToast()
  const [resources, setResources] = useState<Resource[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingResource, setEditingResource] = useState<Resource | null>(null)
  const [resourceToDelete, setResourceToDelete] = useState<Resource | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadResources() {
    setLoading(true)
    fetchResources()
      .then(setResources)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os recursos.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadResources, [])

  async function handleToggleActive(resource: Resource) {
    try {
      await setResourceActive(resource.id, !resource.active)
      showToast({ message: resource.active ? 'Recurso desativado.' : 'Recurso ativado.', tone: 'success' })
      loadResources()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!resourceToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteResource(resourceToDelete.id)
      showToast({ message: 'Recurso removido.', tone: 'success' })
      setResourceToDelete(null)
      loadResources()
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
              setEditingResource(null)
              setIsFormOpen(true)
            }}
          >
            Novo recurso
          </Button>
        }
        subtitle="Salas, equipamentos e outros recursos compartilhados entre profissionais."
        title="Recursos e salas"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : resources.length === 0 ? (
        <Section>
          <EmptyState
            icon={<DoorOpen className="h-10 w-10" />}
            description="Nenhum recurso cadastrado. Cadastre salas ou equipamentos para alocar automaticamente aos agendamentos."
            title="Sem recursos cadastrados"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {resources.map((resource) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={resource.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-bold text-[var(--color-ink)]">{resource.name}</p>
                    <p className="mt-1 truncate text-sm text-[var(--color-ink-soft)]">
                      {TYPE_LABELS[resource.type] ?? resource.type} · capacidade {resource.capacity}
                    </p>
                  </div>
                  <Badge tone={resource.active ? 'success' : 'danger'}>{resource.active ? 'Ativo' : 'Inativo'}</Badge>
                </div>
                {resource.notes ? (
                  <p className="line-clamp-2 text-sm text-[var(--color-ink-soft)]">{resource.notes}</p>
                ) : null}
                <div className="mt-auto flex items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingResource(resource)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => void handleToggleActive(resource)} size="sm" variant="secondary">
                    {resource.active ? 'Desativar' : 'Ativar'}
                  </Button>
                  <Button onClick={() => setResourceToDelete(resource)} size="sm" variant="danger">
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isFormOpen ? (
        <ResourceFormModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadResources()
          }}
          resource={editingResource}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O recurso "${resourceToDelete?.name ?? ''}" será removido. Agendamentos existentes continuam com o histórico.`}
        onCancel={() => setResourceToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(resourceToDelete)}
        title="Remover recurso?"
      />
    </PageContainer>
  )
}

function ResourceFormModal({
  resource,
  onClose,
  onSaved,
}: {
  resource: Resource | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(resource?.name ?? '')
  const [type, setType] = useState(resource?.type ?? 'room')
  const [capacity, setCapacity] = useState(resource?.capacity ?? 1)
  const [notes, setNotes] = useState(resource?.notes ?? '')
  const [active, setActive] = useState(resource?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome do recurso.')
      return
    }
    if (!capacity || capacity <= 0) {
      setError('Informe a capacidade.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = {
        name: name.trim(),
        type,
        capacity: Math.max(1, Math.round(capacity)),
        notes: notes.trim(),
        active,
      }
      if (resource) {
        await updateResource(resource.id, input)
        showToast({ message: 'Recurso atualizado.', tone: 'success' })
      } else {
        await createResource(input)
        showToast({ message: 'Recurso criado.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="resource-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={resource ? 'Atualize os dados do recurso.' : 'Cadastre um novo recurso.'}
      title={resource ? 'Editar recurso' : 'Novo recurso'}
    >
      <Form id="resource-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Tipo">
            <select
              className="mt-2 w-full rounded-xl border border-[var(--color-line)] bg-white px-3 py-2 text-sm font-bold text-[var(--color-ink)] focus:outline-none focus:ring-2 focus:ring-violet-500"
              onChange={(event) => setType(event.target.value)}
              value={type}
            >
              {TYPE_OPTIONS.map((option) => (
                <option key={option} value={option}>
                  {TYPE_LABELS[option] ?? option}
                </option>
              ))}
            </select>
          </Field>
          <Field hint="Quantos agendamentos simultâneos o recurso suporta." label="Capacidade">
            <Input
              className="mt-2"
              min={1}
              onChange={(event) => setCapacity(Number(event.target.value))}
              type="number"
              value={capacity}
            />
          </Field>
        </div>
        <Field label="Observações">
          <Textarea className="mt-2" onChange={(event) => setNotes(event.target.value)} value={notes} />
        </Field>
        <Field label=" ">
          <Checkbox checked={active} onChange={(event) => setActive(event.target.checked)} />
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Recurso ativo (disponível para alocação)</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
