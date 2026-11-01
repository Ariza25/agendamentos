import { useEffect, useRef, useState } from 'react'

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

    function connect() {
      try {
        source = new EventSource('/api/v1/admin/events', { withCredentials: true })

        source.onopen = () => {
          setIsConnected(true)
        }

        source.onmessage = (e) => {
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
          setIsConnected(false)
          if (source) {
            source.close()
            source = null
          }
          // Attempt reconnect after 5s
          retryTimeout = setTimeout(connect, 5000)
        }
      } catch {
        setIsConnected(false)
        retryTimeout = setTimeout(connect, 5000)
      }
    }

    connect()

    return () => {
      if (retryTimeout) clearTimeout(retryTimeout)
      if (source) {
        source.close()
        source = null
      }
      setIsConnected(false)
    }
  }, [])

  return { isConnected }
}

