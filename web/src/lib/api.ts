import { Todo, CreateTodoDto, UpdateTodoDto } from '@/types/todo';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/todos';

export async function getTodos(): Promise<Todo[]> {
  const res = await fetch(API_URL, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch todos');
  const data = await res.json();
  return data.data || [];
}

export async function getTodo(id: string): Promise<Todo> {
  const res = await fetch(`${API_URL}/${id}`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch todo');
  const data = await res.json();
  return data.data;
}

export async function createTodo(todo: CreateTodoDto): Promise<Todo> {
  const res = await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(todo),
  });
  if (!res.ok) throw new Error('Failed to create todo');
  const data = await res.json();
  return data.data;
}

export async function updateTodo(id: string, todo: UpdateTodoDto): Promise<Todo> {
  const res = await fetch(`${API_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(todo),
  });
  if (!res.ok) throw new Error('Failed to update todo');
  const data = await res.json();
  return data.data;
}

export async function deleteTodo(id: string): Promise<void> {
  const res = await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to delete todo');
}
