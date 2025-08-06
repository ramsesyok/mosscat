export interface ToastOptions {
  type?: 'success' | 'error' | 'info'
  timeout?: number
}

export interface ToastEvent {
  message: string
  type: 'success' | 'error' | 'info'
  timeout: number
}
