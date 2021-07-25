import React, { FC, useState } from 'react'
import Input from '@components/input'
import Button from '@components/button'
import Toggle from '@components/toggle'
import Select from '@components/select'
import FormSlider from '@components/formSlider'
import { TypeOption, TypeControlPlane } from '../../../types/cluster'
import './index.scss'

const ControlPlane: FC<TypeControlPlane> = ({
  replicas,
  setReplicas,
  cpu,
  setCpu,
  getCpu,
  getMem,
  memory,
  setMemory,
  machineTypes,
  setMachineTypes,
  getMachineTypes,
  infraNode,
  setInfraNode,
  createWorkers,
  deleteWorkers,
  workers,
  replicasWorkers,
  setReplicasWorkers,
  cpuWorkers,
  setCpuWorkers,
  memoryWorkers,
  setMemoryWorkers,
  machineTypesWorkers,
  setMachineTypesWorkers,
  clusterName,
  isAdvanced,
  keyTaint,
  setKeyTaint,
  valueTaint,
  setValueTaint,
  effectValue,
  setEffectValue,
  effect,
  keyLabel,
  setKeyLabel,
  valueLabel,
  setValueLabel,
  keyProv,
  setKeyProv,
  valueProv,
  setValueProv,
  taints,
  providers,
  labels,
  deleteTaints,
  deleteProviders,
  deleteLabels,
  handleActionTaints,
  handleActionLabel,
  handleActionProv,
  internalLB,
  setInternalLB
}) => {
  const [showTaint, setShowTaint] = useState<boolean>(false)
  const [showLabel, setShowLabel] = useState<boolean>(false)
  const [showProv, setShowProv] = useState<boolean>(false)

  const formReplica = (e: React.FormEvent<HTMLInputElement>) => {
    setReplicas(parseInt(e.currentTarget.value) || 0)
  }

  const formReplicaWorkers = (e: React.FormEvent<HTMLInputElement>) => {
    setReplicasWorkers?.(parseInt(e.currentTarget.value) || 0)
  }

  const formCpu = (option: TypeOption | null) => {
    setCpu(option)
  }

  const formMem = (option: TypeOption | null) => {
    setMemory(option)
  }

  const formMachineTypes = (option: TypeOption | null) => {
    setMachineTypes(option)
  }

  const formCpuWorkers = (option: TypeOption | null) => {
    setCpuWorkers?.(option)
  }

  const formMemWorkers = (option: TypeOption | null) => {
    setMemoryWorkers?.(option)
  }

  const formMachineTypesWorkers = (option: TypeOption | null) => {
    setMachineTypesWorkers?.(option)
  }


  return (
    <>
    <h3 className="title-box">Control plane</h3>
    <div className={!isAdvanced ? 'control-plane' : 'control-plane advanced'}>
      <div className='input-container'>
        <Input value={replicas} onChange={formReplica} type='text' label='replicas' />
        <Select value={cpu} onChange={formCpu} options={getCpu} label='CPU' />
        <Select value={memory} onChange={formMem} options={getMem} label='mem' />
        <Select value={machineTypes} onChange={formMachineTypes} options={getMachineTypes} label='machineType' />
      </div>

      {!isAdvanced ? 
        <div className='workers'>
          <h3 className="title-box">Workers</h3>
          <Toggle label='InfraNode' value={infraNode} onChange={() => setInfraNode?.(!infraNode)} />
          <div className='input-container'>
            <Input type='text' label='replicas' value={replicasWorkers} onChange={formReplicaWorkers} />
            <Select value={cpuWorkers} onChange={formCpuWorkers} options={getCpu} label='CPU' />
            <Select value={memoryWorkers} onChange={formMemWorkers} options={getMem} label='mem' />
            <Select value={machineTypesWorkers} onChange={formMachineTypesWorkers} options={getMachineTypes} label='machineType' />
            <div className='button-container'>
              <Button onClick={() => createWorkers?.()} type='gray' size='small' children='Add' />
            </div>
          </div>

          <ul>
            {(workers || []).map((elm, i = 0) => {
              return (
                <li key={elm.id}>
                  <p>{clusterName}-mp-{i}</p>
                  <i onClick={() => deleteWorkers?.(elm.id)} className='icon-close' />
                </li>
              )
            })}
          </ul>
        </div> :
        <>
          <div className='boxes-container'>
            <div className='box-content'>
              <p className='title'>Taints</p>
              <ul>
                {(taints || []).map((elm: any, i) => {
                  return (
                    <li key={elm}>
                      <p>taint-{i}</p>
                      <i onClick={() => deleteTaints?.(elm)} className='icon-close' />
                    </li>
                  )
                })}
              </ul>
              <i className='icon-plus' onClick={() => setShowTaint(!showTaint)}/>
              {showTaint && 
                <FormSlider
                  direction='left'
                  title='Add taints'
                  key={keyTaint!}
                  setKey={setKeyTaint!}
                  value={valueTaint!}
                  setValue={setValueTaint!}
                  taint={effectValue}
                  setTaint={setEffectValue}
                  options={effect}
                  select
                  handleAction={() => handleActionTaints?.()}
                  handleClose={() => setShowTaint(!showTaint)}
                />}
            </div>

            <div className='box-content'>
              <p className='title'>Labels</p>
              <ul>
                {(labels || []).map((elm: any, i) => {
                  return (
                    <li key={elm}>
                      <p>label-{i}</p>
                      <i onClick={() => deleteLabels?.(elm)} className='icon-close' />
                    </li>
                  )
                })}
              </ul>
              <i className='icon-plus' onClick={() => setShowLabel(!showLabel)}/>
              {showLabel && 
                <FormSlider
                  direction='right'
                  title='Add labels'
                  key={keyLabel!}
                  setKey={setKeyLabel!}
                  value={valueLabel!}
                  setValue={setValueLabel!}
                  handleAction={() => handleActionLabel?.()}
                  handleClose={() => setShowLabel(!showLabel)}
                />}
            </div>

            <div className='box-content'>
              <p className='title'>Provider tags</p>
              <ul>
              <ul>
                {(providers || []).map((elm: any, i) => {
                  return (
                    <li key={elm}>
                      <p>provTag-{i}</p>
                      <i onClick={() => deleteProviders?.(elm)} className='icon-close' />
                    </li>
                  )
                })}
              </ul>
              </ul>
              <i className='icon-plus' onClick={() => setShowProv(!showProv)}/>
              {showProv && 
                <FormSlider
                  direction='right'
                  title='Add provider tags'
                  key={keyProv!}
                  setKey={setKeyProv!}
                  value={valueProv!}
                  setValue={setValueProv!}
                  handleAction={() => handleActionProv?.()}
                  handleClose={() => setShowProv(!showProv)}
                />}
            </div>
          </div>
          <Toggle label='Internal LB' value={internalLB!} onChange={() => setInternalLB?.(!internalLB)} />
        </>}
    </div>
  </>
  )
}

export default ControlPlane