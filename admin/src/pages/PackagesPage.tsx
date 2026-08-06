import { Package as PackageIcon, Plus, TicketPercent } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createPackage,
  createPackageCredit,
  deletePackage,
  fetchPackageCredits,
  fetchPackages,
  setPackageActive,
  updatePackage,
} from '../lib/adminApi'
import { formatMoney } from '../lib/format'
import { useReferenceData } from '../shared/hooks/useReferenceData'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { ConfirmModal } from '../shared/components/ConfirmModal'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select, Textarea } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { Package, PackageCredit } from '../types'

export function PackagesPage() {
  const { showToast } = useToast()
  const referenceData = useReferenceData()
  const [packages, setPackages] = useState<Package[]>([])
  const [packageCredits, setPackageCredits] = useState<PackageCredit[]>([])
  const [loading, setLoading] = useState(true)
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [isCreditOpen, setIsCreditOpen] = useState(false)
  const [editingPackage, setEditingPackage] = useState<Package | null>(null)
  const [packageToDelete, setPackageToDelete] = useState<Package | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  function loadPackages() {
    setLoading(true)
    Promise.allSettled([fetchPackages(), fetchPackageCredits()])
      .then(([packagesResult, creditsResult]) => {
        setPackages(packagesResult.status === 'fulfilled' ? packagesResult.value : [])
        setPackageCredits(creditsResult.status === 'fulfilled' ? creditsResult.value : [])
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os pacotes.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadPackages, [])

  async function handleToggleActive(item: Package) {
    try {
      await setPackageActive(item.id, !item.active)
      showToast({ message: item.active ? 'Pacote desativado.' : 'Pacote ativado.', tone: 'success' })
      loadPackages()
    } catch (error) {
      showToast({ message: (error as Error).message, tone: 'error' })
    }
  }

  async function handleDelete() {
    if (!packageToDelete) {
      return
    }
    setIsDeleting(true)
    try {
      await deletePackage(packageToDelete.id)
      showToast({ message: 'Pacote removido.', tone: 'success' })
      setPackageToDelete(null)
      loadPackages()
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
          <div className="flex flex-wrap gap-2">
            <Button icon={<TicketPercent className="h-4 w-4" />} onClick={() => setIsCreditOpen(true)} variant="secondary">
              Atribuir créditos
            </Button>
            <Button
              icon={<Plus className="h-4 w-4" />}
              onClick={() => {
                setEditingPackage(null)
                setIsFormOpen(true)
              }}
            >
              Novo pacote
            </Button>
          </div>
        }
        subtitle="Pacotes de sessões e créditos atribuídos aos clientes."
        title="Pacotes"
      >
        <div className="mt-2" />
      </Section>

      {loading ? (
        <Section>
          <Spinner />
        </Section>
      ) : packages.length === 0 ? (
        <Section>
          <EmptyState
            description="Crie pacotes de sessões para vender em blocos para os clientes."
            icon={<PackageIcon className="h-10 w-10" />}
            title="Nenhum pacote cadastrado"
          />
        </Section>
      ) : (
        <Section>
          <ul className="mt-2 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
            {packages.map((item) => (
              <li
                className="flex flex-col gap-3 rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4"
                key={item.id}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-bold text-[var(--color-ink)]">{item.name}</p>
                    <p className="mt-0.5 text-sm font-semibold text-[var(--color-ink-soft)]">
                      {item.sessionsCount} sessões{item.durationDays ? ` · válido por ${item.durationDays} dias` : ''}
                    </p>
                  </div>
                  <Badge tone={item.active ? 'success' : 'danger'}>{item.active ? 'Ativo' : 'Inativo'}</Badge>
                </div>
                {item.description ? (
                  <p className="line-clamp-2 text-sm text-[var(--color-ink-soft)]">{item.description}</p>
                ) : null}
                <p className="font-headline text-lg font-black text-[var(--color-ink)]">{formatMoney(item.priceCents)}</p>
                <div className="mt-auto flex flex-wrap items-center gap-2">
                  <Button
                    onClick={() => {
                      setEditingPackage(item)
                      setIsFormOpen(true)
                    }}
                    size="sm"
                    variant="secondary"
                  >
                    Editar
                  </Button>
                  <Button onClick={() => void handleToggleActive(item)} size="sm" variant="secondary">
                    {item.active ? 'Desativar' : 'Ativar'}
                  </Button>
                  <Button onClick={() => setPackageToDelete(item)} size="sm" variant="danger">
                    Remover
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        </Section>
      )}

      <Section subtitle="Sessões de pacotes atribuídas a cada cliente." title="Créditos de pacotes">
        {packageCredits.length === 0 ? (
          <p className="mt-3 text-sm font-semibold text-[var(--color-ink-soft)]">
            Nenhum crédito atribuído ainda. Use "Atribuir créditos" para registrar uma compra.
          </p>
        ) : (
          <ul className="mt-4 grid gap-2">
            {packageCredits.map((credit) => (
              <li
                className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] px-4 py-3"
                key={credit.id}
              >
                <div className="min-w-0">
                  <p className="truncate font-bold text-[var(--color-ink)]">{credit.customerName ?? credit.customerId}</p>
                  <p className="truncate text-sm text-[var(--color-ink-soft)]">{credit.packageName ?? credit.packageId}</p>
                </div>
                <Badge tone="accent">
                  {credit.sessionsUsed}/{credit.sessionsTotal} sessões
                </Badge>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {isFormOpen ? (
        <PackageFormModal
          onClose={() => setIsFormOpen(false)}
          onSaved={() => {
            setIsFormOpen(false)
            loadPackages()
          }}
          item={editingPackage}
        />
      ) : null}

      {isCreditOpen ? (
        <AssignCreditModal
          customers={referenceData.customers}
          onClose={() => setIsCreditOpen(false)}
          onSaved={() => {
            setIsCreditOpen(false)
            loadPackages()
          }}
          packages={packages}
        />
      ) : null}

      <ConfirmModal
        confirmLabel="Remover"
        loading={isDeleting}
        message={`O pacote "${packageToDelete?.name ?? ''}" será removido.`}
        onCancel={() => setPackageToDelete(null)}
        onConfirm={() => void handleDelete()}
        open={Boolean(packageToDelete)}
        title="Remover pacote?"
      />
    </PageContainer>
  )
}

function PackageFormModal({
  item,
  onClose,
  onSaved,
}: {
  item: Package | null
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [name, setName] = useState(item?.name ?? '')
  const [description, setDescription] = useState(item?.description ?? '')
  const [price, setPrice] = useState(item ? String(item.priceCents / 100) : '')
  const [sessionsCount, setSessionsCount] = useState(item?.sessionsCount ?? 1)
  const [durationDays, setDurationDays] = useState(item?.durationDays ?? 0)
  const [active, setActive] = useState(item?.active ?? true)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!name.trim()) {
      setError('Informe o nome do pacote.')
      return
    }
    const priceCents = Math.round(Number(price) * 100)
    if (!Number.isFinite(priceCents) || priceCents <= 0) {
      setError('Informe um preço válido.')
      return
    }
    if (sessionsCount < 1) {
      setError('O pacote precisa ter pelo menos 1 sessão.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      const input = {
        name: name.trim(),
        description: description.trim(),
        priceCents,
        sessionsCount,
        durationDays,
        active,
      }
      if (item) {
        await updatePackage(item.id, input)
        showToast({ message: 'Pacote atualizado.', tone: 'success' })
      } else {
        await createPackage(input)
        showToast({ message: 'Pacote criado.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="package-form">
            Salvar
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle={item ? 'Atualize os dados do pacote.' : 'Cadastre um novo pacote de sessões.'}
      title={item ? 'Editar pacote' : 'Novo pacote'}
    >
      <Form id="package-form" onSubmit={handleSubmit}>
        <Field label="Nome">
          <Input className="mt-2" onChange={(event) => setName(event.target.value)} value={name} />
        </Field>
        <Field label="Descrição">
          <Textarea className="mt-2" onChange={(event) => setDescription(event.target.value)} value={description} />
        </Field>
        <div className="grid gap-4 md:grid-cols-3">
          <Field label="Preço (R$)">
            <Input
              className="mt-2"
              min="0"
              onChange={(event) => setPrice(event.target.value)}
              step="0.01"
              type="number"
              value={price}
            />
          </Field>
          <Field label="Sessões">
            <Input
              className="mt-2"
              min={1}
              onChange={(event) => setSessionsCount(Number(event.target.value))}
              type="number"
              value={sessionsCount}
            />
          </Field>
          <Field hint="0 = sem validade." label="Validade (dias)">
            <Input
              className="mt-2"
              min={0}
              onChange={(event) => setDurationDays(Number(event.target.value))}
              type="number"
              value={durationDays}
            />
          </Field>
        </div>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}

function AssignCreditModal({
  packages,
  customers,
  onClose,
  onSaved,
}: {
  packages: Package[]
  customers: { id: string; name: string }[]
  onClose: () => void
  onSaved: () => void
}) {
  const { showToast } = useToast()
  const [packageId, setPackageId] = useState(packages[0]?.id ?? '')
  const [customerId, setCustomerId] = useState(customers[0]?.id ?? '')
  const [sessionsTotal, setSessionsTotal] = useState(packages[0]?.sessionsCount ?? 1)
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  function handlePackageChange(value: string) {
    setPackageId(value)
    const selected = packages.find((item) => item.id === value)
    if (selected) {
      setSessionsTotal(selected.sessionsCount)
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!packageId || !customerId) {
      setError('Selecione o pacote e o cliente.')
      return
    }
    if (sessionsTotal < 1) {
      setError('Informe a quantidade de sessões.')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      await createPackageCredit({ packageId, customerId, sessionsTotal })
      showToast({ message: 'Créditos atribuídos ao cliente.', tone: 'success' })
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
          <Button loading={isSubmitting} type="submit" form="credit-form">
            Atribuir
          </Button>
        </>
      }
      onClose={onClose}
      open
      size="md"
      subtitle="Registre a compra de um pacote atribuindo as sessões a um cliente."
      title="Atribuir créditos de pacote"
    >
      <Form id="credit-form" onSubmit={handleSubmit}>
        <Field label="Pacote">
          <Select className="mt-2" onChange={(event) => handlePackageChange(event.target.value)} value={packageId}>
            {packages.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name} · {formatMoney(item.priceCents)}
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Cliente">
          <Select className="mt-2" onChange={(event) => setCustomerId(event.target.value)} value={customerId}>
            {customers.map((customer) => (
              <option key={customer.id} value={customer.id}>
                {customer.name}
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Sessões">
          <Input
            className="mt-2"
            min={1}
            onChange={(event) => setSessionsTotal(Number(event.target.value))}
            type="number"
            value={sessionsTotal}
          />
        </Field>
        {error ? <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p> : null}
      </Form>
    </Modal>
  )
}
