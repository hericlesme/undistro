/* eslint-disable jsx-a11y/no-access-key */
/* eslint-disable react-hooks/exhaustive-deps */
import React, { FC, useEffect, useState } from 'react'
import store from '../store'
import CreateCluster from '@components/modals/cluster'
import Infra from '@components/modals/infrastructureProvider'
import ControlPlane from '@components/modals/controlPlane'
import { generateId } from 'util/helpers'
import Steps from './steps'
import Api from 'util/api'
import { TypeOption, TypeWorker, TypeSelectOptions } from '../../types/cluster'
import { TypeModal, apiResponse } from '../../types/generic'

const ClusterWizard: FC<TypeModal> = ({ handleClose }) => {
  const body = store.useState((s: any) => s.body)
  const [accessKey, setAccesskey] = useState<string>('')
  const [secret, setSecret] = useState<string>('')
  const [region, setRegion] = useState<string>('')
  const [clusterName, setClusterName] = useState<string>('')
  const [namespace, setNamespace] = useState<string>('')
  const [provider, setProvider] = useState<string>('')
  const [flavor, setFlavor] = useState<string>('')
  const [k8sVersion, setK8sVersion] = useState<string>('')
  const [replicas, setReplicas] = useState<number>(0)
  const [infraNode, setInfraNode] = useState<boolean>(false)
  const [workers, setWorkers] = useState<TypeWorker[]>([])
  const [machineTypes, setMachineTypes] = useState<string>('')
  const [memory, setMemory] = useState<string>('')
  const [cpu, setCpu] = useState<string>('')
  const [replicasWorkers, setReplicasWorkers] = useState<number>(0)
  const [memoryWorkers, setMemoryWorkers] = useState<string>('')
  const [cpuWorkers, setCpuWorkers] = useState<string>('')
  const [machineTypesWorkers, setMachineTypesWorkers] = useState<string>('')
  const [regionOptions, setRegionOptions] = useState<[]>([])
  const [flavorOptions, setFlavorOptions] = useState<TypeOption[]>([])
  const [cpuOptions, setCpuOptions] = useState<TypeOption[]>()
  const [memOptions, setMemOptions] = useState<TypeOption[]>()
  const [MachineOptions, setMachineOptions] = useState<TypeOption[]>()
  const [k8sOptions, setK8sOptions] = useState<TypeSelectOptions>()
  const [sshKey, setSshKey] = useState<string>('')
  const [session, setSession] = useState<string>('')
  const providerOptions = [{ value: provider, label: 'aws' }]

  const handleAction = () => {
    const getWorkers = workers.map(elm => ({
      machineType: elm.machineType,
      replicas: elm.replicas,
      infraNode: elm.infraNode
    }))

    const data = {
      apiVersion: 'app.undistro.io/v1alpha1',
      kind: 'Cluster',
      metadata: {
        name: clusterName,
        namespace: namespace
      },

      spec: {
        kubernetesVersion: k8sVersion,
        controlPlane: {
          machineType: machineTypes,
          replicas: replicas
        },

        infrastructureProvider: {
          flavor: flavor,
          name: provider,
          region: region,
          sshKey: sshKey
        },

        workers: getWorkers
      }
    }

    const dataPolicies = {
      apiVersion: 'app.undistro.io/v1alpha1',
      kind: 'DefaultPolicies',
      metadata: {
        name: 'defaultpolicies-undistro',
        namespace: namespace
      },
      spec: {
        clusterName: clusterName
      }
    }

    Api.Cluster.post(data, namespace)
    setTimeout(() => {
      Api.Cluster.postPolicies(dataPolicies, namespace)
    }, 600)
  }

  const getSecrets = (secretRef: string) => {
    Api.Secret.list(secretRef).then(res => {
      setAccesskey(atob(res.data.accessKeyID))
      setSecret(atob(res.data.secretAccessKey))
      setRegion(atob(res.data.region))
    })
  }

  const createWorkers = () => {
    setWorkers([
      ...workers,
      {
        id: generateId(),
        machineType: machineTypesWorkers,
        replicas: replicasWorkers,
        infraNode: infraNode
      }
    ])
  }

  const deleteWorkers = (id: any) => {
    const newWorkers = workers.filter(item => item.id !== id)
    setWorkers(newWorkers)
  }

  const getMachines = () => {
    Api.Provider.list('awsmachines')
      .then(res => {
        const name = res.items.map((elm: any) => ({ label: elm.metadata.name, value: elm.metadata.name }))
        const cpu = res.items.map((elm: any) => ({ label: elm.spec.vcpus, value: elm.spec.vcpus }))
        const mem = res.items.map((elm: any) => ({ label: elm.spec.memory, value: elm.spec.memory }))

        setMachineOptions(name)
        setMemOptions(mem)
        setCpuOptions(cpu)
      })
  }

  const getProviders = () => {
    Api.Provider.list('providers')
      .then(res => {
        const newArray = res.items.filter((elm: any) => { return elm.spec.category.includes('infra') })
        setProvider(newArray[0].metadata.name)
        setRegionOptions(newArray[0].status.regionNames.map((elm: string) => ({ value: elm, label: elm })))
        getSecrets(newArray[0].spec.secretRef.name)
        return newArray
      })
  }

  const getFlavors = async () => {
    const res = await Api.Provider.list('flavors')
    const names = res.items.map((elm: any) => ({ label: elm.metadata.name, value: elm.metadata.name }))
    const parse = (data: apiResponse): TypeSelectOptions => {
      return data.reduce<TypeSelectOptions>((acc, curr) => {
        acc[curr.metadata.name] = {
          selectOptions: curr.spec.supportedK8SVersions.map(elm => ({ label: elm, value: elm}))
      }

        return acc
      }, {})
    }

    const parseData = parse(res.items)  
    setK8sOptions(parseData)
    setFlavorOptions(names)
  }

  useEffect(() => {
    getProviders()
    getFlavors()
    getMachines()
  }, [])

  return (
    <>
      <header>
        <h3 className="title">
          <i className="wizard-text">Wizard</i> <span>{body.title}</span>{' '}
          {body.ndTitle}
        </h3>
        <i onClick={handleClose} className="icon-close" />
      </header>
      <div className='box'>
        <Steps wizard handleClose={handleClose} handleAction={() => handleAction()}>
          <CreateCluster 
            clusterName={clusterName}
            setClusterName={setClusterName}
            namespace={namespace}
            setNamespace={setNamespace}
            provider={provider}
            setProvider={setProvider}
            providerOptions={providerOptions}
            region={region}
            setRegion={setRegion}
            regionOptions={regionOptions}
            accessKey={accessKey}
            setAccesskey={setAccesskey}
            secret={secret}
            setSecret={setSecret}
            session={session}
            setSession={setSession}
          />

          <Infra 
            provider={provider}
            setProvider={setProvider}
            providerOptions={providerOptions}
            flavor={flavor}
            setFlavor={setFlavor}
            flavorOptions={flavorOptions}
            region={region}
            setRegion={setRegion}
            regionOptions={regionOptions}
            k8sVersion={k8sVersion}
            setK8sVersion={setK8sVersion}
            k8sOptions={k8sOptions}
            sshKey={sshKey}
            setSshKey={setSshKey}
          />

          <ControlPlane 
            replicas={replicas}
            setReplicas={setReplicas}
            cpu={cpu}
            setCpu={setCpu}
            getCpu={cpuOptions}
            memory={memory}
            setMemory={setMemory}
            getMem={memOptions}
            machineTypes={machineTypes}
            setMachineTypes={setMachineTypes}
            getMachineTypes={MachineOptions}
            infraNode={infraNode}
            setInfraNode={setInfraNode}
            workers={workers}
            createWorkers={createWorkers}
            deleteWorkers={deleteWorkers}
            replicasWorkers={replicasWorkers}
            setReplicasWorkers={setReplicasWorkers}
            cpuWorkers={cpuWorkers}
            setCpuWorkers={setCpuWorkers}
            memoryWorkers={memoryWorkers}
            setMemoryWorkers={setMemoryWorkers}
            machineTypesWorkers={machineTypesWorkers}
            setMachineTypesWorkers={setMachineTypesWorkers}
            clusterName={clusterName}
          />



        </Steps>
      </div>
    </>
  )
}

export default ClusterWizard
