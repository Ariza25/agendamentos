const WIDGET_RESIZE_MESSAGE = 'widget:resize'
const WIDGET_REFRESH_MESSAGE = 'widget:refresh'

export function initWidgetBridge() {
  if (typeof window === 'undefined' || window.self === window.top) {
    return
  }

  let frame: number | null = null

  const sendHeight = () => {
    if (frame !== null) {
      cancelAnimationFrame(frame)
    }
    frame = requestAnimationFrame(() => {
      const height = Math.ceil(document.documentElement.scrollHeight)
      if (height > 0) {
        window.parent.postMessage({ type: WIDGET_RESIZE_MESSAGE, height }, '*')
      }
    })
  }

  const observer = new ResizeObserver(sendHeight)
  observer.observe(document.documentElement)
  window.addEventListener('resize', sendHeight)
  window.addEventListener('message', (event) => {
    if (event.source !== window.parent) {
      return
    }
    const data = event.data
    if (data && data.type === WIDGET_REFRESH_MESSAGE) {
      sendHeight()
    }
  })

  sendHeight()
}
