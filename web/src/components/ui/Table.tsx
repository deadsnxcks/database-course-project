import './Table.css'

export type TableColumn<T> = {
  key: keyof T | 'actions'
  title: string
  render?: (value: T[keyof T] | undefined, row: T) => React.ReactNode
}

type Props<T> = {
  columns: TableColumn<T>[]
  data: T[]
  onAdd?: () => void
  onSearch?: (value: string) => void
}

export function Table<T extends { id: number | string }>({
  columns,
  data,
  onAdd,
  onSearch,
}: Props<T>) {
  return (
    <div className="table-wrapper">
      <div className="table-toolbar">
        {onSearch && (
          <input
            className="table-search"
            placeholder="Поиск..."
            onChange={(e) => onSearch(e.target.value)}
          />
        )}

        {onAdd && (
          <button className="btn btn-primary" onClick={onAdd}>
            + Добавить
          </button>
        )}
      </div>

      <table className="table">
        <thead>
          <tr>
            {columns.map((col) => (
              <th key={String(col.key)}>{col.title}</th>
            ))}
          </tr>
        </thead>

        <tbody>
          {data.length === 0 && (
            <tr>
              <td colSpan={columns.length} className="empty">
                Нет данных
              </td>
            </tr>
          )}

          {data.map((row) => (
            <tr key={row.id}>
              {columns.map((col) => (
                <td key={String(col.key)}>
                  {col.render
                    ? col.render(col.key === 'actions' ? undefined : row[col.key as keyof T], row)
                    : col.key === 'actions'
                    ? null
                    : String(row[col.key as keyof T] ?? '')}
                </td>

              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
