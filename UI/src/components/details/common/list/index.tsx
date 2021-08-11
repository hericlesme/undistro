import React from 'react'

import './index.scss'

type Props = {
  children: React.ReactNode
}

export const List = ({ children }: Props) => {
  return React.Children.count(children) > 0 ? <ul className="details-list">{children}</ul> : null
}
