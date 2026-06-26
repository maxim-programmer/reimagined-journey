const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export async function uploadDocument(file) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await fetch(`${BASE_URL}/api/v1/documents/upload`, {
    method: 'POST',
    body: formData,
  })

  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.error || 'Ошибка при загрузке файла')
  }

  return data
}

export async function listDocuments() {
  const response = await fetch(`${BASE_URL}/api/v1/documents`)

  if (!response.ok) {
    throw new Error('Ошибка при получении списка документов')
  }

  return response.json()
}