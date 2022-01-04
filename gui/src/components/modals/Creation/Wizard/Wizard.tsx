import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { ClusterInfo, InfraProvider, AddOns, Step, ControlPlane } from '@/components/modals/Creation/Wizard/Steps'
import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import classNames from 'classnames'
import { useMutate } from '@/hooks/query'
import { ClusterCreationData } from '@/types/cluster'

const Wizard = ({ step }) => {
  const { watch, register, setValue, getValues, handleSubmit, control } = useForm()
  const [currentSection, setCurrentSection] = useState('')
  const createCluster = useMutate({ url: `/api/clusters/create`, method: 'post' })
  const createPolicy = useMutate({ url: `/api/clusters/policies`, method: 'post' })

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
    setCurrentSection(steps[step.value - 1].title)
  }, [step.value])

  useEffect(() => {
    const subscription = watch((value, { name, type }) => console.log(value, name, type))
    return () => subscription.unsubscribe()
  }, [watch])

  const inputAreaStyles = classNames(styles.modalInputArea, {
    [styles.modalControlPlaneBlock]: currentSection === 'Control Plane'
  })

  const onSubmit = (data, e) => {
    e.preventDefault()

    const clusterData: ClusterCreationData = {
      apiVersion: 'app.undistro.io/v1alpha1',
      kind: 'Cluster',
      metadata: {
        name: data.clusterName,
        namespace: data.clusterNamespace
      },
      spec: {
        kubernetesVersion: data.infraProviderK8sVersion,
        controlPlane: {
          machineType: data.controlPlaneMachineType,
          replicas: data.controlPlaneReplicas
        },
        infrastructureProvider: {
          flavor: data.infraProviderFlavor,
          name: data.clusterProvider,
          region: data.clusterDefaultRegion,
          sshKey: data.infraProviderSshKey
        }
      }
    }

    if (data.workers) {
      clusterData.spec.workers = data.workers
    }

    if (data.infraProviderID && data.infraProviderCIDR) {
      clusterData.spec.network = {
        vpc: {
          id: data.infraProviderID,
          cidrBlock: data.infraProviderCIDR
        }
      }
    }

    const dataPolicies = {
      apiVersion: 'app.undistro.io/v1alpha1',
      kind: 'DefaultPolicies',
      metadata: {
        name: `defaultpolicies-${data.clusterName}`,
        namespace: data.clusterNamespace
      },
      spec: {
        clusterName: data.clusterName
      }
    }

    const identity = {
      apiVersion: 'app.undistro.io/v1alpha1',
      kind: 'Identity',
      metadata: {
        name: `identity-${data.clusterName}`,
        namespace: data.clusterNamespace
      },
      spec: {
        clusterName: data.clusterName
      }
    }

    createCluster.mutate(clusterData)
    createPolicy.mutate(dataPolicies)
  }
  const onError = (errors, e) => {
    console.log(errors)
    console.log(e)
  }

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
          <form onSubmit={handleSubmit(onSubmit, onError)} className={styles.modalForm} id="wizardClusterForm">
            {steps.map(({ component: Component }, index) => (
              <Step step={index + 1} currentStep={step.value}>
                <Component {...formActions} />
              </Step>
            ))}
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
              {step.value === steps.length && (
                <div className={styles.rightButtonContainer}>
                  <button type="submit" className={styles.borderButtonSuccess}>
                    <a>Create</a>
                  </button>
                </div>
              )}
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export { Wizard }
