import type { VFC } from 'react'
import { useState } from 'react'
import { DialogOverlay, DialogContent } from '@reach/dialog'
import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'
import { Wizard } from './Wizard/Wizard'

type ClusterCreationProps = {
  isOpen: boolean
}

enum CreationOptionType {
  wizard = 'Wizard',
  advanced = 'Advanced'
}

type CreationOption = {
  type: CreationOptionType
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

const ClusterCreation: VFC<ClusterCreationProps> = ({ isOpen }: ClusterCreationProps) => {
  const [step, setStep] = useState(1)
  const [creationMode, setCreationMode] = useState('')

  const nextStep = () => {
    setStep(step + 1)
  }

  const prevStep = () => {
    if (step - 1 === 0) {
      setCreationMode('')
    } else {
      setStep(step - 1)
    }
  }

  const creationOptions: CreationOption[] = [
    {
      type: CreationOptionType.wizard,
      subtitle: 'The fastest way to start.',
      description: 'Create a cluster in just a few steps.'
    },
    {
      type: CreationOptionType.advanced,
      subtitle: 'No shortcuts.',
      description: 'Control every aspect of the cluster.'
    }
  ]

  const handleOptionClick = (type: CreationOptionType) => e => {
    setCreationMode(type)
  }

  const renderClusterCreationOptions = () => {
    switch (creationMode) {
      case CreationOptionType.wizard:
        return <Wizard step={{ value: step, next: nextStep, previous: prevStep }} />
      case CreationOptionType.advanced:
        return <div>Advanced</div>
      default:
        return (
          <div className={styles.createClusterOptionsContainer}>
            {creationOptions.map((option: CreationOption) => (
              <CreationOption onClick={handleOptionClick(option.type)} {...option} key={option.type} />
            ))}
          </div>
        )
    }
  }

  return (
    <DialogOverlay isOpen={isOpen} className={styles.dialogOverlay}>
      <DialogContent className={styles.dialogContent}>
        <div className={styles.createClusterOptionscontainer}>
          <div className={styles.modalDialogTitleBar}>
            <div className={styles.modalDialogTitleContainer}>
              <span className={styles.dialogTitlesSpan}>
                <a className={styles.dialogTitleRegular}>CREATE</a>
                <a className={styles.dialogTitleBold}>CLUSTER</a>
              </span>
            </div>

            <div className={styles.closeWindowBtnContainer}>
              <button className={styles.closeWindowBtn}></button>
            </div>
          </div>
          {renderClusterCreationOptions()}
        </div>
      </DialogContent>
    </DialogOverlay>
  )
}

export { ClusterCreation }
