import { ReactComponent as CreateClusterErrorIcon } from '@assets/images/createClusterErrorIcon.svg'

import './createClusterError.scss'

export type ClusterCreationError = {
  code: number | string
  message: string
}

type CreateClusterErrorProps = ClusterCreationError & {
  description: string
  heading: string
}

export const CreateClusterError = ({
  code,
  description,
  heading,
  message
}: CreateClusterErrorProps) => {
  return (
    <>
      <h2 className="title-box">process</h2>
      <div className="cluster-error-container">
        <h3 className="cluster-error-heading">{heading}</h3>
        <CreateClusterErrorIcon className="cluster-error-icon" />
        <p className="cluster-error-description">{description}</p>
        <p className="cluster-error-message">
          {code} - {message}
        </p>
      </div>
    </>
  )
}
