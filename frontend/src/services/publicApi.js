const PUBLIC_BASE = '/api/public/v1'
const USER_BASE = '/api/v1'

async function parseOrThrow(res) {
  const json = await res.json().catch(() => ({}))
  if (!res.ok) {
    const err = new Error(json.message || `HTTP ${res.status}`)
    err.status = res.status
    err.data = json
    throw err
  }
  return json
}

export async function publicGet(path, params = {}) {
  const url = new URL(`${PUBLIC_BASE}${path}`, window.location.origin)
  Object.entries(params).forEach(([k, v]) => {
    if (v !== null && v !== undefined && v !== '') url.searchParams.set(k, v)
  })
  const res = await fetch(url.pathname + url.search, {
    headers: { Accept: 'application/json' },
  })
  return parseOrThrow(res)
}

function getUserToken() {
  return localStorage.getItem('user_token')
}

export async function userApi(path, options = {}) {
  const token = getUserToken()
  const headers = {
    Accept: 'application/json',
    ...(options.body ? { 'Content-Type': 'application/json' } : {}),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...(options.headers || {}),
  }
  const res = await fetch(`${USER_BASE}${path}`, { ...options, headers })
  if (res.status === 401) {
    localStorage.removeItem('user_token')
    localStorage.removeItem('user_data')
    window.location.href = '/connexion'
    throw new Error('Session expirée')
  }
  return parseOrThrow(res)
}

export function userDownloadUrl(path) {
  const token = getUserToken()
  return { url: `${USER_BASE}${path}`, token }
}
