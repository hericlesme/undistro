type Props = {
  children: React.ReactNode
  columns: number | string
  gap?: number
  noMargin?: boolean
}

export const Grid = ({ children, columns, gap = 10, noMargin }: Props) => {
  return (
    <div
      style={{
        display: 'grid',
        gap: `${gap}px`,
        gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))`,
        marginTop: noMargin ? 0 : '30px'
      }}
    >
      {children}
    </div>
  )
}
