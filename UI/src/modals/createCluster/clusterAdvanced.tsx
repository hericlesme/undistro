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
import { LoadOptions } from 'react-select-async-paginate'
import Steps from './steps'
import Api from 'util/api'
import { TypeOption, TypeSelectOptions, TypeSubnet, TypeTaints } from '../../types/cluster'
import { TypeModal } from '../../types/generic'

type TypeWorkers = {
  replicas: number
  machineType: TypeOption | null
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
  const [machineTypes, setMachineTypes] = useState<TypeOption | null>(null)
  const [memory, setMemory] = useState<TypeOption | null>(null)
  const [cpu, setCpu] = useState<TypeOption | null>(null)
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
  const [flavorOptions, setFlavorOptions] = useState<TypeOption[]>([])
  const [k8sOptions, setK8sOptions] = useState<TypeSelectOptions>()
  const [sshKey, setSshKey] = useState<string>('')
  const [subnets, SetSubnets] = useState<TypeSubnet[]>([])
  const [subnetWorkers, SetSubnetWorkers] = useState<string>('')
  const [sshKeyOptions, setSshKeyOptions] = useState<string[]>([])
  const providerOptions = [{ value: 'aws', label: 'aws' }]
  const [replicasWorkers, setReplicasWorkers] = useState<number>(0)
  const [memoryWorkers, setMemoryWorkers] = useState<TypeOption | null>(null)
  const [cpuWorkers, setCpuWorkers] = useState<TypeOption | null>(null)
  const [machineTypesWorkers, setMachineTypesWorkers] = useState<TypeOption | null>(null)
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

  const getSecrets = () => {
    Api.Secret.list()
      .then(res => {
        setAccesskey(atob(res.data.accessKeyID))
        setSecret(atob(res.data.secretAccessKey))
        setRegion(atob(res.data.region))
      })
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

  const getMachineTypes: LoadOptions<TypeOption, { page: number }> = async (value, loadedOptions, additional: any) => {
    const res = await Api.Provider.listMetadata('aws', 'machineTypes', '15', additional ? additional.page : 1, region)
    const totalPages = res.TotalPages
    const machineTypes = res.MachineTypes.map((elm: Partial<{ instanceType: string }>) => ({ value: elm.instanceType, label: elm.instanceType }))
    return {
      options: machineTypes,
      hasMore: additional && totalPages > additional.page,
      additional: { page: additional ? additional.page + 1 : 1 }
    }
  }

  const getCpu: LoadOptions<TypeOption, { page: number }> = async (value, loadedOptions, additional: any) => {
    const res = await Api.Provider.listMetadata('aws', 'machineTypes', '15', additional ? additional.page : 1, region)
    const totalPages = res.TotalPages
    const cpu = res.MachineTypes.map((elm: Partial<{ vcpus: string }>) => ({ value: elm.vcpus, label: elm.vcpus }))
    return {
      options: cpu,
      hasMore: additional && totalPages > additional.page,
      additional: { page: additional ? additional.page + 1 : 1 }
    }
  }

  const getMem: LoadOptions<TypeOption, { page: number }> = async (value, loadedOptions, additional: any) => {
    const res = await Api.Provider.listMetadata('aws', 'machineTypes', '15', additional ? additional.page : 1, region)
    const totalPages = res.TotalPages
    const cpu = res.MachineTypes.map((elm: Partial<{ memory: string }>) => ({ value: elm.memory, label: elm.memory }))
    return {
      options: cpu,
      hasMore: additional && totalPages > additional.page,
      additional: { page: additional ? additional.page + 1 : 1 }
    }
  }

  const getRegion = async () => {
    const res = await Api.Provider.listMetadata('aws', 'regions', '24', 1, region)

    setRegionOptions(res.map((elm: any) => ({ value: elm, label: elm })))
  }

  const getFlavors = async () => {
    const res = await Api.Provider.listMetadata('aws', 'supportedFlavors', '1', 1, region)
    type apiOption = {
      name: string;
      kubernetesVersion: string[];
    };

    type apiResponse = apiOption[];

    const parse = (data: apiResponse): TypeSelectOptions => {
      return data.reduce<TypeSelectOptions>((acc, curr) => {
        acc[curr.name] = {
          selectOptions: curr.kubernetesVersion.map((ver) => ({
            label: ver,
            value: ver,
          })),
        };

        return acc;
      }, {});
    };

    const parseData = parse(res)
    setFlavorOptions(Object.keys(parseData).map(elm => ({ value: elm, label: elm })))
    setK8sOptions(parseData)
  }

  const getKeys = async () => {
    const res = await Api.Provider.listMetadata('aws', 'sshKeys', '1', 1, region)
    setSshKeyOptions(res.map((elm: string) => ({ value: elm, label: elm })))
  }

  useEffect(() => {
    getSecrets()
    getRegion()
    getFlavors()
    getKeys()
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
            sshKeyOptions={sshKeyOptions}
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
            getCpu={getCpu}
            getMem={getMem}
            memory={memory}
            setMemory={setMemory}
            machineTypes={machineTypes}
            setMachineTypes={setMachineTypes}
            getMachineTypes={getMachineTypes}
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
            getCpu={getCpu}
            memory={memory}
            setMemory={setMemory}
            getMem={getMem}
            machineTypes={machineTypes}
            setMachineTypes={setMachineTypes}
            clusterName={clusterName}
            getMachineTypes={getMachineTypes}
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
            getCpu={getCpu}
            getMem={getMem}
            memory={memoryWorkers}
            setMemory={setMemoryWorkers}
            machineTypes={machineTypesWorkers}
            setMachineTypes={setMachineTypesWorkers}
            getMachineTypes={getMachineTypes}
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
