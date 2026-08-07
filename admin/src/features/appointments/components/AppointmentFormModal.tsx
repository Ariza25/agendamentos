import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createAppointment, createCustomer } from '../../../lib/adminApi'
import { APPOINTMENT_STATUS_LABELS } from '../../../lib/appointmentStatus'
import { addMinutes, localDateTimeValue, toRFC3339 } from '../../../lib/format'
import { Button } from '../../../shared/components/Button'
import { Checkbox, Field, Form, Input, Select, Textarea } from '../../../shared/components/FormControls'
import { Modal } from '../../../shared/components/Modal'
import { useToast } from '../../../shared/components/toast/ToastProvider'
import type { Appointment, Customer, Professional, Resource, Service, AppointmentStatus } from '../../../types'

type AppointmentFormModalProps = {
  open: boolean
  onClose: () => void
  onSaved: () => void
  professionals: Professional[]
  services: Service[]
  resources: Resource[]
  customers: Customer[]
  defaultProfessionalId?: string
  defaultDate?: string
}

const STATUS_OPTIONS: AppointmentStatus[] = ['scheduled', 'confirmed', 'in_progress', 'completed', 'cancelled', 'no_show']

export function AppointmentFormModal({
  open,
  onClose,
  onSaved,
  professionals,
  services,
  resources,
  customers,
  defaultProfessionalId,
  defaultDate,
}: AppointmentFormModalProps) {
  const { showToast } = useToast()
  const [professionalId, setProfessionalId] = useState(defaultProfessionalId ?? '')
  const [serviceId, setServiceId] = useState('')
  const [resourceId, setResourceId] = useState('')
  const [customerId, setCustomerId] = useState('')
  const [isCreatingNewCustomer, setIsCreatingNewCustomer] = useState(false)
  const [newCustomerName, setNewCustomerName] = useState('')
  const [newCustomerPhone, setNewCustomerPhone] = useState('')
  const [startsAt, setStartsAt] = useState('')
  const [status, setStatus] = useState<AppointmentStatus>('scheduled')
  const [notes, setNotes] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (!open) {
      return
    }
    setProfessionalId(defaultProfessionalId ?? '')
    setServiceId('')
    setResourceId('')
    setCustomerId('')
    setIsCreatingNewCustomer(false)
    setNewCustomerName('')
    setNewCustomerPhone('')
    setStartsAt(defaultDate ?? localDateTimeValue(new Date().toISOString()))
    setStatus('scheduled')
    setNotes('')
    setError('')
  }, [open, defaultProfessionalId, defaultDate])

  const selectedService = services.find((service) => service.id === serviceId)

  function computeEndsAt(startValue: string): string {
    const startsAtValue = toRFC3339(startValue)
    const duration = selectedService?.durationMinutes ?? 30
    return addMinutes(startsAtValue, duration)
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!professionalId) {
      setError('Selecione o profissional.')
      return
    }
    if (!startsAt) {
      setError('Informe data e horário do agendamento.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      let resolvedCustomerId = customerId
      if (isCreatingNewCustomer) {
        if (!newCustomerName.trim()) {
          setError('Informe o nome do novo cliente.')
          setIsSubmitting(false)
          return
        }
        const customer = await createCustomer({
          name: newCustomerName.trim(),
          phone: newCustomerPhone.trim(),
        })
        resolvedCustomerId = customer.id
      }

      if (!resolvedCustomerId) {
        setError('Selecione ou cadastre um cliente.')
        setIsSubmitting(false)
        return
      }

      await createAppointment({
        professionalId,
        serviceId: serviceId || '',
        resourceId: resourceId || '',
        customerId: resolvedCustomerId,
        startsAt: toRFC3339(startsAt),
        endsAt: computeEndsAt(startsAt),
        status,
        notes: notes.trim(),
      })

      showToast({ message: 'Agendamento criado.', tone: 'success' })
      onSaved()
      onClose()
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
          <Button loading={isSubmitting} onClick={() => undefined} type="submit" form="appointment-form">
            Salvar agendamento
          </Button>
        </>
      }
      onClose={onClose}
      open={open}
      size="lg"
      subtitle="Registre um horário para o cliente."
      title="Novo agendamento"
    >
      <Form id="appointment-form" onSubmit={handleSubmit}>
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Cliente">
            {isCreatingNewCustomer ? (
              <div className="mt-2 grid gap-2">
                <Input
                  onChange={(event) => setNewCustomerName(event.target.value)}
                  placeholder="Nome do novo cliente"
                  value={newCustomerName}
                />
                <Input
                  onChange={(event) => setNewCustomerPhone(event.target.value)}
                  placeholder="Telefone (ex.: 11 98888 7777)"
                  value={newCustomerPhone}
                />
              </div>
            ) : (
              <Select
                className="mt-2"
                onChange={(event) => setCustomerId(event.target.value)}
                value={customerId}
              >
                <option value="">Selecione um cliente...</option>
                {customers.map((customer) => (
                  <option key={customer.id} value={customer.id}>
                    {customer.name}
                    {customer.phone ? ` · ${customer.phone}` : ''}
                  </option>
                ))}
              </Select>
            )}
          </Field>

          <Field className="mt-2" label=" ">
            <Checkbox checked={isCreatingNewCustomer} onChange={(event) => setIsCreatingNewCustomer(event.target.checked)} />
            <span className="ml-2 text-sm font-bold text-[var(--color-ink)]">Criar novo cliente</span>
          </Field>

          <Field label="Profissional">
            <Select className="mt-2" onChange={(event) => setProfessionalId(event.target.value)} value={professionalId}>
              <option value="">Selecione um profissional...</option>
              {professionals.filter((professional) => professional.active).map((professional) => (
                <option key={professional.id} value={professional.id}>
                  {professional.name}
                </option>
              ))}
            </Select>
          </Field>

          <Field label="Serviço (opcional)">
            <Select className="mt-2" onChange={(event) => setServiceId(event.target.value)} value={serviceId}>
              <option value="">Sem serviço (30 min)</option>
              {services.filter((service) => service.active).map((service) => (
                <option key={service.id} value={service.id}>
                  {service.name} · {service.durationMinutes} min
                </option>
              ))}
            </Select>
          </Field>

          <Field label="Data e horário">
            <Input
              className="mt-2"
              onChange={(event) => setStartsAt(event.target.value)}
              type="datetime-local"
              value={startsAt}
            />
          </Field>

          <Field hint="Deixe vazio para alocar automaticamente." label="Recurso (opcional)">
            <Select className="mt-2" onChange={(event) => setResourceId(event.target.value)} value={resourceId}>
              <option value="">Alocar automaticamente</option>
              {resources.filter((resource) => resource.active).map((resource) => (
                <option key={resource.id} value={resource.id}>
                  {resource.name}
                </option>
              ))}
            </Select>
          </Field>

          <Field label="Status">
            <Select className="mt-2" onChange={(event) => setStatus(event.target.value as AppointmentStatus)} value={status}>
              {STATUS_OPTIONS.map((option) => (
                <option key={option} value={option}>
                  {APPOINTMENT_STATUS_LABELS[option]}
                </option>
              ))}
            </Select>
          </Field>
        </div>

        <Field label="Observações">
          <Textarea className="mt-2" onChange={(event) => setNotes(event.target.value)} placeholder="Anotações sobre o atendimento..." value={notes} />
        </Field>

        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
