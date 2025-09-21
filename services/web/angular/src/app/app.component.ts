import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { TaskService } from './task.service';
import { Task } from './types';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="page">
      <h1>Tasks (Angular)</h1>

      <form (submit)="addTask(); $event.preventDefault()" class="form-row">
        <input
          class="input"
          type="text"
          placeholder="Add a new task…"
          [(ngModel)]="title"
          name="title"
          aria-label="Task title" />
        <button class="btn" [disabled]="!canSubmit">{{ submitting ? 'Adding…' : 'Add' }}</button>
      </form>

      <div *ngIf="error" class="error">⚠️ {{ error }}</div>
      <div *ngIf="pending" class="muted">Loading…</div>

      <ul *ngIf="!pending && !error" class="list">
        <li *ngFor="let t of tasks" class="list-item">
          <input type="checkbox" [checked]="t.done" disabled />
          <span [class.done]="t.done">{{ t.title }}</span>
        </li>
        <li *ngIf="(tasks ?? []).length === 0" class="muted">
          No tasks yet. Add your first one above.
        </li>
      </ul>
    </div>
  `
})
export class AppComponent {
  private api = inject(TaskService);

  tasks: Task[] | null = null;
  error: string | null = null;
  title = '';
  submitting = false;

  constructor() {
    this.load();
  }

  get pending(): boolean {
    return this.tasks === null && !this.error;
  }

  get canSubmit(): boolean {
    return this.title.trim().length > 0 && !this.submitting;
  }

  /** Extract `{ "error": "..." }` from backend responses when available. */
  private extractServerError(e: unknown): string | null {
    const err = e as HttpErrorResponse | undefined;
    if (!err) return null;

    const payload = err.error;
    if (payload && typeof payload === 'object' && 'error' in payload) {
      const msg = (payload as any).error;
      if (typeof msg === 'string' && msg.trim().length > 0) return msg;
    }
    if (typeof payload === 'string' && payload.trim().length > 0) {
      return payload;
    }
    return null;
  }

  load() {
    this.error = null;
    this.api.fetchTasks().subscribe({
      next: (items) => {
        this.tasks = items;
      },
      error: (e) => {
        const serverError = this.extractServerError(e);
        this.error = serverError || (e as any)?.message || 'Failed to load tasks';
      }
    });
  }

  addTask() {
    if (!this.canSubmit) return;
    this.submitting = true;
    const title = this.title.trim();
    this.api.createTask(title).subscribe({
      next: (t) => {
        this.tasks = [...(this.tasks ?? []), t];
        this.title = '';
      },
      error: (e) => {
        const serverError = this.extractServerError(e);
        this.error = serverError || (e as any)?.message || 'Failed to add task';
      },
      complete: () => (this.submitting = false)
    });
  }
}
