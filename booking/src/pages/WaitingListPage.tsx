import { ArrowLeft, BellRing } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import { fetchPublicProfessionals, fetchPublicServices, signUpWaitingList } from '../lib/bookingApi'
import type { Professional, Service } from '../lib/bookingApi'
import { formatMoney } from '../lib/format'
import { Button } from '../shared/components/Button'
import { Field, Form, Input, Select } from '../shared/components/FormControls'
import { Spinner } from '../shared/components/Spinner'

export function WaitingListPage() {
  const { slug = '' } = useParams<{ slug: string }>()
  const navigate = useNavigate()

  const [services, setServices] = useState<Service[]>([])
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [serviceId, setServiceId] = useState('')
  const [professionalId, setProfessionalId] = useState('')
  const [name, setName] = useState('')
  const [phone, setPhone] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isDone, setIsDone] = useState(false)

  useEffect(() => {
    let active = true
    setLoading(true)
    Promise.all([fetchPublicServices(slug), fetchPublicProfessionals(slug)])
      .then(([servicesData, professionalsData]) => {
        if (!active) {
          return
        }
        setServices(servicesData)
        setProfessionals(professionalsData)
        if (servicesData[0]) {
          setServiceId(servicesData[0].id)
        }
      })
      .catch((loadError: unknown) => {
        if (active) {
          setError((loadError as Error).message)
        }
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })
    return () => {
      active = false
    }
  }, [slug])

  const selectedService = useMemo(() => services.find((item) => item.id === serviceId) ?? null, [services, serviceId])

  const eligibleProfessionals = useMemo(() => {
    if (!selectedService || !selectedService.professionalIds?.length) {
      return professionals
    }
    return professionals.filter((item) => selectedService.professionalIds?.includes(item.id))
  }, [professionals, selectedService])

  function handleServiceChange(value: string) {
    setServiceId(value)
    const service = services.find((item) => item.id === value)
    if (!service?.professionalIds?.length) {
      setProfessionalId('')
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!serviceId) {
      setError('Escolha o serviço para entrar na lista de espera.')
      return
    }
    if (!name.trim()) {
      setError('Informe seu nome.')
      return
    }
    if (!phone.trim()) {
      setError('Informe seu telefone para contato.')
      return
    }

    setIsSubmitting(true)
    setError('')
    try {
      await signUpWaitingList(slug, {
        serviceId,
        professionalId: professionalId || undefined,
        name: name.trim(),
        phone: phone.trim(),
      })
      setIsDone(true)
    } catch (submitError) {
      setError((submitError as Error).message)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loading) {
    return (
      <main className="flex min-h-dvh items-center justify-center px-4 py-10">
        <Spinner label="Carregando lista de espera..." />
      </main>
    )
  }

  return (
    <main className="flex min-h-dvh items-center justify-center px-4 py-10">
      <section className="w-full max-w-md rounded-3xl border border-[var(--color-line)] bg-white p-8 shadow-sm">
        {isDone ? (
          <div className="text-center">
            <span className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-[var(--color-success)] text-white">
              <BellRing className="h-7 w-7" />
            </span>
            <h1 className="mt-4 font-headline text-2xl font-extrabold text-[var(--color-ink)]">Você está na lista!</h1>
            <p className="mt-2 text-sm font-semibold text-[var(--color-ink-soft)]">
              Assim que um horário ficar livre, entraremos em contato pelo telefone informado.
            </p>
            <Button className="mt-6" onClick={() => navigate('/')} size="lg" variant="secondary">
              Voltar ao início
            </Button>
          </div>
        ) : error && services.length === 0 ? (
          <div className="text-center">
            <h1 className="font-headline text-xl font-extrabold text-[var(--color-ink)]">Não foi possível carregar</h1>
            <p className="mt-2 text-sm font-semibold text-[var(--color-ink-soft)]">{error}</p>
            <Button className="mt-6" onClick={() => navigate('/')} size="lg" variant="secondary">
              Voltar ao início
            </Button>
          </div>
        ) : (
          <>
            <div className="flex items-center gap-3">
              <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--color-primary)] text-white">
                <BellRing className="h-6 w-6" />
              </span>
              <div>
                <h1 className="font-headline text-2xl font-extrabold text-[var(--color-ink)]">Lista de espera</h1>
                <p className="text-sm font-semibold text-[var(--color-ink-soft)]">
                  Avisamos quando surgir um horário livre.
                </p>
              </div>
            </div>

            {error ? (
              <p className="mt-4 rounded-xl border border-red-200 bg-red-50 p-3 text-sm font-bold text-red-700">{error}</p>
            ) : null}

            <Form className="mt-6" onSubmit={handleSubmit}>
              <Field label="Serviço">
                <Select className="mt-2" onChange={(event) => handleServiceChange(event.target.value)} value={serviceId}>
                  {services.map((service) => (
                    <option key={service.id} value={service.id}>
                      {service.name} · {formatMoney(service.priceCents)}
                    </option>
                  ))}
                </Select>
              </Field>
              {eligibleProfessionals.length > 0 ? (
                <Field hint="Opcional" label="Profissional preferido">
                  <Select className="mt-2" onChange={(event) => setProfessionalId(event.target.value)} value={professionalId}>
                    <option value="">Qualquer profissional</option>
                    {eligibleProfessionals.map((professional) => (
                      <option key={professional.id} value={professional.id}>
                        {professional.name}
                      </option>
                    ))}
                  </Select>
                </Field>
              ) : null}
              <Field label="Nome completo">
                <Input
                  className="mt-2"
                  onChange={(event) => setName(event.target.value)}
                  placeholder="Como podemos te chamar?"
                  value={name}
                />
              </Field>
              <Field label="Telefone / WhatsApp">
                <Input
                  className="mt-2"
                  onChange={(event) => setPhone(event.target.value)}
                  placeholder="(00) 00000-0000"
                  value={phone}
                />
              </Field>
              <Button icon={<BellRing className="h-4 w-4" />} loading={isSubmitting} size="lg" type="submit">
                Entrar na lista de espera
              </Button>
              <Button
                className="mt-2 w-full"
                icon={<ArrowLeft className="h-4 w-4" />}
                onClick={() => navigate('/')}
                variant="ghost"
              >
                Voltar ao início
              </Button>
            </Form>
          </>
        )}
      </section>
    </main>
  )
}
