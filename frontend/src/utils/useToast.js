import { reactive } from 'vue'

const toasts = reactive([])
let nextId = 0

export function useToast() {
  function add(message, type = 'info', duration = 3500) {
    const id = ++nextId
    toasts.push({ id, message, type })
    setTimeout(() => remove(id), duration)
  }

  function remove(id) {
    const idx = toasts.findIndex((t) => t.id === id)
    if (idx !== -1) toasts.splice(idx, 1)
  }

  return {
    toasts,
    showSuccess: (msg) => add(msg, 'success'),
    showError:   (msg) => add(msg, 'error', 5000),
    showInfo:    (msg) => add(msg, 'info'),
    remove,
  }
}
