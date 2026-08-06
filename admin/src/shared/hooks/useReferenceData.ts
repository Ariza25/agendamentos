import { useEffect, useMemo, useState } from 'react'

import { fetchCustomers, fetchProfessionals, fetchServices } from '../../lib/adminApi'
import type { Customer, Professional, Service } from '../../types'

type ReferenceData = {
  customers: Customer[]
  services: Service[]
  professionals: Professional[]
  byId: {
    customer: Map<string, Customer>
    service: Map<string, Service>
    professional: Map<string, Professional>
  }
  loading: boolean
  reload: () => void
}

export function useReferenceData(): ReferenceData {
  const [customers, setCustomers] = useState<Customer[]>([])
  const [services, setServices] = useState<Service[]>([])
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [loading, setLoading] = useState(true)
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    let active = true
    setLoading(true)

    Promise.allSettled([fetchProfessionals(), fetchServices(), fetchCustomers()])
      .then(([professionalsResult, servicesResult, customersResult]) => {
        if (!active) {
          return
        }
        setProfessionals(professionalsResult.status === 'fulfilled' ? professionalsResult.value : [])
        setServices(servicesResult.status === 'fulfilled' ? servicesResult.value : [])
        setCustomers(customersResult.status === 'fulfilled' ? customersResult.value : [])
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })

    return () => {
      active = false
    }
  }, [reloadKey])

  const byId = useMemo(
    () => ({
      customer: new Map(customers.map((customer) => [customer.id, customer])),
      service: new Map(services.map((service) => [service.id, service])),
      professional: new Map(professionals.map((professional) => [professional.id, professional])),
    }),
    [customers, services, professionals],
  )

  return {
    customers,
    services,
    professionals,
    byId,
    loading,
    reload: () => setReloadKey((key) => key + 1),
  }
}
