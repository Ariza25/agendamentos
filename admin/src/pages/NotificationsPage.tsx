import { Bell, Send } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { fetchNotifications, sendTestNotification } from '../lib/adminApi'
import { formatDateTime } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select, Textarea } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { NotificationLog } from '../types'

const NOTIFICATION_TYPE_LABELS: Record<string, string> = {
  confirmation: 'Confirmação',
  reminder: 'Lembrete',
  post_service: 'Pós-atendimento',
  no_show_fee: 'Taxa de não comparecimento',
  waiting_list: 'Lista de espera',
  cancellation: 'Cancelamento',
  reschedule: 'Reagendamento',
  test: 'Teste',
}

const CHANNEL_OPTIONS = [
  { value: 'sms', label: 'SMS' },
  { value: 'email', label: 'Email' },
  { value: 'whatsapp', label: 'WhatsApp' },
]

export function NotificationsPage() {
  const { showToast } = useToast()
  const [notifications, setNotifications] = useState<NotificationLog[]>([])
  const [loading, setLoading] = useState(true)
  const [isTestOpen, setIsTestOpen] = useState(false)

  function loadNotifications() {
    setLoading(true)
    fetchNotifications()
      .then(setNotifications)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar as notificações.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadNotifications, [])

  return (
    <PageContainer>
      <Section
        actions={
          <Button icon={<Send className="h-4 w-4" />} onClick={() => setIsTestOpen(true)} variant="secondary">
            Enviar teste
          </Button>
        }
        subtitle="Histórico de lembretes, confirmações e avisos enviados aos clientes."
        title="Notificações"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : notifications.length === 0 ? (
        <Section>
          <EmptyState
            description="Os lembretes e confirmações enviados aos clientes aparecerão aqui."
            icon={<Bell className="h-10 w-10" />}
            title="Nenhuma notificação enviada"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3">
            {notifications.map((notification) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 md:flex-row md:items-start"
                key={notification.id}
              >
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge tone="accent">{NOTIFICATION_TYPE_LABELS[notification.type] ?? notification.type}</Badge>
                    <Badge tone="neutral">{CHANNEL_OPTIONS.find((option) => option.value === notification.channel)?.label ?? notification.channel}</Badge>
                    <Badge tone={notification.status === 'sent' ? 'success' : notification.status === 'failed' ? 'danger' : 'warning'}>
                      {notification.status === 'sent' ? 'Enviada' : notification.status === 'failed' ? 'Falhou' : notification.status}
                    </Badge>
                    <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
                      {notification.createdAt ? formatDateTime(notification.createdAt) : ''}
                    </span>
                  </div>
                  <p className="mt-1.5 font-bold text-[var(--color-ink)]">{notification.subject}</p>
                  <p className="mt-0.5 truncate text-sm text-[var(--color-ink-soft)]">{notification.body}</p>
                  <p className="mt-1 text-xs font-semibold text-[var(--color-ink-soft)]">Para: {notification.recipient}</p>
                  {notification.error ? (
                    <p className="mt-1 truncate text-xs font-semibold text-red-700">Erro: {notification.error}</p>
                  ) : null}
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      {isTestOpen ? (
        <SendTestModal
          onClose={() => setIsTestOpen(false)}
          onSaved={() => {
            setIsTestOpen(false)
            loadNotifications()
          }}
        />
      ) : null}
    </PageContainer>
  )
}

function SendTestModal({ onClose, onSaved }: { onClose: () => void; onSaved: () => void }) {
  const { showToast } = useToast()
  const [channel, setChannel] = useState('email')
  const [recipient, setRecipient] = useState('')
  const [subject, setSubject] = useState('')
  const [body, setBody] = useState('')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!recipient.trim()) {
      setError('Informe o destinatário (telefone ou email).')
      return
    }
    if (!subject.trim() || !body.trim()) {
      setError('Informe o assunto e a mensagem.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      await sendTestNotification({
        channel,
        recipient: recipient.trim(),
        subject: subject.trim(),
        body: body.trim(),
      })
      showToast({ message: 'Notificação de teste enviada.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="send-test-form">
            Enviar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle="Envie uma notificação de teste para validar o canal configurado."
      title="Enviar notificação de teste"
    >
      <Form id="send-test-form" onSubmit={handleSubmit}>
        <Field label="Canal">
          <Select className="mt-2" onChange={(event) => setChannel(event.target.value)} value={channel}>
            {CHANNEL_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </Select>
        </Field>
        <Field hint="Telefone com DDI para SMS/WhatsApp ou email válido." label="Destinatário">
          <Input className="mt-2" onChange={(event) => setRecipient(event.target.value)} value={recipient} />
        </Field>
        <Field label="Assunto">
          <Input className="mt-2" onChange={(event) => setSubject(event.target.value)} value={subject} />
        </Field>
        <Field label="Mensagem">
          <Textarea className="mt-2" onChange={(event) => setBody(event.target.value)} value={body} />
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
