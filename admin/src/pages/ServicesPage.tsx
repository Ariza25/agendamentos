import { Plus, Scissors } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createService, deleteService, fetchServices, setServiceActive, updateService } from '../lib/adminApi'
import { formatMoney } from '../lib/format'
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
import type { Service } from '../types'

const COLOR_OPTIONS = ['#b42318', '#0f766e', '#b45309', '#4338ca', '#be185d', '#15803d', '#c2410c', '#475569']

export function ServicesPage() {
  const { showToast } = useToast()
  const [services, setServices] = useState<Service[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingService, setEditingService] = useState<Service | null>(null)
  const [serviceToDelete, setServiceToDelete] = useState<Service | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadServices() {
    setLoading(true)
    fetchServices()
      .then(setServices)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os serviços.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadServices, [])

  async function handleToggleActive(service: Service) {
    try {
      await setServiceActive(service.id, !service.active)
      showToast({ message: service.active ? 'Serviço desativado.' : 'Serviço ativado.', tone: 'success' })
      loadServices()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!serviceToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteService(serviceToDelete.id)
      showToast({ message: 'Serviço removido.', tone: 'success' })
      setServiceToDelete(null)
      loadServices()
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
              setEditingService(null)
              setIsFormOpen(true)
            }}
          >
            Novo serviço
          </Button>
        }
        subtitle="Serviços disponíveis para agendamento."
        title="Serviços"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : services.length === 0 ? (
        <Section>
          <EmptyState
            icon={<Scissors className="h-10 w-10" />}
            description="Nenhum serviço cadastrado. Cadastre um serviço para permitir agendamentos."
            title="Sem serviços cadastrados"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {services.map((service) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={service.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="h-3 w-3 shrink-0 rounded-full" style={{ backgroundColor: service.color || '#b42318' }} />
                      <p className="truncate font-bold text-[var(--color-ink)]">{service.name}</p>
                    </div>
                    <p className="mt-1 truncate text-sm text-[var(--color-ink-soft)]">
                      {service.durationMinutes} min · {formatMoney(service.priceCents)}
                    </p>
                  </div>
                  <Badge tone={service.active ? 'success' : 'danger'}>{service.active ? 'Ativo' : 'Inativo'}</Badge>
                </div>
                {service.description ? (
                  <p className="line-clamp-2 text-sm text-[var(--color-ink-soft)]">{service.description}</p>
                ) : null}
                <div className="mt-auto flex items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingService(service)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => void handleToggleActive(service)} size="sm" variant="secondary">
                    {service.active ? 'Desativar' : 'Ativar'}
                  </Button>
                  <Button onClick={() => setServiceToDelete(service)} size="sm" variant="danger">
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isFormOpen ? (
        <ServiceFormModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadServices()
          }}
          service={editingService}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O serviço "${serviceToDelete?.name ?? ''}" será removido do catálogo.`}
        onCancel={() => setServiceToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(serviceToDelete)}
        title="Remover serviço?"
      />
    </PageContainer>
  )
}

function ServiceFormModal({
  service,
  onClose,
  onSaved,
}: {
  service: Service | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(service?.name ?? '')
  const [description, setDescription] = useState(service?.description ?? '')
  const [durationMinutes, setDurationMinutes] = useState(service?.durationMinutes ?? 30)
  const [priceCents, setPriceCents] = useState(service?.priceCents ?? 0)
  const [color, setColor] = useState(service?.color ?? COLOR_OPTIONS[0])
  const [active, setActive] = useState(service?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome do serviço.')
      return
    }
    if (!durationMinutes || durationMinutes <= 0) {
      setError('Informe a duração em minutos.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = {
        name: name.trim(),
        description: description.trim(),
        durationMinutes,
        priceCents: Math.max(0, Math.round(priceCents)),
        color,
        active,
      }
      if (service) {
        await updateService(service.id, input)
        showToast({ message: 'Serviço atualizado.', tone: 'success' })
      } else {
        await createService(input)
        showToast({ message: 'Serviço criado.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="service-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={service ? 'Atualize os dados do serviço.' : 'Cadastre um novo serviço.'}
      title={service ? 'Editar serviço' : 'Novo serviço'}
    >
      <Form id="service-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <Field label="Descrição">
          <Textarea className="mt-2" onChange={(event) => setDescription(event.target.value)} value={description} />
        </Field>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Duração (minutos)">
            <Input
              className="mt-2"
              min={1}
              onChange={(event) => setDurationMinutes(Number(event.target.value))}
              type="number"
              value={durationMinutes}
            />
          </Field>
          <Field label="Preço (em centavos)">
            <Input
              className="mt-2"
              min={0}
              onChange={(event) => setPriceCents(Number(event.target.value))}
              type="number"
              value={priceCents}
            />
          </Field>
        </div>
        <Field hint={`Exibido como ${formatMoney(priceCents)}.`} label="Cor">
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
          <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Serviço ativo (visível para agendamento online)</span>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
