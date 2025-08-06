export async function withLoading<T> (state: { loading: boolean }, fn: () => Promise<T>): Promise<T> {
  state.loading = true
  try {
    return await fn()
  } finally {
    state.loading = false
  }
}
