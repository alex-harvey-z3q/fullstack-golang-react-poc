import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TaskService } from './task.service';
import { Task } from './types';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule],
  // TEMP: inline template to rule-out templateUrl path issues
  template: `
    <div style="padding:16px; border: 2px dashed #f36; margin:16px;">
      <h1>Angular is mounted ✅</h1>
      <p>Below is the live Tasks UI.</p>

      <form (submit)="addTask(); $event.preventDefault()" style="display:flex; gap:8px; margin:12px 0;">
        <input [(ngModel)]="title" name="title" placeholder="Add a new task…" />
        <button [disabled]="!canSubmit">{{ submitting ? 'Adding…' : 'Add' }}</button>
      </form>

      <div *ngIf="error" style="color:#b00020;">⚠️ {{ error }}</div>
      <div *ngIf="pending" style="opacity:0.7;">Loading…</div>

      <ul *ngIf="!pending && !error">
        <li *ngFor="let t of tasks">
          <input type="checkbox" [checked]="t.done" disabled />
          <span [style.textDecoration]="t.done ? 'line-through' : 'none'">{{ t.title }}</span>
        </li>
        <li *ngIf="(tasks ?? []).length === 0" style="opacity:.7;">No tasks yet. Add your first one above.</li>
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
    console.log('[AppComponent] constructor');
    this.load();
  }

  get pending(): boolean {
    return this.tasks === null && !this.error;
  }

  get canSubmit(): boolean {
    return this.title.trim().length > 0 && !this.submitting;
  }

  load() {
    this.error = null;
    this.api.fetchTasks().subscribe({
      next: (items) => {
        console.log('[AppComponent] fetched tasks', items);
        this.tasks = items;
      },
      error: (e) => {
        console.error('[AppComponent] fetch error', e);
        this.error = e?.message ?? 'Failed to load tasks';
      }
    });
  }

  addTask() {
    if (!this.canSubmit) return;
    this.submitting = true;
    const title = this.title.trim();
    this.api.createTask(title).subscribe({
      next: (t) => {
        console.log('[AppComponent] created task', t);
        this.tasks = [...(this.tasks ?? []), t];
        this.title = '';
      },
      error: (e) => {
        console.error('[AppComponent] create error', e);
        this.error = e?.message ?? 'Failed to add task';
      },
      complete: () => (this.submitting = false)
    });
  }
}
