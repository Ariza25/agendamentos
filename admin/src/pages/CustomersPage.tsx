import { Download, Plus, Search, Users } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { createCustomer, deleteCustomer, downloadAdminCsv, fetchCustomers, updateCustomer } from '../lib/adminApi'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Textarea } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Customer } from '../types'

export function CustomersPage() {
  const { showToast } = useToast()
  const [search, setSearch] = useState('')
  const [customers, setCustomers] = useState<Customer[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingCustomer, setEditingCustomer] = useState<Customer | null>(null)
  const [customerToDelete, setCustomerToDelete] = useState<Customer | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)
  const [exporting, setExporting] = useState(false)

  function loadCustomers(query = search) {
    setLoading(true)
    fetchCustomers(query)
      .then(setCustomers)
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os clientes.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadCustomers('')
  }, [])

  function handleSearch() {
    loadCustomers(search)
  }

  async function handleDelete() {
    if (!customerToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deleteCustomer(customerToDelete.id)
      showToast({ message: 'Cliente removido.', tone: 'success' })
      setCustomerToDelete(null)
      loadCustomers()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setIsDeleting(false)
    }
  }

  async function handleExport() {
    setExporting(true)
    try {
      await downloadAdminCsv('/api/v1/admin/export/customers.csv', 'clientes.csv')
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    } finally {
      setExporting(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <>
            <Button disabled={exporting} icon={<Download className="h-4 w-4" />} loading={exporting} onClick={() => void handleExport()}>
              Exportar CSV
            </Button>
            <Button
              icon={<Plus className="h-4 w-4" />}
              onClick={() => {
                setEditingCustomer(null)
                setIsFormOpen(true)
              }}
            >
              Novo cliente
            </Button>
          </>
        }
        subtitle="Cadastre e gerencie os clientes do estabelecimento."
        title="Clientes"
      >
        <div className="mt-4 flex max-w-lg items-end gap-2">
          <label className="block flex-1">
            <span className="text-xs font-bold uppercase tracking-wide text-[var(--color-ink-soft)]">Buscar</span>
            <Input
              className="mt-1"
              onChange={(event) => setSearch(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === 'Enter') {
                  handleSearch()
                }
              }}
              placeholder="Nome, telefone ou email"
              value={search}
            />
          </label>
          <Button onClick={handleSearch} variant="secondary">
            <Search className="h-4 w-4" />
            Buscar
          </Button>
        </div>
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : customers.length === 0 ? (
        <Section>
          <EmptyState
            icon={<Users className="h-10 w-10" />}
            description="Nenhum cliente encontrado. Cadastre o primeiro cliente para começar."
            title="Sem clientes cadastrados"
          />
        </Section>
      ) : (
        <Section>
          <div className="mt-2 overflow-x-auto">
            <table className="w-full min-w-0 border-collapse text-left">
              <thead>
                <tr className="text-xs font-black uppercase tracking-wide text-[var(--color-ink-soft)]">
                  <th className="border-b border-[var(--color-line)] py-3 pr-4">Nome</th>
                  <th className="border-b border-[var(--color-line)] px-4 py-3">Telefone</th>
                  <th className="border-b border-[var(--color-line)] px-4 py-3">Email</th>
                  <th className="border-b border-[var(--color-line)] py-3 text-right">Ações</th>
                </tr>
              </thead>
              <tbody>
                {customers.map((customer) => (
                  <tr className="text-sm" key={customer.id}>
                    <td className="border-b border-[var(--color-line)] py-3 pr-4 font-bold text-[var(--color-ink)]">{customer.name}</td>
                    <td className="border-b border-[var(--color-line)] px-4 py-3 text-[var(--color-ink-soft)]">{customer.phone || '—'}</td>
                    <td className="border-b border-[var(--color-line)] px-4 py-3 text-[var(--color-ink-soft)]">{customer.email || '—'}</td>
                    <td className="border-b border-[var(--color-line)] py-3 text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          onClick={() => {
                            setEditingCustomer(customer)
                            setIsFormOpen(true)
                          }}
                          size="sm"
                          variant="secondary"
                        >
                          Editar
                        </Button>
                        <Button onClick={() => setCustomerToDelete(customer)} size="sm" variant="danger">
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
        <CustomerFormModal
          customer={editingCustomer}
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadCustomers()
          }}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O cliente "${customerToDelete?.name ?? ''}" será removido do cadastro.`}
        onCancel={() => setCustomerToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(customerToDelete)}
        title="Remover cliente?"
      />
    </PageContainer>
  )
}

function CustomerFormModal({
  customer,
  onClose,
  onSaved,
}: {
  customer: Customer | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(customer?.name ?? '')
  const [phone, setPhone] = useState(customer?.phone ?? '')
  const [email, setEmail] = useState(customer?.email ?? '')
  const [birthday, setBirthday] = useState(customer?.birthday ?? '')
  const [notes, setNotes] = useState(customer?.notes ?? '')
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome do cliente.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = { name: name.trim(), phone: phone.trim(), email: email.trim(), birthday, notes: notes.trim() }
      if (customer) {
        await updateCustomer(customer.id, input)
        showToast({ message: 'Cliente atualizado.', tone: 'success' })
      } else {
        await createCustomer(input)
        showToast({ message: 'Cliente criado.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="customer-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={customer ? 'Atualize os dados do cliente.' : 'Cadastre um novo cliente.'}
      title={customer ? 'Editar cliente' : 'Novo cliente'}
    >
      <Form id="customer-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <Field label="Telefone">
          <Input className="mt-2" onChange={(event) => setPhone(event.target.value)} placeholder="(11) 98888-7777" value={phone} />
        </Field>
        <Field label="Email">
          <Input className="mt-2" onChange={(event) => setEmail(event.target.value)} type="email" value={email} />
        </Field>
        <Field label="Aniversário">
          <Input className="mt-2" onChange={(event) => setBirthday(event.target.value)} type="date" value={birthday} />
        </Field>
        <Field label="Observações">
          <Textarea className="mt-2" onChange={(event) => setNotes(event.target.value)} value={notes} />
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
