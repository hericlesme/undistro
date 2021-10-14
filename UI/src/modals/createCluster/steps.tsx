import React, { FC, ReactChild, useState } from 'react'
import Button from '@components/button'
import { CreateClusterError, ClusterCreationError } from './createClusterError'
import { CreateClusterProgress } from './createClusterProgress'
import { useEffect } from 'react'

type Props = {
  children: ReactChild[]
  handleAction: () => void
  handleClose: () => void
  handleBack: () => void
  error?: ClusterCreationError
  wizard?: boolean
  messages?: string[]
}

const Steps: FC<Props> = ({ children, handleAction, handleClose, handleBack, error, wizard, messages = [] }) => {
  const [index, setIndex] = useState<number>(0)
  const [hasFinished, setHasFinished] = useState(false)
  const hasError = !!error
  const step = wizard ? 3 : 7
  const isCreating = messages.filter(m => !!m).length > 0

  useEffect(() => {
    setHasFinished(
      messages.some(m => m.includes('SucessfulCreateKubeconfig') || m.includes('SucessfulCreateUserKubeconfig'))
    )
  }, [messages])

  useEffect(() => {
    if (index === step) handleAction()
    if (index < 0) handleBack()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [index])

  return (
    <>
      {hasError ? (
        <CreateClusterError
          code={error?.code!}
          description="An error occurred while creating the cluster."
          heading="The cluster could not be created"
          message={error?.message!}
        />
      ) : index === step ? (
        <CreateClusterProgress hasFinished={hasFinished} messages={messages} />
      ) : (
        children[index]
      )}

      <footer>
        {hasError ? (
          <>
            <div />
            <Button onClick={handleClose} variant="error" size="large" children="Close" />
          </>
        ) : (
          <>
            {isCreating ? (
              <div />
            ) : (
              <Button onClick={() => setIndex(index - 1)} variant="black" size="large" children="back" />
            )}
            <Button
              /* disabled={!hasFinished && isCreating} */
              onClick={() => {
                if (isCreating) handleClose()
                if (index !== step) setIndex(index + 1)
                if (hasFinished) window.location.reload()
              }}
              variant="primary"
              size="large"
              children={hasFinished ? 'Finish' : isCreating ? 'Close' : 'next'}
            />
          </>
        )}
      </footer>
    </>
  )
}

export default Steps
