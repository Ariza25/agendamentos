import { Plus, UserCog } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createAdminUser, deleteAdminUser, fetchAdminUsers, updateAdminUserRole } from '../lib/adminApi'
import { isAdminRole, roleLabel, type AdminRole } from '../lib/adminAuth'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { AdminUser } from '../types'

const ROLE_OPTIONS: AdminRole[] = ['owner', 'admin', 'receptionist', 'professional']

export function TeamPage() {
  const { showToast } = useToast()
  const [users, setUsers] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [userToDelete, setUserToDelete] = useState<AdminUser | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)
  const [lastInvite, setLastInvite] = useState<string>('')

  function loadUsers() {
    setLoading(true)
    fetchAdminUsers()
      .then(setUsers)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os usuários.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadUsers, [])

  async function handleRoleChange(user: AdminUser, role: string) {
    if (!isAdminRole(role)) {
      return
    }
    try {
      await updateAdminUserRole(user.id, role)
      showToast({ message: 'Perfil atualizado.', tone: 'success' })
      loadUsers()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!userToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteAdminUser(userToDelete.id)
      showToast({ message: 'Usuário removido.', tone: 'success' })
      setUserToDelete(null)
      loadUsers()
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
              setLastInvite('')
              setIsFormOpen(true)
            }}
          >
            Novo usuário
          </Button>
        }
        subtitle="Contas com acesso ao painel administrativo."
        title="Equipe"
      >
        <div className="mt-2" />
      </Section>

      {lastInvite ? (
        <Section title="Convite gerado">
          <p className="text-sm text-[var(--color-ink-soft)]">
            Compartilhe o link abaixo com o novo usuário para que ele defina a própria senha:
          </p>
          <p className="mt-3 break-all rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3 font-mono text-sm font-bold text-[var(--color-ink)]">
            {lastInvite}
          </p>
          <Button
            className="mt-4"
            onClick={() => {
              void navigator.clipboard.writeText(lastInvite)
              showToast({ message: 'Link copiado.', tone: 'success' })
            }}
            variant="secondary"
          >
            Copiar link
          </Button>
        </Section>
      ) : null}

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : users.length === 0 ? (
        <Section>
          <EmptyState
            icon={<UserCog className="h-10 w-10" />}
            description="Nenhum usuário com acesso administrativo."
            title="Sem usuários cadastrados"
          />
        </Section>
      ) : (
        <Section>
          <div className="mt-2 overflow-x-auto">
            <table className="w-full min-w-0 border-collapse text-left">
              <thead>
                <tr className="text-xs font-black uppercase tracking-wide text-[var(--color-ink-soft)]">
                  <th className="border-b border-[var(--color-line)] py-3 pr-4">Usuário</th>
                  <th className="border-b border-[var(--color-line)] px-4 py-3">Perfil</th>
                  <th className="border-b border-[var(--color-line)] py-3 text-right">Ações</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr className="text-sm" key={user.id}>
                    <td className="border-b border-[var(--color-line)] py-3 pr-4">
                      <p className="font-bold text-[var(--color-ink)]">{user.name}</p>
                      <p className="text-xs text-[var(--color-ink-soft)]">{user.email}</p>
                    </td>
                    <td className="border-b border-[var(--color-line)] px-4 py-3">
                      <Badge tone={user.mfaEnabled ? 'success' : 'neutral'}>
                        {roleLabel(user.role)}
                        {user.mfaEnabled ? ' · MFA ativo' : ''}
                      </Badge>
                    </td>
                    <td className="border-b border-[var(--color-line)] py-3 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Select
                          className="h-9 w-auto min-w-36 text-xs"
                          onChange={(event) => void handleRoleChange(user, event.target.value)}
                          value={user.role}
                        >
                          {ROLE_OPTIONS.map((role) => (
                            <option key={role} value={role}>
                              {roleLabel(role)}
                            </option>
                          ))}
                        </Select>
                        <Button onClick={() => setUserToDelete(user)} size="sm" variant="danger">
                          Remover
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Section>
      )}

      {isFormOpen ? (
        <UserFormModal
          onClose={() => setIsFormOpen(false)}
          onInvite={(inviteUrl) => {
            setIsFormOpen(false)
            setLastInvite(inviteUrl)
            loadUsers()
          }}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O usuário "${userToDelete?.name ?? ''}" perderá o acesso ao painel.`}
        onCancel={() => setUserToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(userToDelete)}
        title="Remover usuário?"
      />
    </PageContainer>
  )
}

function UserFormModal({ onClose, onInvite }: { onClose: () => void; onInvite: (inviteUrl: string) => void }) {
  const { showToast } = useToast()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [phone, setPhone] = useState('')
  const [role, setRole] = useState<AdminRole>('receptionist')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome do usuário.')
      return
    }
    if (!email.trim()) {
      setError('Informe o email do usuário.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const response = await createAdminUser({ name: name.trim(), email: email.trim(), phone: phone.trim(), role })
      showToast({ message: 'Usuário criado. Compartilhe o link de convite.', tone: 'success' })
      onInvite(response.inviteUrl || '')
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
          <Button loading={isSubmitting} type="submit" form="user-form">
            Criar usuário
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle="O usuário receberá um link para definir a própria senha."
      title="Novo usuário"
    >
      <Form id="user-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <Field label="Email">
          <Input className="mt-2" onChange={(event) => setEmail(event.target.value)} type="email" value={email} />
        </Field>
        <Field label="Telefone (opcional)">
          <Input className="mt-2" onChange={(event) => setPhone(event.target.value)} value={phone} />
        </Field>
        <Field hint="Recepcionistas gerenciam agenda e agendamentos; profissionais veem agenda e agendamentos; admin/owner administram tudo." label="Perfil de acesso">
          <Select className="mt-2" onChange={(event) => setRole(event.target.value as AdminRole)} value={role}>
            {ROLE_OPTIONS.map((option) => (
              <option key={option} value={option}>
                {roleLabel(option)}
              </option>
            ))}
          </Select>
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
