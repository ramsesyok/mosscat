import type { ToastEvent, ToastOptions } from '@/types/toast'
import mitt from 'mitt'

const emitter = mitt<{ toast: ToastEvent }>()

export function useToast () {
  const show = (message: string, options?: ToastOptions) => {
    emitter.emit('toast', {
      message,
      type: options?.type ?? 'info',
      timeout: options?.timeout ?? 3000,
    })
  }

  const success = (msg: string) => show(msg, { type: 'success' })
  const error = (msg: string) => show(msg, { type: 'error' })
  const info = (msg: string) => show(msg, { type: 'info' })

  return { success, error, info }
}

export function onToast (handler: (event: ToastEvent) => void) {
  emitter.on('toast', handler)
}

export function offToast (handler: (event: ToastEvent) => void) {
  emitter.off('toast', handler)
}
