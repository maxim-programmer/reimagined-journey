const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

function getToken() {
  return localStorage.getItem('kb_token') || ''
}

function authHeaders() {
  return {
    Authorization: `Bearer ${getToken()}`,
  }
}

export async function register(login, password) {
  const response = await fetch(`${BASE_URL}/api/v1/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password }),
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error || 'Ошибка регистрации')
  return data
}

export async function login(loginVal, password) {
  const response = await fetch(`${BASE_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login: loginVal, password }),
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error || 'Ошибка входа')
  return data
}

export async function logout() {
  await fetch(`${BASE_URL}/api/v1/auth/logout`, {
    method: 'POST',
    headers: authHeaders(),
  })
  localStorage.removeItem('kb_token')
}

export async function getMe() {
  const response = await fetch(`${BASE_URL}/api/v1/auth/me`, {
    headers: authHeaders(),
  })
  if (!response.ok) return null
  return response.json()
}

export async function uploadDocument(file) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await fetch(`${BASE_URL}/api/v1/documents/upload`, {
    method: 'POST',
    headers: authHeaders(),
    body: formData,
  })

  const data = await response.json()
  if (!response.ok) throw new Error(data.error || 'Ошибка при загрузке файла')
  return data
}

export async function listDocuments() {
  const response = await fetch(`${BASE_URL}/api/v1/documents`, {
    headers: authHeaders(),
  })
  if (!response.ok) throw new Error('Ошибка при получении списка документов')
  return response.json()
}

export async function searchDocuments(query) {
  const response = await fetch(`${BASE_URL}/api/v1/search?q=${encodeURIComponent(query)}`, {
    headers: authHeaders(),
  })
  if (!response.ok) {
    const data = await response.json()
    throw new Error(data.error || 'Ошибка при поиске')
  }
  return response.json()
}

export async function getSearchHistory() {
  const response = await fetch(`${BASE_URL}/api/v1/history`, {
    headers: authHeaders(),
  })
  if (!response.ok) return []
  return response.json()
}

export async function clearSearchHistory() {
  await fetch(`${BASE_URL}/api/v1/history`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
}