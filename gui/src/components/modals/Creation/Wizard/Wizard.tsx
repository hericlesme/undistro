import { Dispatch, useEffect, useState, VFC } from 'react'
import { useForm } from 'react-hook-form'
import { ClusterInfo, InfraProvider, AddOns, Step, ControlPlane } from '@/components/modals/Creation/Wizard/Steps'
import { useMutate } from '@/hooks/query'
import { ClusterCreationData } from '@/types/cluster'
import styles from '@/components/modals/Creation/ClusterCreation.module.css'

type WizardProps = {
  dispatch: Dispatch<any>
  step: {
    value: number
    next: () => void
    previous: () => void
  }
}
const Wizard: VFC<WizardProps> = ({ step, dispatch }: WizardProps) => {
  const { register, setValue, getValues, handleSubmit, control } = useForm()
  const [currentSection, setCurrentSection] = useState('')

  const createCluster = useMutate({
    url: `/api/clusters/create`,
    method: 'post',
    invalidate: '/api/clusters'
  })

  // const createPolicy = useMutate({ url: `/api/clusters/policies`, method: 'post' })

  const steps = [
    // { title: 'Infra Network - VPC', component: InfraNetVPC },
    // { title: 'Kubernetes Network', component: KubernetesNetwork },
    { title: 'Cluster', component: ClusterInfo },
    { title: 'Infrastructure Provider', component: InfraProvider },
    { title: 'Add-Ons', component: AddOns },
    { title: 'Control Plane', component: ControlPlane }
  ]
  useEffect(() => {
    setCurrentSection(steps[step.value - 1].title)
  }, [step.value])

  const inputStyles: { [key: string]: string } = {
    'Control Plane': styles.modalControlPlaneBlock,
    Progress: styles.modalProgressArea
  }

  const inputAreaStyles = inputStyles[currentSection] || styles.modalInputArea

  const onSubmit = async (data, e) => {
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
          replicas: Number(data.controlPlaneReplicas)
        },
        infrastructureProvider: {
          flavor: data.infraProviderFlavor,
          name: data.clusterProvider,
          region: data.clusterDefaultRegion,
          sshKey: data.infraProviderSshKey
        },
        network: {
          vpc: {
            id: data.infraProviderID,
            cidrBlock: data.infraProviderCIDR
          }
        },
        workers: data.workers
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

    let res: any = await createCluster.mutateAsync(JSON.stringify(clusterData))
    let payload = {
      cluster: data.clusterName,
      namespace: data.clusterNamespace,
      status: res.status == 200 ? 'success' : 'failure'
    }
    dispatch({ type: 'SET_STATUS', payload: { progress: payload } })
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

      <form
        autoComplete="off"
        onSubmit={handleSubmit(onSubmit, onError)}
        id="wizardClusterForm"
        className={styles.modalContentContainer}
      >
        <div className={inputAreaStyles}>
          <div className={styles.modalForm}>
            {steps.map(({ component: Component }, index) => (
              <Step key={`step-${index}`} step={index + 1} currentStep={step.value}>
                <Component {...formActions} />
              </Step>
            ))}
          </div>
        </div>
        <div className={styles.modalDialogButtonsContainer}>
          <div className={styles.leftButtonContainer}>
            <button type="button" onClick={step.previous} className={styles.borderButtonDefault}>
              <a>back</a>
            </button>
          </div>
          {step.value !== steps.length && (
            <div className={styles.rightButtonContainer}>
              <button type="button" onClick={step.next} className={styles.borderButtonSuccess}>
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
  )
}

export { Wizard }
