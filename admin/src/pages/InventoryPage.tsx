import { AlertTriangle, ArrowDown, ArrowUp, Boxes, History, Layers, Plus, RefreshCw, Trash2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

import {
  createProduct,
  createStockMovement,
  deleteProduct,
  fetchProducts,
  fetchServices,
  fetchServiceSupplies,
  fetchStockMovements,
  setServiceSupplies,
  updateProduct,
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
import type { Product, Service, ServiceSupply, StockMovement, StockMovementType } from '../types'

type ActiveTab = 'products' | 'movements' | 'supplies'

export function InventoryPage() {
  const { showToast } = useToast()
  const [activeTab, setActiveTab] = useState<ActiveTab>('products')
  const [products, setProducts] = useState<Product[]>([])
  const [movements, setMovements] = useState<StockMovement[]>([])
  const [services, setServices] = useState<Service[]>([])
  const [loading, setLoading] = useState(true)

  // Modals
  const [isProductModalOpen, setIsProductModalOpen] = useState(false)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [isMovementModalOpen, setIsMovementModalOpen] = useState(false)
  const [isSupplyModalOpen, setIsSupplyModalOpen] = useState(false)

  // Product form
  const [pName, setPName] = useState('')
  const [pSku, setPSku] = useState('')
  const [pCategory, setPCategory] = useState('Geral')
  const [pUnit, setPUnit] = useState('un')
  const [pCostPrice, setPCostPrice] = useState('0,00')
  const [pSalePrice, setPSalePrice] = useState('0,00')
  const [pInitialStock, setPInitialStock] = useState('0')
  const [pMinAlert, setPMinAlert] = useState('5')
  const [pIsSupply, setPIsSupply] = useState(false)

  // Movement form
  const [mProductId, setMProductId] = useState('')
  const [mType, setMType] = useState<StockMovementType>('in')
  const [mQuantity, setMQuantity] = useState('1')
  const [mReason, setMReason] = useState('')

  // Service supplies
  const [selectedServiceId, setSelectedServiceId] = useState<string>('')
  const [currentSupplies, setCurrentSupplies] = useState<ServiceSupply[]>([])
  const [addSupplyProductId, setAddSupplyProductId] = useState('')
  const [addSupplyQuantity, setAddSupplyQuantity] = useState('1')

  const [saving, setSaving] = useState(false)

  function loadData() {
    setLoading(true)
    Promise.all([fetchProducts(false), fetchStockMovements(undefined, 100), fetchServices()])
      .then(([prodList, moveList, svcList]) => {
        setProducts(prodList)
        setMovements(moveList)
        setServices(svcList)
        if (svcList.length > 0 && !selectedServiceId) {
          setSelectedServiceId(svcList[0].id)
        }
      })
      .catch(() => {
        showToast({ message: 'Não foi possível carregar os dados de estoque.', tone: 'error' })
      })
      .finally(() => setLoading(false))
  }

  useEffect(loadData, [])

  // Load service supplies when selected service changes
  useEffect(() => {
    if (selectedServiceId) {
      fetchServiceSupplies(selectedServiceId)
        .then(setCurrentSupplies)
        .catch(() => {})
    }
  }, [selectedServiceId])

  function resetProductForm(prod?: Product) {
    if (prod) {
      setEditingProduct(prod)
      setPName(prod.name)
      setPSku(prod.sku)
      setPCategory(prod.category)
      setPUnit(prod.unit)
      setPCostPrice((prod.costPriceCents / 100).toFixed(2))
      setPSalePrice((prod.salePriceCents / 100).toFixed(2))
      setPInitialStock(String(prod.stockQuantity))
      setPMinAlert(String(prod.minStockAlert))
      setPIsSupply(prod.isSupply)
    } else {
      setEditingProduct(null)
      setPName('')
      setPSku('')
      setPCategory('Geral')
      setPUnit('un')
      setPCostPrice('0.00')
      setPSalePrice('0.00')
      setPInitialStock('0')
      setPMinAlert('5')
      setPIsSupply(false)
    }
  }

  async function handleSaveProduct(e: FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      const costCents = Math.round(parseFloat(pCostPrice.replace(',', '.')) * 100) || 0
      const saleCents = Math.round(parseFloat(pSalePrice.replace(',', '.')) * 100) || 0
      const stockQty = parseFloat(pInitialStock.replace(',', '.')) || 0
      const minAlert = parseFloat(pMinAlert.replace(',', '.')) || 0

      if (editingProduct) {
        await updateProduct(editingProduct.id, {
          name: pName.trim(),
          sku: pSku.trim(),
          category: pCategory.trim(),
          unit: pUnit.trim(),
          costPriceCents: costCents,
          salePriceCents: saleCents,
          stockQuantity: stockQty,
          minStockAlert: minAlert,
          isSupply: pIsSupply,
        })
        showToast({ message: 'Produto atualizado.', tone: 'success' })
      } else {
        await createProduct({
          name: pName.trim(),
          sku: pSku.trim(),
          category: pCategory.trim(),
          unit: pUnit.trim(),
          costPriceCents: costCents,
          salePriceCents: saleCents,
          stockQuantity: stockQty,
          minStockAlert: minAlert,
          isSupply: pIsSupply,
          active: true,
        })
        showToast({ message: 'Produto cadastrado.', tone: 'success' })
      }
      setIsProductModalOpen(false)
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSaving(false)
    }
  }

  async function handleDeleteProduct(id: string) {
    if (!confirm('Deseja excluir este item do catálogo de produtos?')) return
    try {
      await deleteProduct(id)
      showToast({ message: 'Produto excluído.', tone: 'success' })
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    }
  }

  async function handleSaveMovement(e: FormEvent) {
    e.preventDefault()
    if (!mProductId) {
      showToast({ message: 'Selecione um produto.', tone: 'error' })
      return
    }
    const qty = parseFloat(mQuantity.replace(',', '.'))
    if (isNaN(qty) || qty <= 0) {
      showToast({ message: 'Informe uma quantidade válida.', tone: 'error' })
      return
    }

    setSaving(true)
    try {
      await createStockMovement({
        productId: mProductId,
        movementType: mType,
        quantity: qty,
        reason: mReason.trim() || 'Ajuste manual',
        referenceType: 'manual',
      })
      showToast({ message: 'Movimentação registrada com sucesso!', tone: 'success' })
      setIsMovementModalOpen(false)
      setMReason('')
      setMQuantity('1')
      loadData()
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSaving(false)
    }
  }

  async function handleAddSupply(e: FormEvent) {
    e.preventDefault()
    if (!selectedServiceId || !addSupplyProductId) return
    const qty = parseFloat(addSupplyQuantity.replace(',', '.'))
    if (isNaN(qty) || qty <= 0) return

    setSaving(true)
    try {
      const updated = [
        ...currentSupplies.filter((s) => s.productId !== addSupplyProductId),
        {
          id: '',
          tenantId: '',
          serviceId: selectedServiceId,
          productId: addSupplyProductId,
          quantity: qty,
        },
      ]
      await setServiceSupplies(
        selectedServiceId,
        updated.map((u) => ({ productId: u.productId, quantity: u.quantity })),
      )
      showToast({ message: 'Insumo vinculado ao serviço.', tone: 'success' })
      const fresh = await fetchServiceSupplies(selectedServiceId)
      setCurrentSupplies(fresh)
      setAddSupplyProductId('')
      setAddSupplyQuantity('1')
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    } finally {
      setSaving(false)
    }
  }

  async function handleRemoveSupply(productId: string) {
    if (!selectedServiceId) return
    try {
      const updated = currentSupplies.filter((s) => s.productId !== productId)
      await setServiceSupplies(
        selectedServiceId,
        updated.map((u) => ({ productId: u.productId, quantity: u.quantity })),
      )
      showToast({ message: 'Insumo desvinculado.', tone: 'success' })
      const fresh = await fetchServiceSupplies(selectedServiceId)
      setCurrentSupplies(fresh)
    } catch (err) {
      showToast({ message: (err as Error).message, tone: 'error' })
    }
  }

  const lowStockCount = products.filter((p) => p.stockQuantity <= p.minStockAlert).length

  return (
    <PageContainer>
      <Section
        actions={
          <div className="flex items-center gap-2">
            <Button
              onClick={() => {
                resetProductForm()
                setIsProductModalOpen(true)
              }}
            >
              <Plus className="mr-2 h-4 w-4" /> Novo Produto
            </Button>
            <Button
              onClick={() => {
                if (products.length > 0) setMProductId(products[0].id)
                setIsMovementModalOpen(true)
              }}
              variant="secondary"
            >
              <RefreshCw className="mr-2 h-4 w-4" /> Movimentar
            </Button>
          </div>
        }
        subtitle="Gerencie estoque de revenda, insumos para serviços e alertas de reposição"
        title="Estoque & Insumos"
      >
        {/* Navigation Tabs */}
        <div className="flex border-b border-border">
          <button
            onClick={() => setActiveTab('products')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-medium transition-colors ${
              activeTab === 'products'
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <Boxes className="h-4 w-4" />
            Catálogo de Produtos ({products.length})
            {lowStockCount > 0 && (
              <span className="rounded-full bg-rose-500/10 px-2 py-0.5 text-xs font-semibold text-rose-600">
                {lowStockCount} baixo
              </span>
            )}
          </button>

          <button
            onClick={() => setActiveTab('supplies')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-medium transition-colors ${
              activeTab === 'supplies'
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <Layers className="h-4 w-4" />
            Insumos por Serviço
          </button>

          <button
            onClick={() => setActiveTab('movements')}
            className={`flex items-center gap-2 border-b-2 px-4 py-2.5 text-sm font-medium transition-colors ${
              activeTab === 'movements'
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <History className="h-4 w-4" />
            Histórico ({movements.length})
          </button>
        </div>
      </Section>

      {/* Tab: Products */}
      {activeTab === 'products' && (
        <Section title="Itens em Estoque">
          {loading ? (
            <Spinner />
          ) : products.length === 0 ? (
            <EmptyState
              description="Nenhum produto cadastrado no estoque ainda."
              title="Estoque vazio"
            />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-border bg-muted/50 text-xs font-medium uppercase text-muted-foreground">
                  <tr>
                    <th className="px-4 py-3">Produto</th>
                    <th className="px-4 py-3">SKU</th>
                    <th className="px-4 py-3">Categoria</th>
                    <th className="px-4 py-3">Tipo</th>
                    <th className="px-4 py-3">Preço Custo</th>
                    <th className="px-4 py-3">Preço Venda</th>
                    <th className="px-4 py-3">Estoque</th>
                    <th className="px-4 py-3 text-right">Ações</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {products.map((p) => {
                    const isLow = p.stockQuantity <= p.minStockAlert
                    return (
                      <tr key={p.id} className="hover:bg-muted/30 transition-colors">
                        <td className="px-4 py-3 font-medium text-foreground">{p.name}</td>
                        <td className="px-4 py-3 font-mono text-xs text-muted-foreground">{p.sku || '—'}</td>
                        <td className="px-4 py-3 text-muted-foreground">{p.category}</td>
                        <td className="px-4 py-3">
                          <span
                            className={`rounded-md px-2 py-0.5 text-xs font-medium ${
                              p.isSupply
                                ? 'bg-amber-500/10 text-amber-700 dark:text-amber-400'
                                : 'bg-primary/10 text-primary'
                            }`}
                          >
                            {p.isSupply ? 'Insumo (Consumo)' : 'Venda / Revenda'}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-muted-foreground">{formatMoney(p.costPriceCents)}</td>
                        <td className="px-4 py-3 font-semibold text-foreground">{formatMoney(p.salePriceCents)}</td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            <span className={`font-bold ${isLow ? 'text-rose-600' : 'text-foreground'}`}>
                              {p.stockQuantity} {p.unit}
                            </span>
                            {isLow && (
                              <span
                                className="inline-flex items-center gap-1 rounded bg-rose-500/10 px-1.5 py-0.5 text-xs font-medium text-rose-600"
                                title={`Mínimo: ${p.minStockAlert}`}
                              >
                                <AlertTriangle className="h-3 w-3" /> Baixo
                              </span>
                            )}
                          </div>
                        </td>
                        <td className="px-4 py-3 text-right whitespace-nowrap">
                          <div className="flex items-center justify-end gap-1">
                            <Button
                              onClick={() => {
                                resetProductForm(p)
                                setIsProductModalOpen(true)
                              }}
                              size="sm"
                              variant="ghost"
                            >
                              Editar
                            </Button>
                            <Button
                              onClick={() => handleDeleteProduct(p.id)}
                              size="icon"
                              variant="ghost"
                            >
                              <Trash2 className="h-4 w-4 text-rose-600" />
                            </Button>
                          </div>
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

      {/* Tab: Service Supplies Mapping */}
      {activeTab === 'supplies' && (
        <Section
          subtitle="Defina quais insumos são baixados automaticamente do estoque ao concluir um agendamento deste serviço"
          title="Consumo Automático de Insumos"
        >
          <div className="space-y-6">
            <div className="max-w-md">
              <Field label="Selecione o Serviço">
                <Select
                  onChange={(e) => setSelectedServiceId(e.target.value)}
                  value={selectedServiceId}
                >
                  {services.map((svc) => (
                    <option key={svc.id} value={svc.id}>
                      {svc.name} ({svc.durationMinutes} min)
                    </option>
                  ))}
                </Select>
              </Field>
            </div>

            {/* Form to add supply to selected service */}
            <form onSubmit={handleAddSupply} className="rounded-xl border border-border bg-card p-4">
              <h4 className="text-sm font-semibold text-foreground mb-3">Vincular Novo Insumo ao Serviço</h4>
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                <Field label="Insumo">
                  <Select
                    onChange={(e) => setAddSupplyProductId(e.target.value)}
                    required
                    value={addSupplyProductId}
                  >
                    <option value="">Selecione o produto/insumo...</option>
                    {products.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.name} ({p.unit})
                      </option>
                    ))}
                  </Select>
                </Field>

                <Field label="Quantidade Consumida por Serviço">
                  <Input
                    onChange={(e) => setAddSupplyQuantity(e.target.value)}
                    placeholder="1"
                    required
                    step="0.01"
                    type="number"
                    value={addSupplyQuantity}
                  />
                </Field>

                <div className="flex items-end">
                  <Button disabled={saving || !addSupplyProductId} type="submit">
                    Vincular Insumo
                  </Button>
                </div>
              </div>
            </form>

            {/* List of current supplies for service */}
            <div className="rounded-xl border border-border bg-card overflow-hidden">
              <div className="border-b border-border bg-muted/40 px-4 py-3 font-semibold text-sm">
                Insumos configurados para este serviço
              </div>
              {currentSupplies.length === 0 ? (
                <div className="p-6 text-center text-sm text-muted-foreground">
                  Nenhum insumo configurado. Ao concluir um agendamento deste serviço, nenhum item de estoque será baixado.
                </div>
              ) : (
                <table className="w-full text-left text-sm">
                  <thead className="border-b border-border bg-muted/50 text-xs uppercase text-muted-foreground">
                    <tr>
                      <th className="px-4 py-2.5">Insumo</th>
                      <th className="px-4 py-2.5">Quantidade por Agendamento</th>
                      <th className="px-4 py-2.5 text-right">Ação</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {currentSupplies.map((supply) => {
                      const prod = products.find((p) => p.id === supply.productId)
                      return (
                        <tr key={supply.productId}>
                          <td className="px-4 py-2.5 font-medium text-foreground">
                            {prod ? prod.name : supply.productId}
                          </td>
                          <td className="px-4 py-2.5 text-muted-foreground">
                            {supply.quantity} {prod?.unit || 'un'}
                          </td>
                          <td className="px-4 py-2.5 text-right">
                            <Button
                              onClick={() => handleRemoveSupply(supply.productId)}
                              size="sm"
                              variant="ghost"
                            >
                              <Trash2 className="h-4 w-4 text-rose-600" />
                            </Button>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              )}
            </div>
          </div>
        </Section>
      )}

      {/* Tab: Movements */}
      {activeTab === 'movements' && (
        <Section title="Histórico de Entradas e Baixas">
          {movements.length === 0 ? (
            <EmptyState
              description="Nenhuma movimentação de estoque registrada até o momento."
              title="Sem movimentações"
            />
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-border bg-muted/50 text-xs font-medium uppercase text-muted-foreground">
                  <tr>
                    <th className="px-4 py-3">Data</th>
                    <th className="px-4 py-3">Produto</th>
                    <th className="px-4 py-3">Tipo</th>
                    <th className="px-4 py-3">Quantidade</th>
                    <th className="px-4 py-3">Estoque Anterior → Novo</th>
                    <th className="px-4 py-3">Motivo</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {movements.map((m) => {
                    const prod = products.find((p) => p.id === m.productId)
                    return (
                      <tr key={m.id} className="hover:bg-muted/30 transition-colors">
                        <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">
                          {m.createdAt ? new Date(m.createdAt).toLocaleString('pt-BR') : '—'}
                        </td>
                        <td className="px-4 py-3 font-medium text-foreground">
                          {prod ? prod.name : m.productName || m.productId}
                        </td>
                        <td className="px-4 py-3">
                          <span
                            className={`inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-semibold ${
                              m.movementType === 'in'
                                ? 'bg-emerald-500/10 text-emerald-600'
                                : m.movementType === 'out'
                                ? 'bg-rose-500/10 text-rose-600'
                                : 'bg-primary/10 text-primary'
                            }`}
                          >
                            {m.movementType === 'in' ? (
                              <ArrowDown className="h-3 w-3" />
                            ) : (
                              <ArrowUp className="h-3 w-3" />
                            )}
                            {m.movementType === 'in'
                              ? 'Entrada'
                              : m.movementType === 'out'
                              ? 'Baixa / Saída'
                              : 'Ajuste'}
                          </span>
                        </td>
                        <td className="px-4 py-3 font-semibold text-foreground">
                          {m.movementType === 'out' ? '-' : '+'}
                          {m.quantity}
                        </td>
                        <td className="px-4 py-3 text-muted-foreground">
                          {m.previousStock} → <span className="font-bold text-foreground">{m.newStock}</span>
                        </td>
                        <td className="px-4 py-3 text-muted-foreground">{m.reason}</td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          )}
        </Section>
      )}

      {/* Modal: Create/Edit Product */}
      <Modal
        subtitle="Configure os dados do produto ou insumo de trabalho."
        open={isProductModalOpen}
        onClose={() => setIsProductModalOpen(false)}
        title={editingProduct ? 'Editar Produto' : 'Novo Produto ou Insumo'}
      >
        <Form onSubmit={handleSaveProduct}>
          <Field label="Nome do Item">
            <Input
              onChange={(e) => setPName(e.target.value)}
              placeholder="Ex: Óleo Essencial Lavanda 100ml"
              required
              value={pName}
            />
          </Field>

          <div className="grid grid-cols-2 gap-3">
            <Field label="Código / SKU">
              <Input
                onChange={(e) => setPSku(e.target.value)}
                placeholder="Ex: OEL-001"
                value={pSku}
              />
            </Field>

            <Field label="Categoria">
              <Input
                onChange={(e) => setPCategory(e.target.value)}
                placeholder="Ex: Cosméticos, Descartáveis..."
                value={pCategory}
              />
            </Field>
          </div>

          <div className="grid grid-cols-3 gap-3">
            <Field label="Unidade">
              <Select onChange={(e) => setPUnit(e.target.value)} value={pUnit}>
                <option value="un">Unidade (un)</option>
                <option value="ml">Mililitros (ml)</option>
                <option value="g">Gramas (g)</option>
                <option value="kg">Quilos (kg)</option>
                <option value="cx">Caixa (cx)</option>
              </Select>
            </Field>

            <Field label="Preço Custo (R$)">
              <Input
                onChange={(e) => setPCostPrice(e.target.value)}
                placeholder="0.00"
                step="0.01"
                type="number"
                value={pCostPrice}
              />
            </Field>

            <Field label="Preço Venda (R$)">
              <Input
                onChange={(e) => setPSalePrice(e.target.value)}
                placeholder="0.00"
                step="0.01"
                type="number"
                value={pSalePrice}
              />
            </Field>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Field label="Qtd. em Estoque">
              <Input
                onChange={(e) => setPInitialStock(e.target.value)}
                placeholder="0"
                step="0.01"
                type="number"
                value={pInitialStock}
              />
            </Field>

            <Field label="Alerta de Estoque Mínimo">
              <Input
                onChange={(e) => setPMinAlert(e.target.value)}
                placeholder="5"
                step="0.01"
                type="number"
                value={pMinAlert}
              />
            </Field>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              id="isSupplyCheckbox"
              type="checkbox"
              checked={pIsSupply}
              onChange={(e) => setPIsSupply(e.target.checked)}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
            <label htmlFor="isSupplyCheckbox" className="text-sm font-medium text-foreground cursor-pointer">
              É insumo interno de serviço (não visível para venda avulsa direta)
            </label>
          </div>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsProductModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button disabled={saving} type="submit">
              {saving ? 'Salvando...' : 'Salvar Item'}
            </Button>
          </div>
        </Form>
      </Modal>

      {/* Modal: Movement */}
      <Modal
        subtitle="Lançar entrada de estoque, avaria ou ajuste de balanço."
        open={isMovementModalOpen}
        onClose={() => setIsMovementModalOpen(false)}
        title="Movimentação de Estoque"
      >
        <Form onSubmit={handleSaveMovement}>
          <Field label="Produto">
            <Select onChange={(e) => setMProductId(e.target.value)} required value={mProductId}>
              {products.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name} (Atual: {p.stockQuantity} {p.unit})
                </option>
              ))}
            </Select>
          </Field>

          <div className="grid grid-cols-2 gap-3">
            <Field label="Tipo de Movimento">
              <Select onChange={(e) => setMType(e.target.value as StockMovementType)} value={mType}>
                <option value="in">Entrada (+)</option>
                <option value="out">Baixa / Perda (-)</option>
                <option value="adjustment">Balanço / Ajuste Direto</option>
              </Select>
            </Field>

            <Field label="Quantidade">
              <Input
                onChange={(e) => setMQuantity(e.target.value)}
                placeholder="1"
                required
                step="0.01"
                type="number"
                value={mQuantity}
              />
            </Field>
          </div>

          <Field label="Motivo / Observação">
            <Input
              onChange={(e) => setMReason(e.target.value)}
              placeholder="Ex: Compra NF 1234, Vencimento de lote, contagem física..."
              required
              value={mReason}
            />
          </Field>

          <div className="mt-4 flex justify-end gap-2">
            <Button onClick={() => setIsMovementModalOpen(false)} type="button" variant="secondary">
              Cancelar
            </Button>
            <Button disabled={saving} type="submit">
              {saving ? 'Registrando...' : 'Confirmar Movimentação'}
            </Button>
          </div>
        </Form>
      </Modal>
    </PageContainer>
  )
}
