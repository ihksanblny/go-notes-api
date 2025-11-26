const API_URL = 'http://localhost:8080'

export async function getNotes() {
  const res = await fetch(`${API_URL}/notes`)
  if (!res.ok) {
    throw new Error('Gagal mengambil notes')
  }
  return res.json()
}

export async function createNote(payload) {
  const res = await fetch(`${API_URL}/notes`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Gagal membuat note')
  }

  return res.json()
}

export async function updateNote(id, payload) {
  const res = await fetch(`${API_URL}/notes/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json'},
    body: JSON.stringify(payload),
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Gagal mengupdate note')
  }

  return res.json()
}

export async function deleteNote(id) {
  const res = await fetch(`${API_URL}/notes/${id}`, {
    method: 'DELETE',
  })

  if (!res.ok && res.status !== 204) {
    const text = await res.text()
    throw new Error(text || 'Gagal menghapus note')
  }

  return true
}
