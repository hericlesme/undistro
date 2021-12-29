import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import {
  ClusterInfo,
  InfraProvider,
  AddOns,
  Step,
  ControlPlane
} from '@/components/overviews/Clusters/Creation/Wizard/Steps'
import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'
import classNames from 'classnames'

const Wizard = ({ step }) => {
  const { watch, register, setValue, getValues, control } = useForm()
  const [currentSection, setCurrentSection] = useState('')

  const steps = [
    {
      title: 'Cluster',
      component: ClusterInfo
    },
    {
      title: 'Infrastructure Provider',
      component: InfraProvider
    },
    { title: 'Add-Ons', component: AddOns },
    {
      title: 'Control Plane',
      component: ControlPlane
    }
  ]

  useEffect(() => {
    console.log(steps[step.value - 1])
    setCurrentSection(steps[step.value - 1].title)
  }, [step.value])

  useEffect(() => {
    const subscription = watch((value, { name, type }) => console.log(value, name, type))
    return () => subscription.unsubscribe()
  }, [watch])

  const inputAreaStyles = classNames(styles.modalInputArea, {
    [styles.modalControlPlaneBlock]: currentSection === 'Control Plane'
  })

  const formActions = {
    register,
    setValue,
    getValues,
    control
  }

  return (
    <div className={styles.createClusterWizContainer}>
      <div className={styles.modalTitleContainer}>
        <a className={styles.modalCreateClusterTitle}>{currentSection}</a>
      </div>

      <div className={styles.modalContentContainer}>
        <div className={inputAreaStyles}>
          <form className={styles.modalForm} id="wizardClusterForm">
            {steps.map(({ component: Component }, index) => (
              <Step step={index + 1} currentStep={step.value}>
                <Component {...formActions} />
              </Step>
            ))}
          </form>
        </div>
      </div>
      <div className={styles.modalDialogButtonsContainer}>
        <div className={styles.leftButtonContainer}>
          <button onClick={step.previous} className={styles.borderButtonDefault}>
            <a>back</a>
          </button>
        </div>
        {step.value !== steps.length && (
          <div className={styles.rightButtonContainer}>
            <button onClick={step.next} className={styles.borderButtonSuccess}>
              <a>next</a>
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

export { Wizard }
