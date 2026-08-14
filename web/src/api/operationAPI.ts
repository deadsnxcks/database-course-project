const BASE_URL = 'http://localhost:8081/operation';

export interface Operation {
  id: number;
  title: string;
  createdAt: string;
}

export async function listOperations(): Promise<Operation[]> {
  const res = await fetch(BASE_URL);
  if (!res.ok) {
    throw await res.json();
  }
  return res.json();
}

export async function getOperation(id: number): Promise<Operation> {
  const res = await fetch(`${BASE_URL}/${id}`);
  if (!res.ok) {
    throw await res.json();
  }
  return res.json();
}

export async function createOperation(title: string): Promise<{ id: number }> {
  const res = await fetch(BASE_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title }),
  });

  if (!res.ok) {
    throw await res.json();
  }

  return res.json();
}

export async function updateOperation(
  id: number,
  payload: { title?: string }
): Promise<void> {
  const res = await fetch(`${BASE_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw await res.json();
  }
}

export async function deleteOperation(id: number): Promise<void> {
  const res = await fetch(`${BASE_URL}/${id}`, {
    method: 'DELETE',
  });

  if (!res.ok) {
    throw await res.json();
  }
}
