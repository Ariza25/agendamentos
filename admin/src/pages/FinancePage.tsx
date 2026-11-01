import {
  AlertTriangle,
  ArrowDownLeft,
  ArrowUpRight,
  CheckCircle,
  Clock,
  DollarSign,
  History,
  Lock,
  Plus,
  Receipt,
  Trash2,
  Unlock,
  Wallet,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  closeCashRegister,
  createCashRegisterEntry,
  createFinancialTransaction,
  deleteFinancialTransaction,
  fetchCashRegisterHistory,
  fetchCurrentCashRegister,
  fetchFinancialSummary,
  fetchFinancialTransactions,
  openCashRegister,
  updateFinancialTransaction,
} from '../lib/adminApi'
import { formatDateTime, formatMoney, formatTime } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type {
  CashRegister,
  CashRegisterEntryType,
  FinancialSummary,
  FinancialTransaction,
  ShiftSummary,
  TransactionStatus,
  TransactionType,
} from '../types'

const STATUS_LABELS: Record<TransactionStatus, string> = {
  pending: 'Pendente',
  paid: 'Pago',
  cancelled: 'Cancelado',
}

function statusTone(status: TransactionStatus): 'neutral' | 'success' | 'warning' | 'danger' {
  if (status === 'paid') return 'success'
  if (status === 'pending') return 'warning'
  return 'danger'
}

type ActiveTab = 'transactions' | 'cash_register'

export function FinancePage() {
  const { showToast } = useToast()
  const [activeTab, setActiveTab] = useState<ActiveTab>('transactions')

  // Transactions state
  const [transactions, setTransactions] = useState<FinancialTransaction[]>([])
  const [summary, setSummary] = useState<FinancialSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [filterType, setFilterType] = useState<string>('')
  const [filterStatus, setFilterStatus] = useState<string>('')

  // Transaction form state
  const [formType, setFormType] = useState<TransactionType>('income')
  const [formCategory, setFormCategory] = useState('Serviços')
  const [formDescription, setFormDescription] = useState('')
  const [formAmount, setFormAmount] = useState('')
  const [formDueDate, setFormDueDate] = useState(() => new Date().toISOString().slice(0, 10))
  const [formStatus, setFormStatus] = useState<TransactionStatus>('paid')
  const [formPaymentMethod, setFormPaymentMethod] = useState('pix')
  const [saving, setSaving] = useState(false)

  // Cash Register (Caixa Diário) state
  const [shiftSummary, setShiftSummary] = useState<ShiftSummary | null>(null)
  const [isShiftOpen, setIsShiftOpen] = useState(false)
  const [shiftHistory, setShiftHistory] = useState<CashRegister[]>([])
  const [loadingShift, setLoadingShift] = useState(false)

  // Cash Register Modals
  const [isOpenModalOpen, setIsOpenModalOpen] = useState(false)
  const [initialFloatInput, setInitialFloatInput] = useState('100.00')
  const [openNotes, setOpenNotes] = useState('')
  const [submittingOpen, setSubmittingOpen] = useState(false)

  const [isEntryModalOpen, setIsEntryModalOpen] = useState(false)
  const [entryType, setEntryType] = useState<CashRegisterEntryType>('sangria')
  const [entryAmountInput, setEntryAmountInput] = useState('')
  const [entryReason, setEntryReason] = useState('')
  const [submittingEntry, setSubmittingEntry] = useState(false)

  const [isCloseModalOpen, setIsCloseModalOpen] = useState(false)
  const [closingCountInput, setClosingCountInput] = useState('')
  const [closingNotes, setClosingNotes] = useState('')
  const [submittingClose, setSubmittingClose] = useState(false)

  function loadData() {
    setLoading(true)
    Promise.all([
      fetchFinancialTransactions({
        type: filterType || undefined,
        status: filterStatus || undefined,
      }),
      fetchFinancialSummary(),
    ])
      .then(([txList, sumData]) => {
        setTransactions(txList)
        setSummary(sumData)
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar as finanças.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  function loadCashRegister() {
    setLoadingShift(true)
    Promise.all([fetchCurrentCashRegister(), fetchCashRegisterHistory()])
      .then(([currentRes, history]) => {
        setIsShiftOpen(currentRes.open && !!currentRes.summary)
        setShiftSummary(currentRes.summary || null)
        setShiftHistory(history)
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os dados do caixa.', tone: 'error' })
      })
      .finally(() => setLoadingShift(false))
  }

  useEffect(() => {
    loadData()
  }, [filterType, filterStatus])

  useEffect(() => {
    if (activeTab === 'cash_register') {
      loadCashRegister()
    }
  }, [activeTab])

  function resetForm() {
    setFormType('income')
    setFormCategory('Serviços')
    setFormDescription('')
    setFormAmount('')
    setFormDueDate(new Date().toISOString().slice(0, 10))
    setFormStatus('paid')
    setFormPaymentMethod('pix')
  }

  async function handleCreate(e: FormEvent) {
    e.preventDefault()
    const amountVal = parseFloat(formAmount.replace(',', '.'))
    if (isNaN(amountVal) || amountVal <= 0) {
      showToast({ message: 'Informe um valor válido maior que zero.', tone: 'error' })
      return
    }

    setSaving(true)
    try {
      await createFinancialTransaction({
        type: formType,
        category: formCategory,
        description: formDescription.trim(),
        amountCents: Math.round(amountVal * 100),
        dueDate: formDueDate,
        status: formStatus,
        paymentMethod: formPaymentMethod,
      })
      showToast({ message: 'Transação registrada com sucesso!', tone: 'success' })
      setIsModalOpen(false)
      resetForm()
      loadData()
      if (activeTab === 'cash_register') loadCashRegister()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSaving(false)
    }
  }

  async function handleMarkPaid(tx: FinancialTransaction) {
    try {
      await updateFinancialTransaction(tx.id, {
        status: 'paid',
        paidAt: new Date().toISOString(),
      })
      showToast({ message: 'Transação marcada como paga.', tone: 'success' })
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Deseja excluir esta transação?')) return
    try {
      await deleteFinancialTransaction(id)
      showToast({ message: 'Transação excluída.', tone: 'success' })
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    }
  }

  // --- Cash Register Actions ---
  async function handleOpenShift(e: FormEvent) {
    e.preventDefault()
    const initialVal = parseFloat(initialFloatInput.replace(',', '.'))
    if (isNaN(initialVal) || initialVal < 0) {
      showToast({ message: 'Informe um valor inicial válido.', tone: 'error' })
      return
    }

    setSubmittingOpen(true)
    try {
      await openCashRegister(Math.round(initialVal * 100), openNotes)
      showToast({ message: 'Caixa aberto com sucesso!', tone: 'success' })
      setIsOpenModalOpen(false)
      loadCashRegister()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSubmittingOpen(false)
    }
  }

  async function handleCreateEntry(e: FormEvent) {
    e.preventDefault()
    const val = parseFloat(entryAmountInput.replace(',', '.'))
    if (isNaN(val) || val <= 0) {
      showToast({ message: 'Informe um valor maior que zero.', tone: 'error' })
      return
    }
    if (!entryReason.trim()) {
      showToast({ message: 'Informe o motivo da movimentação.', tone: 'error' })
      return
    }

    setSubmittingEntry(true)
    try {
      await createCashRegisterEntry(entryType, Math.round(val * 100), entryReason.trim())
      showToast({
        message: `${entryType === 'sangria' ? 'Sangria' : 'Suprimento'} registrado com sucesso!`,
        tone: 'success',
      })
      setIsEntryModalOpen(false)
      loadCashRegister()
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSubmittingEntry(false)
    }
  }

  async function handleCloseShift(e: FormEvent) {
    e.preventDefault()
    const countVal = parseFloat(closingCountInput.replace(',', '.'))
    if (isNaN(countVal) || countVal < 0) {
      showToast({ message: 'Informe a contagem física do dinheiro em caixa.', tone: 'error' })
      return
    }

    setSubmittingClose(true)
    try {
      const closed = await closeCashRegister(Math.round(countVal * 100), closingNotes)
      const diffCents = closed.differenceCents
      let diffMsg = 'Caixa fechado com sucesso.'
      if (diffMsg && diffCents > 0) {
        diffMsg += ` Sobra identificada: ${formatMoney(diffCents)}.`
      } else if (diffCents < 0) {
        diffMsg += ` Quebra/falta identificada: ${formatMoney(Math.abs(diffCents))}.`
      }
      showToast({ message: diffMsg, tone: 'success' })
      setIsCloseModalOpen(false)
      loadCashRegister()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSubmittingClose(false)
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          activeTab === 'transactions' ? (
            <Button
              onClick={() => {
                resetForm()
                setIsModalOpen(true)
              }}
            >
              <Plus className="mr-2 h-4 w-4" /> Nova Transação
            </Button>
          ) : isShiftOpen ? (
            <div className="flex items-center gap-2">
              <Button
                variant="secondary"
                size="sm"
                onClick={() => {
                  setEntryType('sangria')
                  setEntryAmountInput('')
                  setEntryReason('')
                  setIsEntryModalOpen(true)
                }}
              >
                <ArrowUpRight className="mr-1 h-3.5 w-3.5 text-rose-500" />
                Sangria
              </Button>
              <Button
                variant="secondary"
                size="sm"
                onClick={() => {
                  setEntryType('suprimento')
                  setEntryAmountInput('')
                  setEntryReason('')
                  setIsEntryModalOpen(true)
                }}
              >
                <ArrowDownLeft className="mr-1 h-3.5 w-3.5 text-emerald-500" />
                Suprimento
              </Button>
              <Button
                size="sm"
                className="bg-slate-900 text-white hover:bg-black"
                onClick={() => {
                  setClosingCountInput(
                    shiftSummary ? (shiftSummary.currentDrawerCents / 100).toFixed(2) : '0.00'
                  )
                  setClosingNotes('')
                  setIsCloseModalOpen(true)
                }}
              >
                <Lock className="mr-1 h-3.5 w-3.5" />
                Fechar Caixa
              </Button>
            </div>
          ) : (
            <Button
              className="bg-emerald-600 text-white hover:bg-emerald-700"
              onClick={() => {
                setInitialFloatInput('100.00')
                setOpenNotes('')
                setIsOpenModalOpen(true)
              }}
            >
              <Unlock className="mr-1 h-3.5 w-3.5" />
              Abrir Caixa Hoje
            </Button>
          )
        }
        subtitle="Controle fluxo de caixa, despesas, turnos de caixa diário e conciliação de gaveta."
        title="Financeiro & Caixa Diário"
      >
        {/* Navigation Tabs */}
        <div className="mt-4 flex border-b border-border">
          <button
            onClick={() => setActiveTab('transactions')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'transactions'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <DollarSign className="h-4 w-4" />
            Lançamentos & DRE
          </button>
          <button
            onClick={() => setActiveTab('cash_register')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'cash_register'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <Wallet className="h-4 w-4" />
            Caixa Diário (Frente)
            {isShiftOpen ? (
              <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 text-xs font-bold text-emerald-600">
                <span className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
                Aberto
              </span>
            ) : (
              <span className="rounded-full bg-slate-500/10 px-2 py-0.5 text-xs font-semibold text-slate-500">
                Fechado
              </span>
            )}
          </button>
        </div>
      </Section>

      {activeTab === 'transactions' ? (
        <>
          {/* KPI Cards */}
          {summary && (
            <Section title="Resumo do Período">
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-muted-foreground">Total Receitas</span>
                    <span className="rounded-lg bg-emerald-500/10 p-2 text-emerald-600 dark:text-emerald-400">
                      <ArrowDownLeft className="h-4 w-4" />
                    </span>
                  </div>
                  <p className="mt-2 text-2xl font-bold text-emerald-600 dark:text-emerald-400">
                    {formatMoney(summary.totalIncomeCents)}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    Pendente: {formatMoney(summary.pendingIncomeCents)}
                  </p>
                </div>

                <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-muted-foreground">Total Despesas</span>
                    <span className="rounded-lg bg-rose-500/10 p-2 text-rose-600 dark:text-rose-400">
                      <ArrowUpRight className="h-4 w-4" />
                    </span>
                  </div>
                  <p className="mt-2 text-2xl font-bold text-rose-600 dark:text-rose-400">
                    {formatMoney(summary.totalExpenseCents)}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    A pagar: {formatMoney(summary.pendingExpenseCents)}
                  </p>
                </div>

                <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-muted-foreground">Lucro Líquido</span>
                    <span className="rounded-lg bg-primary/10 p-2 text-primary">
                      <DollarSign className="h-4 w-4" />
                    </span>
                  </div>
                  <p
                    className={`mt-2 text-2xl font-bold ${
                      summary.netProfitCents >= 0 ? 'text-primary' : 'text-rose-600'
                    }`}
                  >
                    {formatMoney(summary.netProfitCents)}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    Receitas realizadas - Despesas pagas
                  </p>
                </div>

                <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-muted-foreground">Saldo Pendente</span>
                    <span className="rounded-lg bg-amber-500/10 p-2 text-amber-600">
                      <DollarSign className="h-4 w-4" />
                    </span>
                  </div>
                  <p className="mt-2 text-2xl font-bold text-foreground">
                    {formatMoney(summary.pendingIncomeCents - summary.pendingExpenseCents)}
                  </p>
                  <p className="mt-1 text-xs text-muted-foreground">Previsão de recebíveis futuros</p>
                </div>
              </div>
            </Section>
          )}

          {/* Filter and Transactions List */}
          <Section
            actions={
              <div className="flex flex-wrap items-center gap-2">
                <Select
                  className="w-36 text-sm"
                  onChange={(e) => setFilterType(e.target.value)}
                  value={filterType}
                >
                  <option value="">Tipo: Todos</option>
                  <option value="income">Receitas</option>
                  <option value="expense">Despesas</option>
                </Select>

                <Select
                  className="w-36 text-sm"
                  onChange={(e) => setFilterStatus(e.target.value)}
                  value={filterStatus}
                >
                  <option value="">Status: Todos</option>
                  <option value="paid">Pago</option>
                  <option value="pending">Pendente</option>
                  <option value="cancelled">Cancelado</option>
                </Select>
              </div>
            }
            title="Lançamentos Financeiros"
          >
            {loading ? (
              <Spinner />
            ) : transactions.length === 0 ? (
              <EmptyState
                description="Nenhuma transação financeira encontrada para os filtros selecionados."
                title="Nenhum lançamento"
              />
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full border-collapse text-left text-sm">
                  <thead>
                    <tr className="border-b border-border text-xs uppercase text-muted-foreground">
                      <th className="py-3 px-4">Tipo</th>
                      <th className="py-3 px-4">Categoria</th>
                      <th className="py-3 px-4">Descrição</th>
                      <th className="py-3 px-4">Vencimento</th>
                      <th className="py-3 px-4">Status</th>
                      <th className="py-3 px-4">Forma</th>
                      <th className="py-3 px-4 text-right">Valor</th>
                      <th className="py-3 px-4 text-right">Ações</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {transactions.map((tx) => (
                      <tr className="hover:bg-muted/50 transition-colors" key={tx.id}>
                        <td className="py-3 px-4">
                          {tx.type === 'income' ? (
                            <span className="inline-flex items-center gap-1 font-semibold text-emerald-600 dark:text-emerald-400">
                              <ArrowDownLeft className="h-3.5 w-3.5" /> Receita
                            </span>
                          ) : (
                            <span className="inline-flex items-center gap-1 font-semibold text-rose-600 dark:text-rose-400">
                              <ArrowUpRight className="h-3.5 w-3.5" /> Despesa
                            </span>
                          )}
                        </td>
                        <td className="py-3 px-4 font-medium text-foreground">{tx.category}</td>
                        <td className="py-3 px-4 text-muted-foreground">{tx.description}</td>
                        <td className="py-3 px-4 text-muted-foreground">
                          {new Date(tx.dueDate).toLocaleDateString('pt-BR')}
                        </td>
                        <td className="py-3 px-4">
                          <Badge tone={statusTone(tx.status)}>{STATUS_LABELS[tx.status]}</Badge>
                        </td>
                        <td className="py-3 px-4 uppercase text-xs text-muted-foreground">
                          {tx.paymentMethod || '—'}
                        </td>
                        <td className="py-3 px-4 text-right font-bold">
                          <span
                            className={
                              tx.type === 'income'
                                ? 'text-emerald-600 dark:text-emerald-400'
                                : 'text-rose-600 dark:text-rose-400'
                            }
                          >
                            {tx.type === 'income' ? '+' : '-'} {formatMoney(tx.amountCents)}
                          </span>
                        </td>
                        <td className="py-3 px-4 text-right">
                          <div className="flex items-center justify-end gap-1">
                            {tx.status === 'pending' && (
                              <Button
                                aria-label="Marcar como pago"
                                onClick={() => handleMarkPaid(tx)}
                                size="sm"
                                variant="secondary"
                              >
                                <CheckCircle className="h-3.5 w-3.5 text-emerald-600" />
                              </Button>
                            )}
                            <Button
                              aria-label="Excluir transação"
                              onClick={() => handleDelete(tx.id)}
                              size="sm"
                              variant="danger"
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </Section>
        </>
      ) : (
        /* Cash Register Tab */
        <div className="space-y-6">
          {loadingShift ? (
            <Section>
              <Spinner />
            </Section>
          ) : !isShiftOpen || !shiftSummary ? (
            <Section>
              <EmptyState
                description="O caixa diário deste estabelecimento está fechado no momento. Abra o caixa informando o troco inicial para registrar vendas e gerenciar a gaveta."
                title="Caixa Diário Fechado"
              />
              <div className="mt-4 flex justify-center">
                <Button
                  className="bg-emerald-600 text-white hover:bg-emerald-700"
                  onClick={() => {
                    setInitialFloatInput('100.00')
                    setOpenNotes('')
                    setIsOpenModalOpen(true)
                  }}
                >
                  <Unlock className="mr-2 h-4 w-4" />
                  Abrir Caixa Hoje
                </Button>
              </div>
            </Section>
          ) : (
            <>
              {/* Caixa Aberto KPIs */}
              <Section
                subtitle={`Turno iniciado em ${formatDateTime(
                  shiftSummary.cashRegister.openedAt
                )} · Usuário: ${shiftSummary.cashRegister.openedByUserId || 'Admin'}`}
                title="Gaveta e Saldo em Espécie do Turno"
              >
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
                  <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                    <span className="text-xs font-semibold text-muted-foreground">Fundo Inicial (Troco)</span>
                    <p className="mt-2 text-xl font-bold text-foreground">
                      {formatMoney(shiftSummary.cashRegister.initialBalanceCents)}
                    </p>
                  </div>

                  <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                    <span className="text-xs font-semibold text-muted-foreground">Entradas Dinheiro (Vendas)</span>
                    <p className="mt-2 text-xl font-bold text-emerald-600">
                      +{formatMoney(shiftSummary.totalCashIncomeCents)}
                    </p>
                  </div>

                  <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                    <span className="text-xs font-semibold text-muted-foreground">Suprimentos (Aportes)</span>
                    <p className="mt-2 text-xl font-bold text-blue-600">
                      +{formatMoney(shiftSummary.totalSuprimentosCents)}
                    </p>
                  </div>

                  <div className="rounded-xl border border-border bg-card p-4 shadow-sm">
                    <span className="text-xs font-semibold text-muted-foreground">Sangrias (Retiradas)</span>
                    <p className="mt-2 text-xl font-bold text-rose-600">
                      -{formatMoney(shiftSummary.totalSangriasCents)}
                    </p>
                  </div>

                  <div className="rounded-xl border-2 border-emerald-500 bg-emerald-50/50 p-4 shadow-sm dark:bg-emerald-950/20">
                    <span className="text-xs font-black uppercase tracking-wider text-emerald-700 dark:text-emerald-400">
                      Saldo Atual da Gaveta
                    </span>
                    <p className="mt-2 text-2xl font-black text-emerald-700 dark:text-emerald-400">
                      {formatMoney(shiftSummary.currentDrawerCents)}
                    </p>
                  </div>
                </div>

                {/* Additional non-cash totals */}
                <div className="mt-4 rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
                  <span className="text-xs font-bold uppercase tracking-wider text-[var(--color-ink-soft)]">
                    Outros Recebimentos no Turno (Cartões / Pix)
                  </span>
                  <div className="mt-2 flex flex-wrap items-center gap-6 text-sm font-semibold text-[var(--color-ink)]">
                    <div>
                      Cartões: <span className="font-bold">{formatMoney(shiftSummary.totalCardIncomeCents)}</span>
                    </div>
                    <div>
                      Pix: <span className="font-bold">{formatMoney(shiftSummary.totalPixIncomeCents)}</span>
                    </div>
                    {shiftSummary.totalOtherIncomeCents > 0 && (
                      <div>
                        Outros: <span className="font-bold">{formatMoney(shiftSummary.totalOtherIncomeCents)}</span>
                      </div>
                    )}
                  </div>
                </div>
              </Section>

              {/* Sangrias & Suprimentos do Turno */}
              <Section
                actions={
                  <div className="flex items-center gap-2">
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => {
                        setEntryType('sangria')
                        setEntryAmountInput('')
                        setEntryReason('')
                        setIsEntryModalOpen(true)
                      }}
                    >
                      <ArrowUpRight className="mr-1 h-3.5 w-3.5 text-rose-500" />
                      Sangria
                    </Button>
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => {
                        setEntryType('suprimento')
                        setEntryAmountInput('')
                        setEntryReason('')
                        setIsEntryModalOpen(true)
                      }}
                    >
                      <ArrowDownLeft className="mr-1 h-3.5 w-3.5 text-emerald-500" />
                      Suprimento
                    </Button>
                  </div>
                }
                subtitle="Registro de retiradas e aportes realizados durante o turno atual."
                title="Movimentações de Gaveta do Turno"
              >
                {shiftSummary.entries.length === 0 ? (
                  <p className="text-sm font-semibold text-muted-foreground">
                    Nenhuma sangria ou suprimento registrado neste turno.
                  </p>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full border-collapse text-left text-sm">
                      <thead>
                        <tr className="border-b border-border text-xs uppercase text-muted-foreground">
                          <th className="py-2.5 px-3">Hora</th>
                          <th className="py-2.5 px-3">Tipo</th>
                          <th className="py-2.5 px-3">Motivo / Descrição</th>
                          <th className="py-2.5 px-3 text-right">Valor</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-border">
                        {shiftSummary.entries.map((entry) => (
                          <tr key={entry.id} className="hover:bg-muted/50">
                            <td className="py-2.5 px-3 text-xs text-muted-foreground font-mono">
                              {formatTime(entry.createdAt)}
                            </td>
                            <td className="py-2.5 px-3">
                              <Badge tone={entry.type === 'sangria' ? 'danger' : 'accent'}>
                                {entry.type === 'sangria' ? 'Sangria (Saída)' : 'Suprimento (Aporte)'}
                              </Badge>
                            </td>
                            <td className="py-2.5 px-3 font-medium text-foreground">{entry.reason}</td>
                            <td className="py-2.5 px-3 text-right font-bold">
                              <span className={entry.type === 'sangria' ? 'text-rose-600' : 'text-blue-600'}>
                                {entry.type === 'sangria' ? '-' : '+'} {formatMoney(entry.amountCents)}
                              </span>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </Section>
            </>
          )}

          {/* Histórico de Turnos Fechados */}
          <Section
            subtitle="Histórico de auditoria dos últimos caixas encerrados com contagem e quebra/sobra."
            title="Histórico de Turnos Anteriores"
          >
            {shiftHistory.length === 0 ? (
              <p className="text-sm font-semibold text-muted-foreground">
                Nenhum turno anterior registrado ainda.
              </p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full border-collapse text-left text-sm">
                  <thead>
                    <tr className="border-b border-border text-xs uppercase text-muted-foreground">
                      <th className="py-3 px-3">Abertura</th>
                      <th className="py-3 px-3">Fechamento</th>
                      <th className="py-3 px-3">Fundo Inicial</th>
                      <th className="py-3 px-3">Saldo Sistema</th>
                      <th className="py-3 px-3">Contagem Física</th>
                      <th className="py-3 px-3 text-right">Diferença</th>
                      <th className="py-3 px-3">Observações</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {shiftHistory.map((shift) => (
                      <tr key={shift.id} className="hover:bg-muted/50">
                        <td className="py-2.5 px-3 text-xs text-foreground font-medium">
                          {formatDateTime(shift.openedAt)}
                        </td>
                        <td className="py-2.5 px-3 text-xs text-muted-foreground">
                          {shift.closedAt ? formatDateTime(shift.closedAt) : 'Em aberto'}
                        </td>
                        <td className="py-2.5 px-3">{formatMoney(shift.initialBalanceCents)}</td>
                        <td className="py-2.5 px-3 font-semibold">{formatMoney(shift.systemBalanceCents)}</td>
                        <td className="py-2.5 px-3 font-bold">{formatMoney(shift.closingBalanceCents)}</td>
                        <td className="py-2.5 px-3 text-right">
                          {shift.differenceCents === 0 ? (
                            <Badge tone="success">Caixa Exato (R$ 0,00)</Badge>
                          ) : shift.differenceCents > 0 ? (
                            <Badge tone="accent">+{formatMoney(shift.differenceCents)} sobra</Badge>
                          ) : (
                            <Badge tone="danger">-{formatMoney(Math.abs(shift.differenceCents))} falta</Badge>
                          )}
                        </td>
                        <td className="py-2.5 px-3 text-xs text-muted-foreground italic">
                          {shift.notes || '—'}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </Section>
        </div>
      )}

      {/* Modal: Nova Transação */}
      <Modal
        onClose={() => setIsModalOpen(false)}
        open={isModalOpen}
        title="Novo Lançamento Financeiro"
      >
        <Form className="grid gap-4" onSubmit={handleCreate}>
          <Field label="Tipo de Lançamento">
            <Select
              onChange={(e) => {
                const t = e.target.value as TransactionType
                setFormType(t)
                setFormCategory(t === 'income' ? 'Serviços' : 'Suprimentos')
              }}
              value={formType}
            >
              <option value="income">Receita (Entrada)</option>
              <option value="expense">Despesa (Saída)</option>
            </Select>
          </Field>

          <Field label="Categoria">
            <Select onChange={(e) => setFormCategory(e.target.value)} value={formCategory}>
              {formType === 'income' ? (
                <>
                  <option value="Serviços">Serviços e Agendamentos</option>
                  <option value="Vendas">Venda de Produtos</option>
                  <option value="Planos">Assinaturas e Mensalidades</option>
                  <option value="Outros">Outras Receitas</option>
                </>
              ) : (
                <>
                  <option value="Suprimentos">Suprimentos e Estoque</option>
                  <option value="Comissões">Comissões de Profissionais</option>
                  <option value="Aluguel">Aluguel e Condomínio</option>
                  <option value="Utilidades">Água / Luz / Internet</option>
                  <option value="Marketing">Marketing e Anúncios</option>
                  <option value="Impostos">Impostos e Taxas</option>
                  <option value="Outros">Outras Despesas</option>
                </>
              )}
            </Select>
          </Field>

          <Field label="Descrição">
            <Input
              onChange={(e) => setFormDescription(e.target.value)}
              placeholder="Ex: Compra de luvas e descartáveis"
              required
              value={formDescription}
            />
          </Field>

          <div className="grid grid-cols-2 gap-3">
            <Field label="Valor (R$)">
              <Input
                onChange={(e) => setFormAmount(e.target.value)}
                placeholder="0,00"
                required
                type="number"
                step="0.01"
                min="0.01"
                value={formAmount}
              />
            </Field>

            <Field label="Data de Vencimento">
              <Input
                onChange={(e) => setFormDueDate(e.target.value)}
                required
                type="date"
                value={formDueDate}
              />
            </Field>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Field label="Status">
              <Select
                onChange={(e) => setFormStatus(e.target.value as TransactionStatus)}
                value={formStatus}
              >
                <option value="paid">Pago / Liquidado</option>
                <option value="pending">Pendente</option>
              </Select>
            </Field>

            <Field label="Forma de Pagamento">
              <Select
                onChange={(e) => setFormPaymentMethod(e.target.value)}
                value={formPaymentMethod}
              >
                <option value="pix">Pix</option>
                <option value="card">Cartão de Crédito / Débito</option>
                <option value="cash">Dinheiro em Espécie</option>
                <option value="transfer">Transferência / Boleto</option>
              </Select>
            </Field>
          </div>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button disabled={saving} type="submit">
              {saving ? 'Salvando...' : 'Salvar Transação'}
            </Button>
          </div>
        </Form>
      </Modal>

      {/* Modal: Abrir Caixa */}
      <Modal
        onClose={() => setIsOpenModalOpen(false)}
        open={isOpenModalOpen}
        title="Abertura de Caixa Diário"
        subtitle="Informe o fundo de troco inicial em dinheiro presente na gaveta."
      >
        <Form className="mt-4 grid gap-4" onSubmit={handleOpenShift}>
          <Field label="Fundo de Troco Inicial (R$)">
            <Input
              type="number"
              step="0.01"
              min="0"
              required
              placeholder="100.00"
              value={initialFloatInput}
              onChange={(e) => setInitialFloatInput(e.target.value)}
            />
          </Field>

          <Field label="Observações de Abertura (Opcional)">
            <Input
              placeholder="Ex.: Troco conferido pela manhã..."
              value={openNotes}
              onChange={(e) => setOpenNotes(e.target.value)}
            />
          </Field>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsOpenModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button
              className="bg-emerald-600 text-white hover:bg-emerald-700"
              disabled={submittingOpen}
              type="submit"
            >
              {submittingOpen ? 'Abrindo...' : 'Confirmar Abertura'}
            </Button>
          </div>
        </Form>
      </Modal>

      {/* Modal: Sangria / Suprimento */}
      <Modal
        onClose={() => setIsEntryModalOpen(false)}
        open={isEntryModalOpen}
        title={entryType === 'sangria' ? 'Realizar Sangria (Retirada)' : 'Realizar Suprimento (Aporte)'}
        subtitle={
          entryType === 'sangria'
            ? 'Retire dinheiro da gaveta para cofre, pagamentos avulsos ou despesas rápidas.'
            : 'Adicione troco ou aporte de dinheiro em espécie na gaveta.'
        }
      >
        <Form className="mt-4 grid gap-4" onSubmit={handleCreateEntry}>
          <Field label="Valor (R$)">
            <Input
              type="number"
              step="0.01"
              min="0.01"
              required
              placeholder="0,00"
              value={entryAmountInput}
              onChange={(e) => setEntryAmountInput(e.target.value)}
            />
          </Field>

          <Field label="Motivo / Justificativa">
            <Input
              required
              placeholder={
                entryType === 'sangria'
                  ? 'Ex: Sangria para cofre de segurança / Pagamento motoboy'
                  : 'Ex: Aporte de moedas e cédulas para troco'
              }
              value={entryReason}
              onChange={(e) => setEntryReason(e.target.value)}
            />
          </Field>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsEntryModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button
              className={
                entryType === 'sangria'
                  ? 'bg-rose-600 text-white hover:bg-rose-700'
                  : 'bg-blue-600 text-white hover:bg-blue-700'
              }
              disabled={submittingEntry}
              type="submit"
            >
              {submittingEntry ? 'Registrando...' : 'Confirmar Movimentação'}
            </Button>
          </div>
        </Form>
      </Modal>

      {/* Modal: Fechar Caixa */}
      <Modal
        onClose={() => setIsCloseModalOpen(false)}
        open={isCloseModalOpen}
        title="Fechamento de Caixa Diário"
        subtitle="Conte o dinheiro físico presente na gaveta e confirme o encerramento do turno."
      >
        <Form className="mt-4 grid gap-4" onSubmit={handleCloseShift}>
          {shiftSummary && (
            <div className="rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-3.5">
              <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
                Saldo Esperado pelo Sistema na Gaveta:
              </span>
              <p className="mt-0.5 text-xl font-black text-[var(--color-ink)]">
                {formatMoney(shiftSummary.currentDrawerCents)}
              </p>
            </div>
          )}

          <Field label="Dinheiro Contado na Gaveta (R$)">
            <Input
              type="number"
              step="0.01"
              min="0"
              required
              placeholder="0.00"
              value={closingCountInput}
              onChange={(e) => setClosingCountInput(e.target.value)}
            />
          </Field>

          {/* Live Difference feedback */}
          {shiftSummary && closingCountInput && (
            (() => {
              const countCents = Math.round(parseFloat(closingCountInput.replace(',', '.')) * 100) || 0
              const diffCents = countCents - shiftSummary.currentDrawerCents
              return (
                <div
                  className={`rounded-xl p-3 text-sm font-bold border ${
                    diffCents === 0
                      ? 'border-emerald-200 bg-emerald-50 text-emerald-800'
                      : diffCents > 0
                      ? 'border-blue-200 bg-blue-50 text-blue-800'
                      : 'border-rose-200 bg-rose-50 text-rose-800'
                  }`}
                >
                  {diffCents === 0 && '✓ Caixa confere perfeitamente com o sistema!'}
                  {diffCents > 0 && `↑ Sobra de caixa: +${formatMoney(diffCents)}`}
                  {diffCents < 0 && `⚠️ Quebra / Falta de caixa: -${formatMoney(Math.abs(diffCents))}`}
                </div>
              )
            })()
          )}

          <Field label="Observações de Fechamento (Opcional)">
            <Input
              placeholder="Ex.: Justificativa para sobra/falta ou observação final..."
              value={closingNotes}
              onChange={(e) => setClosingNotes(e.target.value)}
            />
          </Field>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsCloseModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button
              className="bg-slate-900 text-white hover:bg-black"
              disabled={submittingClose}
              type="submit"
            >
              {submittingClose ? 'Encerrando...' : 'Confirmar Fechamento de Caixa'}
            </Button>
          </div>
        </Form>
      </Modal>
    </PageContainer>
  )
}
