import type { Task } from './types'

const JSON_HEADERS = { 'Content-Type': 'application/json' }

export async function fetchTasks(signal?: AbortSignal) {
  const res = await fetch('/api/tasks', { signal })
  if (!res.ok) {
    let msg = `GET /api/tasks failed: ${res.status}`
    try {
      const body = await res.json()
      if (body?.error) msg += ` ${body.error}`
    } catch {}
    throw new Error(msg)
  }
  return res.json()
}

// Requires POST /api/tasks on the backend.
// If not implemented yet, this will throw and the UI will show the error.
export async function createTask(title: string): Promise<Task> {
  const res = await fetch('/api/tasks', {
    method: 'POST',
    headers: JSON_HEADERS,
    body: JSON.stringify({ title })
  })
  if (!res.ok) {
    const body = await res.text().catch(() => '')
    throw new Error(`POST /api/tasks failed: ${res.status} ${body}`)
  }
  return res.json()
}
