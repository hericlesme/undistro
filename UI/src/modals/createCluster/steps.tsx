import React, { FC, ReactChild, useState } from 'react'
import Button from '@components/button'
import { CreateClusterError } from './createClusterError'
import type {ClusterCreationError} from './createClusterError'

type Props = {
  children: ReactChild[]
  handleAction: () => void
  handleClose: () => void
  handleBack: () => void
  error?: ClusterCreationError,
  wizard?: boolean
}

const Steps: FC<Props> = ({
  children,
  handleAction,
  handleClose,
  handleBack,
  error,
  wizard
}) => {
  const [index, setIndex] = useState<number>(0)
  const hasError=!!error
  const step = wizard ? 3 : 7

  if (index === step) handleAction()
  if (index < 0) handleBack()

  return (
    <>
      {hasError ? (
        <CreateClusterError
          code={error?.code!}
          description="An error occurred while creating the cluster."
          heading="The cluster could not be created"
          message={error?.message!}
        />
      ) : (
        children[index]
      )}

      <footer>
        {hasError ? (
          <>
            <div />
            <Button
              onClick={handleClose}
              type="error"
              size="large"
              children="Close"
            />
          </>
        ) : (
          <>
            <Button
              onClick={() => setIndex(index - 1)}
              type="black"
              size="large"
              children="back"
            />
            <Button
              onClick={() => setIndex(index + 1)}
              type="primary"
              size="large"
              children="next"
            />
          </>
        )}
      </footer>
    </>
  )
}

export default Steps
