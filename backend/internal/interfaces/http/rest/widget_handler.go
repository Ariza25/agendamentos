package rest

import (
	"net/http"
	"os"
	"strings"
)


const widgetScript = `(function () {
  'use strict'

  var BOOKING_ORIGIN = '__ERP_WIDGET_BOOKING_ORIGIN__'

  var script = document.currentScript
  if (!script) {
    return
  }

  function readAttr(name, fallback) {
    var value = script.getAttribute('data-' + name)
    if (value === null || value === undefined) {
      return fallback
    }
    value = value.trim()
    return value === '' ? fallback : value
  }

  var slug = readAttr('slug', '').toLowerCase()
  var mode = readAttr('mode', 'button') === 'inline' ? 'inline' : 'button'
  var label = readAttr('label', 'Agendar horario')
  var position = readAttr('position', 'bottom-right') === 'bottom-left' ? 'bottom-left' : 'bottom-right'
  var rawHeight = parseInt(readAttr('height', '640'), 10)
  var primary = readAttr('primary', '#4f46e5')

  if (!/^[a-z0-9]([a-z0-9-]{0,60}[a-z0-9])?$/.test(slug)) {
    return
  }
  if (isNaN(rawHeight) || rawHeight < 320) {
    rawHeight = 640
  }
  if (!/^#[0-9a-fA-F]{3,8}$/.test(primary)) {
    primary = '#4f46e5'
  }

  var root = document.createElement('div')
  root.id = 'erp-widget-root-' + slug
  root.setAttribute('data-erp-widget', 'root')

  var iframe = document.createElement('iframe')
  iframe.setAttribute('frameborder', '0')
  iframe.setAttribute('allowtransparency', 'true')
  iframe.setAttribute('allow', 'camera; microphone; payment')
  iframe.setAttribute('title', 'Agendamento online')
  iframe.src = BOOKING_ORIGIN + '/#/booking/' + slug
  iframe.style.width = '100%'
  iframe.style.border = '0'
  iframe.style.display = 'block'
  iframe.style.background = 'transparent'
  iframe.style.height = rawHeight + 'px'

  function applyResize(event) {
    if (event.source !== iframe.contentWindow) {
      return
    }
    var data = event.data
    if (!data || data.type !== 'widget:resize' || typeof data.height !== 'number') {
      return
    }
    var next = Math.max(200, Math.min(Math.round(data.height), 3000))
    iframe.style.height = next + 'px'
  }

  function mount() {
    if (mode === 'inline') {
      root.appendChild(iframe)
      script.parentNode.insertBefore(root, script.nextSibling)
      return
    }

    var panel = document.createElement('div')
    panel.setAttribute('role', 'dialog')
    panel.setAttribute('aria-label', 'Agendamento online')
    panel.style.cssText = [
      'position:fixed',
      'z-index:2147483000',
      position === 'bottom-left' ? 'left:12px' : 'right:12px',
      'bottom:84px',
      'width:min(420px, calc(100vw - 24px))',
      'max-height:calc(100vh - 96px)',
      'display:none',
      'overflow:hidden',
      'border-radius:16px',
      'box-shadow:0 24px 48px -12px rgba(0,0,0,0.35)',
      'background:#ffffff',
      'border:1px solid rgba(0,0,0,0.08)'
    ].join(';')
    panel.appendChild(iframe)
    root.appendChild(panel)

    var button = document.createElement('button')
    button.type = 'button'
    button.textContent = label
    button.style.cssText = [
      'position:fixed',
      'z-index:2147483000',
      'bottom:24px',
      position === 'bottom-left' ? 'left:24px' : 'right:24px',
      'padding:14px 22px',
      'border:0',
      'border-radius:999px',
      'cursor:pointer',
      'font:600 15px/1.2 -apple-system, Segoe UI, Roboto, Helvetica, Arial, sans-serif',
      'color:#ffffff',
      'background:' + primary,
      'box-shadow:0 12px 24px -8px rgba(0,0,0,0.4)'
    ].join(';')

    var open = false
    button.addEventListener('click', function () {
      open = !open
      panel.style.display = open ? 'block' : 'none'
      if (open) {
        iframe.contentWindow.postMessage({ type: 'widget:refresh' }, BOOKING_ORIGIN)
      }
    })
    root.appendChild(button)
    document.body.appendChild(root)
  }

  if (document.body) {
    mount()
  } else {
    document.addEventListener('DOMContentLoaded', mount)
  }

  window.addEventListener('message', applyResize)
})()
`

type WidgetHandler struct{}

func NewWidgetHandler() *WidgetHandler {
	return &WidgetHandler{}
}

func (h *WidgetHandler) ServeScript(w http.ResponseWriter, r *http.Request) {
	script := strings.ReplaceAll(widgetScript, "__ERP_WIDGET_BOOKING_ORIGIN__", resolveWidgetBookingOrigin())
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	_, _ = w.Write([]byte(script))
}

func resolveWidgetBookingOrigin() string {
	for _, envKey := range []string{"WIDGET_BOOKING_ORIGIN", "BOOKING_BASE_URL", "BASE_URL"} {
		if raw := strings.TrimSpace(os.Getenv(envKey)); raw != "" {
			return strings.TrimSuffix(raw, "/")
		}
	}
	return ""
}
