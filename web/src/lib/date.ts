export const formatDate = (value?: string | number | null): string => {
  if (!value) return '—';

  const date = new Date(value);
  if (isNaN(date.getTime())) return '—';

  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date);
};

export const toServerDate = (value?: string | null): string | undefined => {
  if (!value) return undefined
 
  // Если пояс уже указан (Z или ±HH:MM) — строку не трогаем
  const hasZone = /(?:Z|[+-]\d{2}:\d{2})$/.test(value)
  const date = hasZone ? new Date(value) : new Date(value.replace(/Z$/, ''))
 
  if (isNaN(date.getTime())) return undefined
 
  return date.toISOString()
}