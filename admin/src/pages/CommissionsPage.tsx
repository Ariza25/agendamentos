import {
  Calendar,
  Check,
  CheckCircle,
  Clock,
  DollarSign,
  FileText,
  History,
  Layers,
  Percent,
  Plus,
  Printer,
  Receipt,
  Settings2,
  Trash2,
  Users,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createCommissionPayout,
  createCommissionRule,
  deleteCommissionRule,
  fetchCommissionPayouts,
  fetchCommissionRecords,
  fetchCommissionRules,
  fetchCommissionSummaries,
  fetchProfessionals,
  fetchServices,
} from '../lib/adminApi'
import { formatDateTime, formatDate, formatMoney } from '../lib/format'
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
  CommissionPayout,
  CommissionRecord,
  CommissionRule,
  CommissionSummary,
  CommissionType,
  Professional,
  Service,
} from '../types'

type ActiveTab = 'summaries' | 'payouts' | 'rules'

export function CommissionsPage() {
  const { showToast } = useToast()
  const [activeTab, setActiveTab] = useState<ActiveTab>('summaries')

  const [summaries, setSummaries] = useState<CommissionSummary[]>([])
  const [rules, setRules] = useState<CommissionRule[]>([])
  const [records, setRecords] = useState<CommissionRecord[]>([])
  const [payouts, setPayouts] = useState<CommissionPayout[]>([])
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [services, setServices] = useState<Service[]>([])
  const [loading, setLoading] = useState(true)

  // Rule Modal state
  const [isRuleModalOpen, setIsRuleModalOpen] = useState(false)
  const [ruleProfId, setRuleProfId] = useState('')
  const [ruleServiceId, setRuleServiceId] = useState('')
  const [ruleType, setRuleType] = useState<CommissionType>('percentage')
  const [ruleValue, setRuleValue] = useState('30')
  const [savingRule, setSavingRule] = useState(false)

  // Payout Modal state
  const [payoutModalProf, setPayoutModalProf] = useState<CommissionSummary | null>(null)
  const [viewingReceiptPayout, setViewingReceiptPayout] = useState<CommissionPayout | null>(null)

  // Filter
  const [selectedProfFilter, setSelectedProfFilter] = useState('')

  function loadData() {
    setLoading(true)
    Promise.all([
      fetchCommissionSummaries(),
      fetchCommissionRules(),
      fetchCommissionRecords(selectedProfFilter || undefined),
      fetchCommissionPayouts(selectedProfFilter || undefined),
      fetchProfessionals(),
      fetchServices(),
    ])
      .then(([sumList, ruleList, recList, payoutList, profList, svcList]) => {
        setSummaries(sumList)
        setRules(ruleList)
        setRecords(recList)
        setPayouts(payoutList)
        setProfessionals(profList)
        setServices(svcList)
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar as comissões.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadData, [selectedProfFilter])

  async function handleCreateRule(e: FormEvent) {
    e.preventDefault()
    const val = parseFloat(ruleValue.replace(',', '.'))
    if (isNaN(val) || val <= 0) {
      showToast({ message: 'Informe um valor válido maior que zero.', tone: 'error' })
      return
    }

    setSavingRule(true)
    try {
      await createCommissionRule({
        professionalId: ruleProfId || undefined,
        serviceId: ruleServiceId || undefined,
        commissionType: ruleType,
        value: val,
        active: true,
      })
      showToast({ message: 'Regra de comissão criada com sucesso!', tone: 'success' })
      setIsRuleModalOpen(false)
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSavingRule(false)
    }
  }

  async function handleDeleteRule(id: string) {
    if (!confirm('Deseja excluir esta regra de comissão?')) return
    try {
      await deleteCommissionRule(id)
      showToast({ message: 'Regra excluída.', tone: 'success' })
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <div className="flex items-center gap-2">
            <Button
              icon={<Plus className="h-4 w-4" />}
              onClick={() => setIsRuleModalOpen(true)}
              variant="secondary"
            >
              Nova Regra
            </Button>
          </div>
        }
        subtitle="Controle de comissionamento da equipe, fechamentos quinzenais/mensais e lançamento de despesas no DRE."
        title="Comissões & Repasses da Equipe"
      >
        {/* Navigation Tabs */}
        <div className="mt-4 flex border-b border-border">
          <button
            onClick={() => setActiveTab('summaries')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'summaries'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <Users className="h-4 w-4" />
            Equipe & Repasses Pendentes
          </button>
          <button
            onClick={() => setActiveTab('payouts')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'payouts'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <History className="h-4 w-4" />
            Histórico de Fechamentos Pagos ({payouts.length})
          </button>
          <button
            onClick={() => setActiveTab('rules')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-semibold transition-colors ${
              activeTab === 'rules'
                ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <Settings2 className="h-4 w-4" />
            Regras de Comissão ({rules.length})
          </button>
        </div>
      </Section>

      {activeTab === 'summaries' && (
        <>
          {/* Summaries by Professional */}
          <Section subtitle="Valores de comissão acumulados aguardando fechamento." title="Saldos por Profissional">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {summaries.map((s) => (
                <div
                  key={s.professionalId}
                  className="rounded-2xl border border-border bg-card p-5 shadow-sm flex flex-col justify-between"
                >
                  <div>
                    <div className="flex items-center justify-between">
                      <span className="font-bold text-foreground flex items-center gap-2 text-base">
                        <Users className="h-4 w-4 text-[var(--color-primary)]" /> {s.professionalName}
                      </span>
                      <Badge tone={s.pendingCents > 0 ? 'warning' : 'success'}>
                        {s.pendingCents > 0 ? 'Pendente' : 'Quitado'}
                      </Badge>
                    </div>

                    <div className="mt-4 space-y-2 text-sm">
                      <div className="flex justify-between items-baseline">
                        <span className="text-xs font-semibold text-muted-foreground">Pendente a pagar:</span>
                        <span className="text-xl font-black text-amber-600 dark:text-amber-400">
                          {formatMoney(s.pendingCents)}
                        </span>
                      </div>
                      <div className="flex justify-between text-xs">
                        <span className="text-muted-foreground">Já liquidado:</span>
                        <span className="font-medium text-emerald-600 dark:text-emerald-400">
                          {formatMoney(s.settledCents || (s.totalCents - s.pendingCents))}
                        </span>
                      </div>
                      <div className="flex justify-between border-t border-border pt-1.5 text-xs">
                        <span className="font-medium text-foreground">Total gerado histórico:</span>
                        <span className="font-bold text-foreground">{formatMoney(s.totalCents)}</span>
                      </div>
                    </div>
                  </div>

                  <div className="mt-5 pt-3 border-t border-border">
                    <Button
                      disabled={s.pendingCents <= 0}
                      onClick={() => setPayoutModalProf(s)}
                      size="sm"
                      className={`w-full ${
                        s.pendingCents > 0
                          ? 'bg-emerald-600 text-white hover:bg-emerald-700'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      <Receipt className="mr-1.5 h-3.5 w-3.5" />
                      {s.pendingCents > 0 ? 'Fechar Repasse & Pagar' : 'Sem Comissões Pendentes'}
                    </Button>
                  </div>
                </div>
              ))}

              {summaries.length === 0 && !loading && (
                <div className="col-span-full py-8 text-center text-sm text-muted-foreground">
                  Nenhuma comissão acumulada até o momento. Conclua agendamentos ou feche comandas para gerar comissões automáticas.
                </div>
              )}
            </div>
          </Section>

          {/* Commission Records Table */}
          <Section
            actions={
              <Select
                className="w-48 text-xs font-medium"
                onChange={(e) => setSelectedProfFilter(e.target.value)}
                value={selectedProfFilter}
              >
                <option value="">Todos os profissionais</option>
                {professionals.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </Select>
            }
            subtitle="Extrato detalhado de cada serviço e atendimento que gerou comissão."
            title="Extrato Individual de Lançamentos"
          >
            {records.length === 0 ? (
              <EmptyState
                description="Nenhum registro de comissão encontrado para o filtro selecionado."
                title="Sem lançamentos"
              />
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-border text-xs font-semibold uppercase text-muted-foreground">
                    <tr>
                      <th className="px-4 py-3">Data</th>
                      <th className="px-4 py-3">Profissional</th>
                      <th className="px-4 py-3">Serviço</th>
                      <th className="px-4 py-3 text-right">Valor do Serviço</th>
                      <th className="px-4 py-3 text-right">Comissão</th>
                      <th className="px-4 py-3">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {records.map((rec) => {
                      const amount = rec.commissionAmountCents ?? rec.amountCents ?? 0
                      return (
                        <tr key={rec.id} className="hover:bg-muted/30 transition-colors">
                          <td className="px-4 py-3 text-xs text-muted-foreground whitespace-nowrap">
                            {rec.createdAt ? new Date(rec.createdAt).toLocaleDateString('pt-BR') : '—'}
                          </td>
                          <td className="px-4 py-3 font-semibold text-foreground">
                            {rec.professionalName || rec.professionalId}
                          </td>
                          <td className="px-4 py-3 text-xs text-muted-foreground">
                            {rec.serviceName || 'Serviço'}
                          </td>
                          <td className="px-4 py-3 text-xs text-right font-medium text-muted-foreground">
                            {formatMoney(rec.baseAmountCents || 0)}
                          </td>
                          <td className="px-4 py-3 text-right font-bold text-foreground">
                            {formatMoney(amount)}
                          </td>
                          <td className="px-4 py-3">
                            <Badge tone={rec.status === 'paid' || rec.status === 'settled' ? 'success' : 'warning'}>
                              {rec.status === 'paid' || rec.status === 'settled' ? 'Pago' : 'Pendente'}
                            </Badge>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </Section>
        </>
      )}

      {activeTab === 'payouts' && (
        <Section
          subtitle="Auditoria de todos os lotes de fechamentos quitados com comprovante e despesa financeira registrada."
          title="Histórico de Fechamentos Realizados"
        >
          {payouts.length === 0 ? (
            <EmptyState
              description="Nenhum fechamento de comissão foi realizado ainda. Realize fechamentos quinzenais ou mensais na aba 'Equipe'."
              title="Sem histórico de fechamentos"
            />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-border text-xs uppercase text-muted-foreground">
                  <tr>
                    <th className="px-4 py-3">Data do Fechamento</th>
                    <th className="px-4 py-3">Profissional</th>
                    <th className="px-4 py-3">Atendimentos</th>
                    <th className="px-4 py-3">Forma de Pagamento</th>
                    <th className="px-4 py-3 text-right">Total Pago</th>
                    <th className="px-4 py-3 text-right">Comprovante</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {payouts.map((p) => (
                    <tr key={p.id} className="hover:bg-muted/30 transition-colors">
                      <td className="px-4 py-3 font-medium text-foreground text-xs">
                        {formatDateTime(p.paidAt || p.createdAt)}
                      </td>
                      <td className="px-4 py-3 font-bold text-foreground">
                        {p.professionalName || 'Profissional'}
                      </td>
                      <td className="px-4 py-3 text-xs text-muted-foreground">
                        {p.recordsCount} lançamento(s)
                      </td>
                      <td className="px-4 py-3 uppercase text-xs font-semibold text-muted-foreground">
                        {p.paymentMethod}
                      </td>
                      <td className="px-4 py-3 text-right font-black text-emerald-600">
                        {formatMoney(p.totalCents)}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <Button
                          size="sm"
                          variant="secondary"
                          icon={<FileText className="h-3.5 w-3.5" />}
                          onClick={() => setViewingReceiptPayout(p)}
                        >
                          Ver Recibo
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Section>
      )}

      {activeTab === 'rules' && (
        <Section
          subtitle="Configure percentuais de repasse customizados por profissional ou por tipo de serviço."
          title="Regras Ativas de Repasse"
        >
          {rules.length === 0 ? (
            <EmptyState
              description="Cadastre regras para automatizar o cálculo de comissões por porcentagem (%) ou valor fixo (R$)."
              title="Nenhuma regra cadastrada"
            />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-border bg-muted/50 text-xs font-medium uppercase text-muted-foreground">
                  <tr>
                    <th className="px-4 py-3">Profissional</th>
                    <th className="px-4 py-3">Serviço Aplicável</th>
                    <th className="px-4 py-3">Tipo de Repasse</th>
                    <th className="px-4 py-3">Valor / Alíquota</th>
                    <th className="px-4 py-3 text-right">Ação</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {rules.map((r) => {
                    const prof = professionals.find((p) => p.id === r.professionalId)
                    const svc = services.find((s) => s.id === r.serviceId)
                    return (
                      <tr key={r.id} className="hover:bg-muted/30 transition-colors">
                        <td className="px-4 py-3 font-medium text-foreground">
                          {prof ? prof.name : <span className="text-muted-foreground italic">Todos os profissionais</span>}
                        </td>
                        <td className="px-4 py-3 text-muted-foreground">
                          {svc ? svc.name : <span className="italic">Todos os serviços</span>}
                        </td>
                        <td className="px-4 py-3">
                          <span className="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-0.5 text-xs font-semibold text-foreground">
                            {r.commissionType === 'percentage' ? <Percent className="h-3 w-3" /> : <DollarSign className="h-3 w-3" />}
                            {r.commissionType === 'percentage' ? 'Percentual (%)' : 'Valor Fixo (R$)'}
                          </span>
                        </td>
                        <td className="px-4 py-3 font-bold text-foreground">
                          {r.commissionType === 'percentage' ? `${r.value}%` : formatMoney(Math.round(r.value * 100))}
                        </td>
                        <td className="px-4 py-3 text-right">
                          <Button onClick={() => handleDeleteRule(r.id)} size="icon" variant="ghost">
                            <Trash2 className="h-4 w-4 text-rose-600" />
                          </Button>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          )}
        </Section>
      )}

      {/* Modal: Fechamento de Comissões (Payout) */}
      {payoutModalProf && (
        <CommissionPayoutModal
          professional={payoutModalProf}
          onClose={() => setPayoutModalProf(null)}
          onSuccess={(payout) => {
            setPayoutModalProf(null)
            setViewingReceiptPayout(payout)
            loadData()
          }}
        />
      )}

      {/* Modal: Recibo / Comprovante de Pagamento */}
      {viewingReceiptPayout && (
        <CommissionReceiptModal
          payout={viewingReceiptPayout}
          onClose={() => setViewingReceiptPayout(null)}
        />
      )}

      {/* Modal: Nova Regra de Comissão */}
      <Modal
        subtitle="Defina como o cálculo do repasse será efetuado para a equipe."
        open={isRuleModalOpen}
        onClose={() => setIsRuleModalOpen(false)}
        title="Nova Regra de Comissão"
      >
        <Form onSubmit={handleCreateRule}>
          <Field label="Profissional">
            <Select onChange={(e) => setRuleProfId(e.target.value)} value={ruleProfId}>
              <option value="">Todos os profissionais (Regra Geral)</option>
              {professionals.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </Select>
          </Field>

          <Field label="Serviço">
            <Select onChange={(e) => setRuleServiceId(e.target.value)} value={ruleServiceId}>
              <option value="">Todos os serviços</option>
              {services.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name}
                </option>
              ))}
            </Select>
          </Field>

          <div className="grid grid-cols-2 gap-3">
            <Field label="Tipo de Cálculo">
              <Select onChange={(e) => setRuleType(e.target.value as CommissionType)} value={ruleType}>
                <option value="percentage">Porcentagem (%)</option>
                <option value="fixed">Valor Fixo (R$)</option>
              </Select>
            </Field>

            <Field label={ruleType === 'percentage' ? 'Percentual (%)' : 'Valor Fixo (R$)'}>
              <Input
                onChange={(e) => setRuleValue(e.target.value)}
                placeholder={ruleType === 'percentage' ? 'Ex: 30' : 'Ex: 45.00'}
                required
                step="0.01"
                type="number"
                value={ruleValue}
              />
            </Field>
          </div>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsRuleModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button disabled={savingRule} type="submit">
              {savingRule ? 'Salvando...' : 'Salvar Regra'}
            </Button>
          </div>
        </Form>
      </Modal>
    </PageContainer>
  )
}

function CommissionPayoutModal({
  professional,
  onClose,
  onSuccess,
}: {
  professional: CommissionSummary
  onClose: () => void
  onSuccess: (payout: CommissionPayout) => void
}) {
  const { showToast } = useToast()
  const [paymentMethod, setPaymentMethod] = useState('pix')
  const [notes, setNotes] = useState(`Acerto de comissões - ${new Date().toLocaleDateString('pt-BR')}`)
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleConfirm(e: FormEvent) {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      const payout = await createCommissionPayout({
        professionalId: professional.professionalId,
        paymentMethod,
        notes: notes.trim(),
      })
      showToast({
        message: `Fechamento realizado com sucesso! Total pago: ${formatMoney(payout.totalCents)}`,
        tone: 'success',
      })
      onSuccess(payout)
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      open
      onClose={onClose}
      size="md"
      title={`Fechar Repasse: ${professional.professionalName}`}
      subtitle="Confirme a liquidação das comissões acumuladas e lance automaticamente a despesa no Financeiro."
    >
      <Form className="mt-4 grid gap-4" onSubmit={handleConfirm}>
        <div className="rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
              Total Acumulado a Pagar
            </span>
            <span className="text-2xl font-black text-emerald-600">
              {formatMoney(professional.pendingCents)}
            </span>
          </div>
          <p className="mt-2 text-xs text-[var(--color-ink-soft)]">
            * Ao confirmar, todos os lançamentos pendentes serão marcados como quitados e uma despesa com categoria <strong>Comissões</strong> será lançada no DRE financeiro.
          </p>
        </div>

        <Field label="Forma de Pagamento">
          <Select value={paymentMethod} onChange={(e) => setPaymentMethod(e.target.value)}>
            <option value="pix">Pix (Chave do Profissional)</option>
            <option value="cash">Dinheiro em Espécie (Gaveta do Caixa)</option>
            <option value="transfer">Transferência Bancária / TED</option>
            <option value="other">Outro</option>
          </Select>
        </Field>

        <Field label="Observações do Fechamento">
          <Input
            placeholder="Ex.: Acerto 1ª Quinzena de Novembro..."
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
          />
        </Field>

        <div className="mt-2 flex justify-end gap-2 border-t pt-4">
          <Button type="button" variant="secondary" onClick={onClose}>
            Cancelar
          </Button>
          <Button
            type="submit"
            loading={isSubmitting}
            className="bg-emerald-600 text-white hover:bg-emerald-700"
            icon={<CheckCircle className="h-4 w-4" />}
          >
            Confirmar Pagamento & Gerar Recibo
          </Button>
        </div>
      </Form>
    </Modal>
  )
}

function CommissionReceiptModal({
  payout,
  onClose,
}: {
  payout: CommissionPayout
  onClose: () => void
}) {
  function handlePrint() {
    window.print()
  }

  return (
    <Modal
      open
      onClose={onClose}
      size="md"
      title="Recibo de Repasse de Comissões"
      subtitle={`Comprovante de pagamento de comissões - #${payout.id.slice(0, 8)}`}
    >
      <div className="mt-4 rounded-2xl border border-[var(--color-line)] bg-white p-6 shadow-sm print:border-0 print:p-0">
        <div className="flex items-center justify-between border-b pb-4">
          <div>
            <h3 className="text-lg font-black text-slate-900">ERP Agendamentos</h3>
            <p className="text-xs text-slate-500">Comprovante de Liquidação de Comissões</p>
          </div>
          <span className="rounded-full bg-emerald-100 px-3 py-1 text-xs font-bold text-emerald-800">
            PAGO
          </span>
        </div>

        <div className="mt-4 space-y-3 text-sm">
          <div className="flex justify-between">
            <span className="text-slate-500">Profissional Beneficiário:</span>
            <span className="font-bold text-slate-900">{payout.professionalName}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-500">Data de Pagamento:</span>
            <span className="font-medium text-slate-700">
              {formatDateTime(payout.paidAt || payout.createdAt)}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-500">Atendimentos / Lançamentos:</span>
            <span className="font-medium text-slate-700">{payout.recordsCount} atendimento(s)</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-500">Forma de Pagamento:</span>
            <span className="font-semibold uppercase text-slate-800">{payout.paymentMethod}</span>
          </div>
          {payout.notes && (
            <div className="flex justify-between border-t pt-2 text-xs">
              <span className="text-slate-500">Observações:</span>
              <span className="text-slate-700 italic">{payout.notes}</span>
            </div>
          )}

          <div className="mt-4 flex items-center justify-between border-t-2 border-slate-900 pt-3">
            <span className="text-base font-bold text-slate-900">VALOR TOTAL PAGO:</span>
            <span className="text-2xl font-black text-emerald-600">
              {formatMoney(payout.totalCents)}
            </span>
          </div>
        </div>

        <div className="mt-8 border-t border-dashed pt-6 text-center text-xs text-slate-400">
          <div className="mx-auto mb-1 h-0.5 w-48 bg-slate-300" />
          <span>Assinatura do Profissional</span>
        </div>
      </div>

      <div className="mt-4 flex justify-between items-center print:hidden">
        <Button variant="secondary" icon={<Printer className="h-4 w-4" />} onClick={handlePrint}>
          Imprimir Recibo
        </Button>
        <Button onClick={onClose}>Fechar</Button>
      </div>
    </Modal>
  )
}
