import { Check, DollarSign, Percent, Plus, Settings2, Trash2, Users } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createCommissionRule,
  deleteCommissionRule,
  fetchCommissionRecords,
  fetchCommissionRules,
  fetchCommissionSummaries,
  fetchProfessionals,
  fetchServices,
  settleCommissions,
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
import type {
  CommissionRecord,
  CommissionRule,
  CommissionSummary,
  CommissionType,
  Professional,
  Service,
} from '../types'

export function CommissionsPage() {
  const { showToast } = useToast()
  const [summaries, setSummaries] = useState<CommissionSummary[]>([])
  const [rules, setRules] = useState<CommissionRule[]>([])
  const [records, setRecords] = useState<CommissionRecord[]>([])
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

  // Filter
  const [selectedProfFilter, setSelectedProfFilter] = useState('')

  function loadData() {
    setLoading(true)
    Promise.all([
      fetchCommissionSummaries(),
      fetchCommissionRules(),
      fetchCommissionRecords(selectedProfFilter || undefined),
      fetchProfessionals(),
      fetchServices(),
    ])
      .then(([sumList, ruleList, recList, profList, svcList]) => {
        setSummaries(sumList)
        setRules(ruleList)
        setRecords(recList)
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
      showToast({ message: 'Regra de comissão criada!', tone: 'success' })
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

  async function handleSettleProfessional(summary: CommissionSummary) {
    if (!confirm(`Deseja liquidar todas as comissões pendentes (${formatMoney(summary.pendingCents)}) de ${summary.professionalName}?`)) {
      return
    }
    try {
      await settleCommissions(summary.professionalId, [])
      showToast({ message: `Comissões de ${summary.professionalName} liquidadas com sucesso!`, tone: 'success' })
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    }
  }

  return (
    <PageContainer>
      <Section
        actions={
          <Button onClick={() => setIsRuleModalOpen(true)}>
            <Plus className="mr-2 h-4 w-4" /> Nova Regra de Comissão
          </Button>
        }
        subtitle="Regras de repasse, extrato de comissionamento e liquidação para a equipe"
        title="Comissões"
      >
        {/* Summaries by Professional */}
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {summaries.map((s) => (
            <div key={s.professionalId} className="rounded-xl border border-border bg-card p-4 shadow-sm flex flex-col justify-between">
              <div>
                <div className="flex items-center justify-between">
                  <span className="font-semibold text-foreground flex items-center gap-2">
                    <Users className="h-4 w-4 text-primary" /> {s.professionalName}
                  </span>
                  <Badge tone={s.pendingCents > 0 ? 'warning' : 'neutral'}>
                    {s.count} atendimento{s.count !== 1 ? 's' : ''}
                  </Badge>
                </div>
                <div className="mt-4 space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Pendente a pagar:</span>
                    <span className="font-bold text-amber-600 dark:text-amber-400">{formatMoney(s.pendingCents)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Já liquidado:</span>
                    <span className="font-medium text-emerald-600 dark:text-emerald-400">{formatMoney(s.settledCents)}</span>
                  </div>
                  <div className="flex justify-between border-t border-border pt-1">
                    <span className="font-medium text-foreground">Total gerado:</span>
                    <span className="font-bold text-foreground">{formatMoney(s.totalCents)}</span>
                  </div>
                </div>
              </div>

              <div className="mt-4 pt-3 border-t border-border">
                <Button
                  disabled={s.pendingCents <= 0}
                  onClick={() => handleSettleProfessional(s)}
                  size="sm"
                  variant={s.pendingCents > 0 ? 'primary' : 'secondary'}
                  className="w-full"
                >
                  <Check className="mr-1.5 h-3.5 w-3.5" />
                  {s.pendingCents > 0 ? 'Liquidar Comissões' : 'Totalmente Liquidado'}
                </Button>
              </div>
            </div>
          ))}

          {summaries.length === 0 && !loading && (
            <div className="col-span-full py-6 text-center text-sm text-muted-foreground">
              Nenhuma comissão acumulada até o momento. Conclua agendamentos para gerar lançamentos automáticos.
            </div>
          )}
        </div>
      </Section>

      {/* Commission Rules */}
      <Section
        subtitle="As comissões são calculadas automaticamente ao marcar o status do agendamento como 'Concluído'"
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

      {/* Commission Records History */}
      <Section
        actions={
          <Select
            className="w-48 text-sm"
            onChange={(e) => setSelectedProfFilter(e.target.value)}
            value={selectedProfFilter}
          >
            <option value="">Filtrar: Todos profissionais</option>
            {professionals.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </Select>
        }
        title="Histórico de Comissões"
      >
        {records.length === 0 ? (
          <EmptyState
            description="Nenhum registro de comissão encontrado."
            title="Sem lançamentos"
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="border-b border-border bg-muted/50 text-xs font-medium uppercase text-muted-foreground">
                <tr>
                  <th className="px-4 py-3">Data</th>
                  <th className="px-4 py-3">Profissional</th>
                  <th className="px-4 py-3">Valor da Comissão</th>
                  <th className="px-4 py-3">Status</th>
                  <th className="px-4 py-3">Data de Liquidação</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {records.map((rec) => {
                  const prof = professionals.find((p) => p.id === rec.professionalId)
                  return (
                    <tr key={rec.id} className="hover:bg-muted/30 transition-colors">
                      <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">
                        {rec.createdAt ? new Date(rec.createdAt).toLocaleDateString('pt-BR') : '—'}
                      </td>
                      <td className="px-4 py-3 font-medium text-foreground">
                        {prof ? prof.name : rec.professionalId}
                      </td>
                      <td className="px-4 py-3 font-bold text-foreground">
                        {formatMoney(rec.amountCents)}
                      </td>
                      <td className="px-4 py-3">
                        <Badge tone={rec.status === 'settled' ? 'success' : 'warning'}>
                          {rec.status === 'settled' ? 'Liquidado' : 'Pendente'}
                        </Badge>
                      </td>
                      <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">
                        {rec.settledAt ? new Date(rec.settledAt).toLocaleDateString('pt-BR') : '—'}
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}
      </Section>

      {/* Modal: New Commission Rule */}
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
