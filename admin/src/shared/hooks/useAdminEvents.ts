import { useEffect, useRef, useState } from 'react'
import { buildApiUrl } from '../../lib/apiBase'
import { getAdminAuthToken } from '../../lib/adminAuth'

export type AdminEvent = {
  type: string
  tenantId: string
  data: any
  timestamp: string
}

type EventListener = (event: AdminEvent) => void

export function useAdminEvents(onEvent?: EventListener) {
  const [isConnected, setIsConnected] = useState(false)
  const onEventRef = useRef(onEvent)
  onEventRef.current = onEvent

  useEffect(() => {
    let source: EventSource | null = null
    let retryTimeout: ReturnType<typeof setTimeout> | null = null
    let isDisposed = false

    function cleanup() {
      if (retryTimeout) {
        clearTimeout(retryTimeout)
        retryTimeout = null
      }
      if (source) {
        source.close()
        source = null
      }
      setIsConnected(false)
    }

    function scheduleReconnect() {
      if (isDisposed) return
      cleanup()
      const token = getAdminAuthToken()
      if (token) {
        retryTimeout = setTimeout(connect, 10000)
      }
    }

    function connect() {
      if (isDisposed) return
      const token = getAdminAuthToken()
      if (!token) {
        cleanup()
        return
      }

      cleanup()

      try {
        const url = buildApiUrl(`/api/v1/admin/events?token=${encodeURIComponent(token)}`)
        source = new EventSource(url, { withCredentials: true })

        source.onopen = () => {
          if (!isDisposed) {
            setIsConnected(true)
          }
        }

        source.onmessage = (e) => {
          if (isDisposed) return
          try {
            const parsed: AdminEvent = JSON.parse(e.data)
            if (onEventRef.current) {
              onEventRef.current(parsed)
            }
          } catch {
            // Heartbeat ping or non-json message, ignore
          }
        }

        source.onerror = () => {
          scheduleReconnect()
        }
      } catch {
        scheduleReconnect()
      }
    }

    connect()

    const handleAuthChange = () => {
      connect()
    }

    window.addEventListener('erpagendamentos:admin:auth:changed', handleAuthChange)

    return () => {
      isDisposed = true
      window.removeEventListener('erpagendamentos:admin:auth:changed', handleAuthChange)
      cleanup()
    }
  }, [])

  return { isConnected }
}

