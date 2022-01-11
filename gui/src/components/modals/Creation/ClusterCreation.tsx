import { useReducer, VFC } from 'react'
import { useState } from 'react'
import { DialogOverlay, DialogContent } from '@reach/dialog'
import { Wizard } from './Wizard/Wizard'
import { useModalContext } from '@/contexts/ModalContext'
import { Progress } from '@/components/modals/Creation/Progress'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'

enum CREATION_OPTION {
  WIZARD = 'Wizard',
  ADVANCED = 'Advanced'
}

enum CREATION_STEP {
  OPTIONS = 'OPTIONS',
  STEPS = 'STEPS',
  STATUS = 'STATUS'
}

type CreationOption = {
  type: CREATION_OPTION
  subtitle: string
  description: string
  onClick?: (e: any) => void
}

const CreationOption: VFC<CreationOption> = ({ subtitle, description, type, ...props }: CreationOption) => {
  const createOptionStyles = {
    container: styles[`createCluster${type}Container`],
    iconContainer: styles[`${type.toLowerCase()}IconContainer`],
    titleContainer: styles[`${type.toLowerCase()}TitlesContainer`],
    title: styles.titleCreateStrong,
    description: styles.titleCreateDescription
  }

  return (
    <div className={createOptionStyles.container} {...props}>
      <div className={createOptionStyles.iconContainer}></div>
      <div className={createOptionStyles.titleContainer}>
        <div className={createOptionStyles.title}>{type}</div>
        <div className={createOptionStyles.description}>
          <a>{subtitle}</a>
          <br />
          <a>{description}</a>
        </div>
      </div>
    </div>
  )
}

const ClusterCreation: VFC = () => {
  const [step, setStep] = useState(1)
  const { hideModal } = useModalContext()

  const reducer = (state, action) => {
    switch (action.type) {
      case 'SET_STEP':
        return { ...state, step: action.payload.step }
      case 'SET_CREATION_MODE':
        return {
          ...state,
          step: CREATION_STEP.STEPS,
          creationMode: action.payload.creationMode
        }
      case 'SET_STATUS':
        return {
          ...state,
          step: CREATION_STEP.STATUS,
          progress: action.payload.progress
        }
    }
  }

  const [creationState, dispatch] = useReducer(reducer, {
    step: CREATION_STEP.OPTIONS,
    creationMode: ''
  })

  const nextStep = () => {
    setStep(step + 1)
  }

  const prevStep = () => {
    if (step - 1 === 0) {
      dispatch({ type: 'SET_STEP', payload: { step: CREATION_STEP.OPTIONS } })
    } else {
      setStep(step - 1)
    }
  }

  const creationOptions: CreationOption[] = [
    {
      type: CREATION_OPTION.WIZARD,
      subtitle: 'The fastest way to start.',
      description: 'Create a cluster in just a few steps.'
    },
    {
      type: CREATION_OPTION.ADVANCED,
      subtitle: 'No shortcuts.',
      description: 'Control every aspect of the cluster.'
    }
  ]

  const setCreationOption = (type: CREATION_OPTION) => () => {
    dispatch({ type: 'SET_CREATION_MODE', payload: { creationMode: type } })
  }

  const renderClusterCreationOptions = () => {
    switch (creationState.step) {
      case CREATION_STEP.OPTIONS:
        return (
          <div className={styles.createClusterOptionsContainer}>
            {creationOptions.map((option: CreationOption) => (
              <CreationOption onClick={setCreationOption(option.type)} {...option} key={option.type} />
            ))}
          </div>
        )
      case CREATION_STEP.STEPS:
        return <Wizard dispatch={dispatch} step={{ value: step, next: nextStep, previous: prevStep }} />
      case CREATION_STEP.STATUS:
        return <Progress creationInfo={creationState.progress} />
    }
  }

  return (
    <DialogOverlay isOpen={true} className={styles.dialogOverlay}>
      <DialogContent aria-label="Create Cluster" className={styles.dialogContent}>
        <div className={styles.createClusterOptionscontainer}>
          <div className={styles.modalDialogTitleBar}>
            <div className={styles.modalDialogTitleContainer}>
              <span className={styles.dialogTitlesSpan}>
                <a className={styles.dialogTitleRegular}>CREATE</a>
                <a className={styles.dialogTitleBold}>CLUSTER</a>
              </span>
            </div>
            <div className={styles.closeWindowBtnContainer}>
              <button onClick={hideModal} className={styles.closeWindowBtn}></button>
            </div>
          </div>
          {renderClusterCreationOptions()}
        </div>
      </DialogContent>
    </DialogOverlay>
  )
}

export { ClusterCreation }
