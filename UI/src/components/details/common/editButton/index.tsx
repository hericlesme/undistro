import { ComponentPropsWithoutRef } from 'react'

import './index.scss'

type Props = Omit<ComponentPropsWithoutRef<'button'>, 'className'> & {
  editing: boolean
}

export const EditButton = ({ editing, ...otherProps }: Props) => {
  return (
    <button
      {...otherProps}
      className={`edit-button ${editing ? 'editing' : ''}`}
    >
      <svg
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
        />
      </svg>
    </button>
  )
}
