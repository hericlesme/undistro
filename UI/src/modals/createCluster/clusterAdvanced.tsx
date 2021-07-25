/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable jsx-a11y/no-access-key */
import React, { FC, useState, useEffect } from 'react'
import store from '../store'
import CreateCluster from '@components/modals/cluster'
import InfraNetwork from '@components/modals/infraNetwork'
import K8sNetwork from '@components/modals/kubernetesNetwork'
import Infra from '@components/modals/infrastructureProvider'
import ControlPlane from '@components/modals/controlPlane'
import Bastion from '@components/modals/bastion'
import Workers from '@components/modals/workersAdvanced'
import Steps from './steps'
import Api from 'util/api'
import { TypeOption, TypeSelectOptions, TypeSubnet, TypeTaints } from '../../types/cluster'
import { TypeModal, apiResponse } from '../../types/generic'

type TypeWorkers = {
  replicas: number
  machineType: string
  subnet: string
  labels: {}[]
  providerTags: {}[]
  taints: TypeTaints[]
  autoscaling: {
    enabled: boolean
    minSize: number
    maxSize: number
  }
}

const ClusterAdvanced: FC<TypeModal> = ({ handleClose }) => {
  const body = store.useState((s: any) => s.body)
  const [accessKey, setAccesskey] = useState<string>('')
  const [secret, setSecret] = useState<string>('')
  const [region, setRegion] = useState<string>('')
  const [clusterName, setClusterName] = useState<string>('')
  const [namespace, setNamespace] = useState<string>('')
  const [provider, setProvider] = useState<string>('')
  const [flavor, setFlavor] = useState<string>('')
  const [k8sVersion, setK8sVersion] = useState<string>('')
  const [id, setId] = useState<string>('')
  const [idSubnet, setIdSubnet] = useState<string>('')
  const [isPublic, setIsPublic] = useState<boolean>(false)
  const [zone, setZone] = useState<string>('')
  const [cidr, setCidr] = useState<string>('')
  const [cidrBastion, setCidrBastion] = useState<string>('')
  const [cidrs, setCidrs] = useState<string[]>()
  const [cidrSubnet, setCidrSubnet] = useState<string>('')
  const [serverPort, setServerPort] = useState<string>('')
  const [serviceDomain, setServiceDomain] = useState<string>('')
  const [podsRanges, setPodsRanges] = useState<string>('')
  const [serviceRanges, setServiceRanges] = useState<string>('')
  const [enabled, setEnabled] = useState<boolean>(false)
  const [ingress, setIngress] = useState<boolean>(false)
  const [internalLB, setInternalLB] = useState<boolean>(false)
  const [replicas, setReplicas] = useState<number>(0)
  const [machineTypes, setMachineTypes] = useState<string>('')
  const [memory, setMemory] = useState<string>('')
  const [cpu, setCpu] = useState<string>('')
  const [keyTaint, setKeyTaint] = useState<string>('')
  const [valueTaint, setValueTaint] = useState<string>('')
  const [keyLabel, setKeyLabel] = useState<string>('')
  const [valueLabel, setValueLabel] = useState<string>('')
  const [keyProv, setKeyProv] = useState<string>('')
  const [valueProv, setValueProv] = useState<string>('')
  const [effect, setEffect] = useState<string>('')
  const effectOptions = [
    { value: 'No_Schedule', label: 'No schedule'},
    { value: 'Prefer_No_Schedule', label: 'Prefer no schedule'},
    { value: 'No_Execute', label: 'No execute'}
  ]
  const [taints, setTaints] = useState<TypeTaints[]>()
  const [labels, setLabels] = useState<{}[]>()
  const [providerTags, setProviderTags] = useState<{}[]>()
  const [regionOptions, setRegionOptions] = useState<TypeOption[]>([])
  const [k8sOptions, setK8sOptions] = useState<TypeSelectOptions>()
  const [sshKey, setSshKey] = useState<string>('')
  const [subnets, SetSubnets] = useState<TypeSubnet[]>([])
  const [subnetWorkers, SetSubnetWorkers] = useState<string>('')
  const providerOptions = [{ value: 'aws', label: 'aws' }]
  const [replicasWorkers, setReplicasWorkers] = useState<number>(0)
  const [memoryWorkers, setMemoryWorkers] = useState<string>('')
  const [cpuWorkers, setCpuWorkers] = useState<string>('')
  const [machineTypesWorkers, setMachineTypesWorkers] = useState<string>('')
  const [keyTaintWorkers, setKeyTaintWorkers] = useState<string>('')
  const [valueTaintWorkers, setValueTaintWorkers] = useState<string>('')
  const [keyLabelWorkers, setKeyLabelWorkers] = useState<string>('')
  const [valueLabelWorkers, setValueLabelWorkers] = useState<string>('')
  const [keyProvWorkers, setKeyProvWorkers] = useState<string>('')
  const [valueProvWorkers, setValueProvWorkers] = useState<string>('')
  const [effectWorkers, setEffectWorkers] = useState<string>('')
  const [taintsWorkers, setTaintsWorkers] = useState<TypeTaints[]>()
  const [autoScale, setAutoScale] = useState<boolean>(false)
  const [maxSize, setMaxSize] = useState<number>(0)
  const [minSize, setMinSize] = useState<number>(0)
  const [groupId, setGroupId] = useState<string>('')
  const [providerTagsWorkers, setProviderTagsWorkers] = useState<{}[]>()
  const [labelsWorkers, setLabelsWorkers] = useState<{}[]>()
  const groupIdOptions = [{ value: '', label: ''}]
  const [groups, setGroups] = useState<TypeWorkers[]>()
  const [flavorOptions, setFlavorOptions] = useState<TypeOption[]>([])
  const [cpuOptions, setCpuOptions] = useState<TypeOption[]>()
  const [memOptions, setMemOptions] = useState<TypeOption[]>()
  const [MachineOptions, setMachineOptions] = useState<TypeOption[]>()
  const [session, setSession] = useState<string>('')
  // const handleAction = () => {
  //   handleClose()
  //   if (body.handleAction) body.handleAction()
  // }

  // const showModal = () => {
  //   handleClose()
  //   Modals.show('create-cluster', {
  //     title: 'Create',
  // 		ndTitle: 'Cluster'
  //   })
  // }

  const handleAction = () => {
    const data = {
      "apiVersion": "app.undistro.io/v1alpha1",
      "kind": "Cluster",
      "metadata": {
        "name": clusterName,
        "namespace": namespace
      },
      "spec": {
        "kubernetesVersion": k8sVersion,
        "controlPlane": {
          "internalLB": internalLB,
          "replicas": replicas,
          "machineType": machineTypes,
          "subnet": subnetWorkers,
          "labels": labels,
          "providerTags": providerTags,
          "taints": taints
        },
        "workers": groups,
        "bastion": {
          "enabled": enabled,
          "ingress": ingress,
          "instanceType": machineTypes,
          "allowedCIDRBlocks": cidrs
        },
        "infrastructureProvider": {
          "name": provider,
          "sshKey": sshKey,
          "flavor": flavor,
          "region": region
        },
        "network": {
          "apiServerPort": serverPort,
          "pods": podsRanges,
          "serviceDomain": serviceDomain,
          "multiZone": isPublic,
          "vpc": {
            "id": id,
            "cidrBlock": cidr,
            "zone": zone
          },

          "subnets": subnets
        }
      }
    }

    const dataPolicies = {
      "apiVersion": "app.undistro.io/v1alpha1",
      "kind": "DefaultPolicies",
      "metadata": {
        "name": "defaultpolicies-undistro",
        "namespace": namespace
      },
      "spec": {
        "clusterName": clusterName
      }
    }

    Api.Cluster.post(data, namespace)
      .then(_ => console.log('success'))
    
    setTimeout(() => {
      Api.Cluster.postPolicies(dataPolicies, namespace)
    }, 600)
  }

  const saveGroup = () => {
    setGroups([...(groups || []), {
      replicas: replicasWorkers,
      machineType: machineTypesWorkers,
      subnet: subnetWorkers,
      labels: labelsWorkers!,
      providerTags: providerTagsWorkers!,
      taints: taintsWorkers!,
      autoscaling: {
        enabled: autoScale,
        minSize: minSize,
        maxSize: maxSize
      }
    }])
  }


  const createCidrs = () => {
    setCidrs([...(cidrs || []), cidrBastion])
  }

  const deleteCidrs = () => {
    const newCidrs = (cidrs || []).slice(1, 1)
    setCidrs(newCidrs)
  }

  const createSubnets = () => {
    SetSubnets([...subnets, {
      id: idSubnet,
      cidrBlock: cidrSubnet,
      zone: zone,
      isPublic: isPublic
    }])
  }

  const createTaints = (onChange: Function, data: TypeTaints[]) => {
    onChange([...data, {
      key: keyTaint,
      value: valueTaint,
      effect: effect
    }])
  }

  const createMap = (onChange: Function, data: {}[], key: string, value: string) => {
    onChange([...data, { key: value }])
  }

  const deleteSubnets = (id: any) => {
    const newSubnets = subnets.filter(item => item.id !== id)
    SetSubnets(newSubnets)
  }

  const getSecrets = (secretRef: string) => {
    Api.Secret.list(secretRef).then(res => {
      setAccesskey(atob(res.data.accessKeyID))
      setSecret(atob(res.data.secretAccessKey))
      setRegion(atob(res.data.region))
    })
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
    getFlavors()
    getProviders()
    getMachines()
  }, [])

  console.log(subnets)


  return (
    <>
    <header>
      <h3 className="title"><span>{body.title}</span> {body.ndTitle}</h3>
      <i onClick={handleClose} className="icon-close" />
    </header>
      <div className='box'>
        <Steps handleClose={handleClose} handleAction={() => handleAction()}>
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

          <InfraNetwork 
            id={id}
            setId={setId}
            idSubnet={idSubnet}
            setIdSubnet={setIdSubnet}
            isPublic={isPublic}
            setIsPublic={setIsPublic}
            zone={zone}
            setZone={setZone}
            cidr={cidr}
            setCidr={setCidr}
            cidrSubnet={cidrSubnet}
            setCidrSubnet={setCidrSubnet}
            addSubnet={createSubnets}
            subnets={subnets}
            deleteSubnet={deleteSubnets}
          />

          <K8sNetwork 
            serverPort={serverPort}
            setServerPort={setServerPort}
            serviceDomain={serviceDomain}
            setServiceDomain={setServiceDomain}
            podsRanges={podsRanges}
            setPodsRanges={setPodsRanges}
            serviceRanges={serviceRanges}
            setServiceRanges={setServiceRanges}
          />

          <Bastion
            enabled={enabled}
            setEnabled={setEnabled}
            ingress={ingress}
            setIngress={setIngress}
            replicas={replicas}
            setReplicas={setReplicas}
            cpu={cpu}
            setCpu={setCpu}
            getCpu={cpuOptions}
            getMem={memOptions}
            memory={memory}
            setMemory={setMemory}
            machineTypes={machineTypes}
            setMachineTypes={setMachineTypes}
            getMachineTypes={MachineOptions}
            cidr={cidrBastion}
            setCidr={setCidrBastion}
            cidrs={(cidrs || [])}
            handleEvent={() => createCidrs()}
            deleteCidr={() => deleteCidrs()}
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
            clusterName={clusterName}
            getMachineTypes={MachineOptions}
            isAdvanced
            keyTaint={keyTaint}
            setKeyTaint={setKeyTaint}
            valueTaint={valueTaint}
            setValueTaint={setValueTaint}
            effectValue={effect}
            setEffectValue={setEffect}
            effect={effectOptions}
            keyLabel={keyLabel}
            setKeyLabel={setKeyLabel}
            valueLabel={valueLabel}
            setValueLabel={setValueLabel}
            keyProv={keyProv}
            setKeyProv={setKeyProv}
            valueProv={valueProv}
            setValueProv={setValueProv}
            taints={taints}
            providers={providerTags}
            labels={labels}
            internalLB={internalLB}
            setInternalLB={setInternalLB}
            handleActionTaints={() => createTaints(setTaints, (taints || []))}
            handleActionLabel={() => createMap(setLabels, (labels || []), keyLabel, valueLabel)}
            handleActionProv={() => createMap(setProviderTags, (providerTags || []), keyProv, valueProv)}
          />

          <Workers
            subnet={subnetWorkers}
            setSubnet={SetSubnetWorkers}
            handleAction={saveGroup}
            groupIdOptions={groupIdOptions}
            groupId={groupId}
            setGroupId={setGroupId}
            replicas={replicasWorkers}
            setReplicas={setReplicasWorkers}
            cpu={cpuWorkers}
            setCpu={setCpuWorkers}
            getCpu={cpuOptions}
            getMem={memOptions}
            memory={memoryWorkers}
            setMemory={setMemoryWorkers}
            machineTypes={machineTypesWorkers}
            setMachineTypes={setMachineTypesWorkers}
            getMachineTypes={MachineOptions}
            autoScale={autoScale}
            setAutoScale={setAutoScale}
            maxSize={maxSize}
            setMaxSize={setMaxSize}
            minSize={minSize}
            setMinSize={setMinSize}
            keyTaint={keyTaintWorkers}
            setKeyTaint={setKeyTaintWorkers}
            valueTaint={valueTaintWorkers}
            setValueTaint={setValueTaintWorkers}
            effectValue={effectWorkers}
            setEffectValue={setEffectWorkers}
            effect={effectOptions}
            keyLabel={keyLabelWorkers}
            setKeyLabel={setKeyLabelWorkers}
            valueLabel={valueLabelWorkers}
            setValueLabel={setValueLabelWorkers}
            keyProv={keyProvWorkers}
            setKeyProv={setKeyProvWorkers}
            valueProv={valueProvWorkers}
            setValueProv={setValueProvWorkers}
            taints={taintsWorkers}
            providers={providerTagsWorkers}
            labels={labelsWorkers}
            handleActionTaints={() => createTaints(setTaintsWorkers, (taintsWorkers || []))}
            handleActionLabel={() => createMap(setLabelsWorkers, (labelsWorkers || []), keyLabelWorkers, valueLabelWorkers)}
            handleActionProv={() => createMap(setProviderTagsWorkers, (providerTagsWorkers || []), keyProvWorkers, valueProvWorkers)}
          />
        </Steps>
      </div>
    </>
  )
}

export default ClusterAdvanced
