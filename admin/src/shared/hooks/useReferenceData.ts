import { useEffect, useMemo, useState } from 'react'

import { fetchCustomers, fetchProfessionals, fetchResources, fetchServices } from '../../lib/adminApi'
import type { Customer, Professional, Resource, Service } from '../../types'

type ReferenceData = {
  customers: Customer[]
  services: Service[]
  professionals: Professional[]
  resources: Resource[]
  byId: {
    customer: Map<string, Customer>
    service: Map<string, Service>
    professional: Map<string, Professional>
    resource: Map<string, Resource>
  }
  loading: boolean
  reload: () => void
}

export function useReferenceData(): ReferenceData {
  const [customers, setCustomers] = useState<Customer[]>([])
  const [services, setServices] = useState<Service[]>([])
  const [professionals, setProfessionals] = useState<Professional[]>([])
  const [resources, setResources] = useState<Resource[]>([])
  const [loading, setLoading] = useState(true)
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    let active = true
    setLoading(true)

    Promise.allSettled([fetchProfessionals(), fetchServices(), fetchResources(), fetchCustomers()])
      .then(([professionalsResult, servicesResult, resourcesResult, customersResult]) => {
        if (!active) {
          return
        }
        setProfessionals(professionalsResult.status === 'fulfilled' ? professionalsResult.value : [])
        setServices(servicesResult.status === 'fulfilled' ? servicesResult.value : [])
        setResources(resourcesResult.status === 'fulfilled' ? resourcesResult.value : [])
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
      resource: new Map(resources.map((resource) => [resource.id, resource])),
    }),
    [customers, services, professionals, resources],
  )

  return {
    customers,
    services,
    professionals,
    resources,
    byId,
    loading,
    reload: () => setReloadKey((key) => key + 1),
  }
}
