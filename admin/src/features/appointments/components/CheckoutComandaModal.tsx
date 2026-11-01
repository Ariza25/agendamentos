import {
  Banknote,
  CheckCircle,
  CreditCard,
  DollarSign,
  PackagePlus,
  Percent,
  Receipt,
  Trash2,
  User,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import { checkoutAppointment, fetchProducts } from '../../../lib/adminApi'
import { formatMoney } from '../../../lib/format'
import { Button } from '../../../shared/components/Button'
import { Field, Form, Input, Select } from '../../../shared/components/FormControls'
import { Modal } from '../../../shared/components/Modal'
import { useToast } from '../../../shared/components/toast/ToastProvider'
import type { Appointment, Customer, Product, Professional, Service } from '../../../types'

type CheckoutComandaModalProps = {
  open: boolean
  onClose: () => void
  appointment: Appointment | null
  customer?: Customer
  service?: Service
  professional?: Professional
  onSuccess: () => void
}

type SelectedProduct = {
  product: Product
  quantity: number
}

export function CheckoutComandaModal({
  open,
  onClose,
  appointment,
  customer,
  service,
  professional,
  onSuccess,
}: CheckoutComandaModalProps) {
  const { showToast } = useToast()

  const [availableProducts, setAvailableProducts] = useState<Product[]>([])
  const [loadingProducts, setLoadingProducts] = useState(false)

  // Selected items
  const [selectedProducts, setSelectedProducts] = useState<SelectedProduct[]>([])
  const [addingProductId, setAddingProductId] = useState('')
  const [addingQuantity, setAddingQuantity] = useState(1)

  // Pricing adjustments
  const [discountInput, setDiscountInput] = useState('')
  const [surchargeInput, setSurchargeInput] = useState('')
  const [paymentMethod, setPaymentMethod] = useState<'pix' | 'cash' | 'credit_card' | 'debit_card' | 'other'>('pix')
  const [notes, setNotes] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (open) {
      setLoadingProducts(true)
      fetchProducts(true)
        .then((items) => {
          // Filter products with stock > 0 and not marked purely as supplies
          setAvailableProducts(items.filter((p) => p.active))
        })
        .catch(() => {
          showToast({ message: 'Não foi possível carregar os produtos.', tone: 'error' })
        })
        .finally(() => {
          setLoadingProducts(false)
        })

      // Reset state
      setSelectedProducts([])
      setAddingProductId('')
      setAddingQuantity(1)
      setDiscountInput('')
      setSurchargeInput('')
      setPaymentMethod('pix')
      setNotes('')
    }
  }, [open])

  if (!open || !appointment) {
    return null
  }

  const servicePriceCents = service?.priceCents || 0
  const productsTotalCents = selectedProducts.reduce(
    (acc, item) => acc + item.product.salePriceCents * item.quantity,
    0
  )

  const discountCents = Math.round(parseFloat(discountInput.replace(',', '.') || '0') * 100) || 0
  const surchargeCents = Math.round(parseFloat(surchargeInput.replace(',', '.') || '0') * 100) || 0

  const totalChargedCents = Math.max(0, servicePriceCents + productsTotalCents - discountCents + surchargeCents)

  const commissionPct = professional?.commissionPercentage ?? 0
  const estimatedCommissionCents = Math.round((servicePriceCents * commissionPct) / 100)

  function handleAddProduct() {
    if (!addingProductId) return
    const found = availableProducts.find((p) => p.id === addingProductId)
    if (!found) return

    if (addingQuantity <= 0) {
      showToast({ message: 'Quantidade deve ser maior que zero.', tone: 'error' })
      return
    }

    const existingIndex = selectedProducts.findIndex((item) => item.product.id === found.id)
    if (existingIndex >= 0) {
      setSelectedProducts((prev) =>
        prev.map((item, idx) =>
          idx === existingIndex ? { ...item, quantity: item.quantity + addingQuantity } : item
        )
      )
    } else {
      setSelectedProducts((prev) => [...prev, { product: found, quantity: addingQuantity }])
    }

    setAddingProductId('')
    setAddingQuantity(1)
  }

  function handleRemoveProduct(productId: string) {
    setSelectedProducts((prev) => prev.filter((item) => item.product.id !== productId))
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    if (!appointment) return
    setIsSubmitting(true)

    try {
      const res = await checkoutAppointment({
        appointmentId: appointment.id,
        products: selectedProducts.map((p) => ({
          productId: p.product.id,
          quantity: p.quantity,
        })),
        discountCents,
        surchargeCents,
        paymentMethod,
        notes: notes.trim(),
      })

      showToast({
        message: `Comanda fechada com sucesso! Total recebido: ${formatMoney(res.totalChargedCents)}`,
        tone: 'success',
      })
      onSuccess()
      onClose()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      open={open}
      onClose={onClose}
      size="lg"
      title="Fechar Comanda / Frente de Caixa"
      subtitle="Finalize o atendimento, adicione produtos do estoque, aplique descontos e liquide a comanda."
    >
      <Form className="mt-4 grid gap-5" onSubmit={handleSubmit}>
        {/* Atendimento Header Card */}
        <div className="rounded-2xl border border-[var(--color-line)] bg-[var(--color-bg-muted)] p-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <span className="text-xs font-bold uppercase tracking-wider text-[var(--color-ink-soft)]">
                Cliente & Atendimento
              </span>
              <h4 className="text-base font-bold text-[var(--color-ink)]">
                {customer?.name || 'Cliente'}
                {customer?.phone ? ` · ${customer.phone}` : ''}
              </h4>
              <p className="mt-0.5 text-xs text-[var(--color-ink-soft)]">
                {service?.name || 'Serviço'}
                {service?.durationMinutes ? ` (${service.durationMinutes} min)` : ''}
              </p>
            </div>

            <div className="flex items-center gap-4 border-t border-[var(--color-line)] pt-3 sm:border-0 sm:pt-0">
              <div className="text-right">
                <span className="text-xs font-semibold text-[var(--color-ink-soft)]">Profissional</span>
                <p className="text-sm font-bold text-[var(--color-ink)]">{professional?.name || 'Não vinculado'}</p>
                {commissionPct > 0 ? (
                  <span className="inline-flex items-center gap-1 text-[11px] font-semibold text-emerald-600">
                    <Percent className="h-3 w-3" /> {commissionPct}% comissão ({formatMoney(estimatedCommissionCents)})
                  </span>
                ) : (
                  <span className="text-[11px] text-[var(--color-ink-soft)]">Sem comissão cadastrada</span>
                )}
              </div>
              <div className="rounded-xl bg-[var(--color-primary)]/10 px-3 py-2 text-right">
                <span className="block text-[10px] font-bold uppercase tracking-wider text-[var(--color-primary)]">
                  Serviço
                </span>
                <span className="text-base font-black text-[var(--color-primary)]">
                  {formatMoney(servicePriceCents)}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Adicionar Produtos do Estoque */}
        <div className="rounded-2xl border border-[var(--color-line)] p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <PackagePlus className="h-4 w-4 text-[var(--color-primary)]" />
              <h4 className="text-sm font-bold text-[var(--color-ink)]">
                Venda de Produtos Adicionais (Baixa no Estoque)
              </h4>
            </div>
            <span className="text-xs font-semibold text-[var(--color-ink-soft)]">
              {selectedProducts.length} produto(s) na comanda
            </span>
          </div>

          <div className="mt-3 grid gap-2 sm:grid-cols-[1fr_auto_auto]">
            <Select
              className="text-xs"
              value={addingProductId}
              onChange={(e) => setAddingProductId(e.target.value)}
              disabled={loadingProducts}
            >
              <option value="">Selecione um produto para adicionar...</option>
              {availableProducts.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name} — {formatMoney(p.salePriceCents)} (Estoque: {p.stockQuantity} {p.unit})
                </option>
              ))}
            </Select>

            <Input
              type="number"
              min={1}
              className="w-24 text-xs"
              value={addingQuantity}
              onChange={(e) => setAddingQuantity(Math.max(1, parseInt(e.target.value) || 1))}
              placeholder="Qtd"
            />

            <Button
              type="button"
              variant="secondary"
              size="sm"
              onClick={handleAddProduct}
              disabled={!addingProductId}
            >
              Adicionar
            </Button>
          </div>

          {selectedProducts.length > 0 && (
            <ul className="mt-3 divide-y divide-[var(--color-line)] rounded-xl border border-[var(--color-line)] bg-[var(--color-bg-muted)]">
              {selectedProducts.map(({ product, quantity }) => (
                <li
                  key={product.id}
                  className="flex items-center justify-between px-3 py-2 text-xs font-medium"
                >
                  <div className="flex items-center gap-2">
                    <span className="font-bold text-[var(--color-ink)]">{product.name}</span>
                    <span className="text-[var(--color-ink-soft)]">
                      ({quantity}x {formatMoney(product.salePriceCents)})
                    </span>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="font-bold text-[var(--color-ink)]">
                      {formatMoney(product.salePriceCents * quantity)}
                    </span>
                    <button
                      type="button"
                      onClick={() => handleRemoveProduct(product.id)}
                      className="text-rose-500 hover:text-rose-700"
                      aria-label="Remover produto"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Totais, Desconto e Acréscimo */}
        <div className="grid gap-3 sm:grid-cols-2">
          <Field label="Desconto (R$)">
            <Input
              type="number"
              step="0.01"
              min="0"
              placeholder="0,00"
              value={discountInput}
              onChange={(e) => setDiscountInput(e.target.value)}
            />
          </Field>

          <Field label="Acréscimo / Taxa (R$)">
            <Input
              type="number"
              step="0.01"
              min="0"
              placeholder="0,00"
              value={surchargeInput}
              onChange={(e) => setSurchargeInput(e.target.value)}
            />
          </Field>
        </div>

        {/* Forma de Pagamento */}
        <div>
          <label className="text-xs font-bold uppercase tracking-wider text-[var(--color-ink-soft)]">
            Forma de Pagamento
          </label>
          <div className="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-5">
            {[
              { id: 'pix', label: 'Pix', icon: Receipt },
              { id: 'cash', label: 'Dinheiro', icon: Banknote },
              { id: 'credit_card', label: 'Crédito', icon: CreditCard },
              { id: 'debit_card', label: 'Débito', icon: CreditCard },
              { id: 'other', label: 'Outros', icon: DollarSign },
            ].map((method) => {
              const Icon = method.icon
              const isSelected = paymentMethod === method.id
              return (
                <button
                  key={method.id}
                  type="button"
                  onClick={() => setPaymentMethod(method.id as any)}
                  className={`flex flex-col items-center justify-center gap-1.5 rounded-xl border p-2.5 text-xs font-bold transition-colors ${
                    isSelected
                      ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/10 text-[var(--color-primary)]'
                      : 'border-[var(--color-line)] bg-transparent text-[var(--color-ink-soft)] hover:bg-[var(--color-bg-muted)]'
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {method.label}
                </button>
              )
            })}
          </div>
          {paymentMethod === 'cash' && (
            <p className="mt-1.5 text-xs font-medium text-amber-700">
              * O pagamento em dinheiro será somado ao saldo em espécie na gaveta do caixa diário.
            </p>
          )}
        </div>

        {/* Observações */}
        <Field label="Observações da Comanda">
          <Input
            placeholder="Ex.: Cliente solicitou nota ou observação adicional..."
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
          />
        </Field>

        {/* Resumo Financeiro Final */}
        <div className="flex items-center justify-between rounded-2xl bg-[var(--color-ink)] p-4 text-white">
          <div>
            <span className="block text-xs font-semibold text-white/70">Total a Pagar</span>
            <div className="flex items-baseline gap-2">
              <span className="text-2xl font-black">{formatMoney(totalChargedCents)}</span>
              {discountCents > 0 && (
                <span className="text-xs text-rose-300">(-{formatMoney(discountCents)})</span>
              )}
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button type="button" variant="secondary" onClick={onClose}>
              Cancelar
            </Button>
            <Button
              type="submit"
              loading={isSubmitting}
              className="bg-emerald-600 text-white hover:bg-emerald-700"
              icon={<CheckCircle className="h-4 w-4" />}
            >
              Concluir & Receber
            </Button>
          </div>
        </div>
      </Form>
    </Modal>
  )
}
