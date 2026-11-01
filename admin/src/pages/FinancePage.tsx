import { ArrowDownLeft, ArrowUpRight, DollarSign, Plus, Trash2, CheckCircle } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createFinancialTransaction,
  deleteFinancialTransaction,
  fetchFinancialSummary,
  fetchFinancialTransactions,
  updateFinancialTransaction,
} from '../lib/adminApi'
import { formatMoney } from '../lib/format'
import { Badge } from '../shared/components/Badge'
import { Button } from '../shared/components/Button'
import { EmptyState } from '../shared/components/EmptyState'
import { Field, Form, Input, Select } from '../shared/components/FormControls'
import { Modal } from '../shared/components/Modal'
import { PageContainer } from '../shared/components/PageContainer'
import { Section } from '../shared/components/Section'
import { Spinner } from '../shared/components/Spinner'
import { useToast } from '../shared/components/toast/ToastProvider'
import type { FinancialSummary, FinancialTransaction, TransactionStatus, TransactionType } from '../types'

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

export function FinancePage() {
  const { showToast } = useToast()
  const [transactions, setTransactions] = useState<FinancialTransaction[]>([])
  const [summary, setSummary] = useState<FinancialSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [filterType, setFilterType] = useState<string>('')
  const [filterStatus, setFilterStatus] = useState<string>('')

  // Form state
  const [formType, setFormType] = useState<TransactionType>('income')
  const [formCategory, setFormCategory] = useState('Serviços')
  const [formDescription, setFormDescription] = useState('')
  const [formAmount, setFormAmount] = useState('')
  const [formDueDate, setFormDueDate] = useState(() => new Date().toISOString().slice(0, 10))
  const [formStatus, setFormStatus] = useState<TransactionStatus>('paid')
  const [formPaymentMethod, setFormPaymentMethod] = useState('pix')
  const [saving, setSaving] = useState(false)

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

  useEffect(loadData, [filterType, filterStatus])

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

  return (
    <PageContainer>
      <Section
        actions={
          <Button onClick={() => { resetForm(); setIsModalOpen(true) }}>
            <Plus className="mr-2 h-4 w-4" /> Nova Transação
          </Button>
        }
        subtitle="Controle fluxo de caixa, despesas operacionais e DRE resumido"
        title="Financeiro"
      >
        {/* KPI Cards */}
        {summary && (
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
              <p className={`mt-2 text-2xl font-bold ${summary.netProfitCents >= 0 ? 'text-primary' : 'text-rose-600'}`}>
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
              <p className="mt-1 text-xs text-muted-foreground">
                Previsão de recebíveis futuros
              </p>
            </div>
          </div>
        )}
      </Section>

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
            <table className="w-full text-left text-sm">
              <thead className="border-b border-border bg-muted/50 text-xs font-medium uppercase text-muted-foreground">
                <tr>
                  <th className="px-4 py-3">Data</th>
                  <th className="px-4 py-3">Tipo</th>
                  <th className="px-4 py-3">Categoria</th>
                  <th className="px-4 py-3">Descrição</th>
                  <th className="px-4 py-3">Valor</th>
                  <th className="px-4 py-3">Status</th>
                  <th className="px-4 py-3 text-right">Ações</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {transactions.map((tx) => (
                  <tr key={tx.id} className="hover:bg-muted/30 transition-colors">
                    <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">
                      {tx.dueDate ? tx.dueDate.slice(0, 10) : '—'}
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-semibold ${
                          tx.type === 'income'
                            ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                            : 'bg-rose-500/10 text-rose-600 dark:text-rose-400'
                        }`}
                      >
                        {tx.type === 'income' ? <ArrowDownLeft className="h-3 w-3" /> : <ArrowUpRight className="h-3 w-3" />}
                        {tx.type === 'income' ? 'Receita' : 'Despesa'}
                      </span>
                    </td>
                    <td className="px-4 py-3 font-medium text-foreground">{tx.category || 'Geral'}</td>
                    <td className="px-4 py-3 text-muted-foreground">{tx.description}</td>
                    <td className="px-4 py-3 font-semibold text-foreground whitespace-nowrap">
                      <span className={tx.type === 'income' ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400'}>
                        {tx.type === 'income' ? '+' : '-'} {formatMoney(tx.amountCents)}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <Badge tone={statusTone(tx.status)}>
                        {STATUS_LABELS[tx.status] || tx.status}
                      </Badge>
                    </td>
                    <td className="px-4 py-3 text-right whitespace-nowrap">
                      <div className="flex items-center justify-end gap-1">
                        {tx.status === 'pending' && (
                          <Button
                            onClick={() => handleMarkPaid(tx)}
                            size="icon"
                            title="Marcar como Pago"
                            variant="ghost"
                          >
                            <CheckCircle className="h-4 w-4 text-emerald-600" />
                          </Button>
                        )}
                        <Button
                          onClick={() => handleDelete(tx.id)}
                          size="icon"
                          title="Excluir Transação"
                          variant="ghost"
                        >
                          <Trash2 className="h-4 w-4 text-rose-600" />
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

      {/* Modal Nova Transação */}
      <Modal
        subtitle="Preencha os detalhes para lançar uma receita ou despesa manual."
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Nova Transação Financeira"
      >
        <Form onSubmit={handleCreate}>
          <Field label="Tipo de Lançamento">
            <Select
              onChange={(e) => {
                const t = e.target.value as TransactionType
                setFormType(t)
                setFormCategory(t === 'income' ? 'Serviços' : 'Suprimentos')
              }}
              value={formType}
            >
              <option value="income">Receita (+)</option>
              <option value="expense">Despesa (-)</option>
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
              <Select onChange={(e) => setFormStatus(e.target.value as TransactionStatus)} value={formStatus}>
                <option value="paid">Pago / Liquidado</option>
                <option value="pending">Pendente</option>
              </Select>
            </Field>

            <Field label="Forma de Pagamento">
              <Select onChange={(e) => setFormPaymentMethod(e.target.value)} value={formPaymentMethod}>
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
    </PageContainer>
  )
}
