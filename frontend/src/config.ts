export interface AppConfig {
  apiBaseUrl: string
}

let config: AppConfig = {
  apiBaseUrl: 'https://api.example.local/api/v1',
}

export async function loadConfig (): Promise<AppConfig> {
  try {
    const res = await fetch('/config.json')
    if (res.ok) {
      config = await res.json()
    }
  } catch {
    // ignore errors and use defaults
  }
  return config
}

export function getConfig (): AppConfig {
  return config
}
