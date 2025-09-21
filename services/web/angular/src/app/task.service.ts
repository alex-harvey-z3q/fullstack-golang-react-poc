import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import type { Task } from './types';

@Injectable({ providedIn: 'root' })
export class TaskService {
  private http = inject(HttpClient);

  fetchTasks(): Observable<Task[]> {
    return this.http.get<Task[]>('/api/tasks');
  }

  createTask(title: string): Observable<Task> {
    return this.http.post<Task>('/api/tasks', { title });
  }
}
