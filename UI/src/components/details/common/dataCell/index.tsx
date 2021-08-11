import React from 'react'

import './index.scss'

type DataCellProps = {
  borderless?: boolean
  children: React.ReactNode
  label: string
}

export const DataCell = ({ borderless, children, label }: DataCellProps) => {
  return (
    <div className={`data-cell-container ${borderless ? 'borderless' : ''}`}>
      <div className="data-cell-label">{label}</div>
      <div className="data-cell-content">{children}</div>
    </div>
  )
}
