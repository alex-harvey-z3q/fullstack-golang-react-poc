import { useEffect, useMemo, useRef, useState } from 'react'
import type { Task } from './types'
import { createTask, fetchTasks } from './api'

export default function App() {
  const [tasks, setTasks] = useState<Task[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [title, setTitle] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const acRef = useRef<AbortController | null>(null)

  // load tasks on mount (and whenever you choose to refresh)
  useEffect(() => {
    acRef.current?.abort()
    const ac = new AbortController()
    acRef.current = ac
    setError(null)
    fetchTasks(ac.signal)
      .then(setTasks)
      .catch((e) => setError(e.message))
    return () => ac.abort()
  }, [])

  const pending = tasks === null && !error
  const canSubmit = useMemo(
    () => title.trim().length > 0 && !submitting,
    [title, submitting]
  )

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!canSubmit) return
    setSubmitting(true)
    setError(null)
    try {
      const t = await createTask(title.trim())
      setTasks((prev) => (prev ? [...prev, t] : [t]))
      setTitle('')
    } catch (e: any) {
      setError(e.message ?? 'Failed to add task')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="page">
      <header className="header">
        <h1>Tasks</h1>
      </header>

      <section className="card">
        <form className="form-row" onSubmit={onSubmit}>
          <input
            className="input"
            type="text"
            placeholder="Add a new task…"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            aria-label="Task title"
          />
          <button
            className="btn"
            type="submit"
            disabled={!canSubmit}
            title="Add task">
            {submitting ? 'Adding…' : 'Add'}
          </button>
        </form>

        {error && <div className="error">⚠️ {error}</div>}

        {pending && <div className="muted">Loading…</div>}

        {!pending && !error && (
          <ul className="list">
            {tasks?.map((t) => (
              <li key={t.id} className="list-item">
                <input type="checkbox" checked={t.done} readOnly />
                <span className={t.done ? 'done' : ''}>{t.title}</span>
              </li>
            ))}
            {tasks?.length === 0 && (
              <li className="muted">No tasks yet. Add your first one above.</li>
            )}
          </ul>
        )}
      </section>
    </div>
  )
}
