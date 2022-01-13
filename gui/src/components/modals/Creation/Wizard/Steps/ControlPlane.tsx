import type { MouseEventHandler, VFC } from 'react'
import type { MachineType, Worker as IWorker } from '@/types/cluster'
import type { FormActions } from '@/types/utils'

import { useState, useEffect } from 'react'
import { useWatch } from 'react-hook-form'
import classNames from 'classnames'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import { useFetch } from '@/hooks/query'
import { Select } from '@/components/forms/Select'
import { TextInput } from '@/components/forms'

type WorkerProps = {
  worker: IWorker
  index: number
  removeWorker: (index: number) => MouseEventHandler<HTMLButtonElement>
}

enum MACHINE_ATTR {
  MEM = 'mem',
  CPU = 'cpu',
  NAME = 'name'
}

const Worker: VFC<WorkerProps> = ({ worker, index, removeWorker }: WorkerProps) => (
  <tr>
    <td>
      <div className={styles.modalTableRow}>
        <div className={styles.modalTableItem}>
          <a>{worker.name}</a>
        </div>
        <div className={styles.modalDeleteTableItemContainer}>
          <button onClick={removeWorker(index)} className={styles.deleteTableItem}></button>
        </div>
      </div>
    </td>
  </tr>
)

const ControlPlane: VFC<FormActions> = ({ register, setValue, control }: FormActions) => {
  const clusterName = useWatch({
    control,
    name: 'clusterName'
  })

  const controlPlaneCPU = useWatch({
    control,
    name: 'controlPlaneCPU'
  })

  const controlPlaneMem = useWatch({
    control,
    name: 'controlPlaneMem'
  })

  const selectedMachineType = useWatch({
    control,
    name: 'controlPlaneMachineType'
  })

  useEffect(() => {
    if (machineTypes && selectedMachineType) {
      let selectedMachine = machineTypes.find(m => m.name === selectedMachineType)
      setValue('controlPlaneCPU', selectedMachine.cpu)
      setValue('controlPlaneMem', selectedMachine.mem)
    }
  }, [selectedMachineType])

  const [workers, setWorkers] = useState([])
  const { data: machineTypes } = useFetch<MachineType[]>('/api/metadata/machinetypes')

  const sortUniques = (arr: (string | number)[]): (string | number)[] => {
    const sortAlphaNum = (a, b) => a.localeCompare(b, 'en', { numeric: true })
    return Array.from(new Set(arr)).sort(sortAlphaNum)
  }

  const controlPlaneMachineTypeFilter = (machineType: string) => {
    let machineFilter = true
    let machine = machineTypes.find(m => m.name === machineType)
    if (controlPlaneCPU !== undefined) {
      machineFilter = machineFilter && machine.cpu == controlPlaneCPU
    }
    if (controlPlaneMem !== undefined) {
      machineFilter = machineFilter && machine.mem == controlPlaneMem
    }
    return machineFilter
  }

  const workersMachineTypeFilter = (machineType: string) => {
    let machineFilter = true
    let machine = machineTypes.find(m => m.name === machineType)
    if (workerConfig.workersCPU !== undefined) {
      machineFilter = machineFilter && machine.cpu == workerConfig.workersCPU
    }
    if (workerConfig.workersMem !== undefined) {
      machineFilter = machineFilter && machine.mem == workerConfig.workersMem
    }
    return machineFilter
  }

  const getMachineAttr = (attr: MACHINE_ATTR, filter = undefined) => {
    if (!machineTypes) return []
    let machineAttrs = sortUniques(machineTypes.map(machineType => machineType[attr]))
    if (filter) machineAttrs = machineAttrs.filter(e => filter(e))
    return machineAttrs
  }

  const workerDefaults = {
    workersInfraNodeSwitch: false,
    workersReplicas: 3,
    workersMachineType: '',
    workersCPU: undefined,
    workersMem: undefined
  }

  const [workerConfig, setWorkerConfig] = useState(workerDefaults)

  useEffect(() => {
    let workersData = workers.map(w => {
      return {
        infraNode: w.workersInfraNodeSwitch ? true : false,
        machineType: w.workersMachineType,
        replicas: w.workersReplicas
      }
    })

    setValue('workers', workersData)
  }, [workers])

  const addWorker = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()

    setWorkers([...workers, { name: `${clusterName}-mp-${workers.length}`, ...workerConfig }])
  }

  const removeWorker = (index: number) => (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    const newWorkers = [...workers]
    newWorkers.splice(index, 1)
    setWorkers(newWorkers)
  }

  const handleWorkerConfigChange = event => {
    setWorkerConfig(prevState => ({
      ...prevState,
      [event.target.id]: event.target.value
    }))
  }

  return (
    <>
      {/* Control Plane */}
      <div className={styles.controlPlaneInputRow}>
        {/* Replicas */}
        <TextInput
          type="number"
          inputSize="sm"
          label="Replicas"
          placeholder="Replicas"
          fieldName="controlPlaneReplicas"
          register={register}
          defaultValue={3}
          min={3}
        />
        {/* CPU */}
        <Select
          label="CPU"
          inputSize="sm"
          fieldName="controlPlaneCPU"
          placeholder="CPU"
          register={register}
          options={getMachineAttr(MACHINE_ATTR.CPU)}
        />
        {/* Memory */}
        <Select
          label="Mem"
          inputSize="sm"
          fieldName="controlPlaneMem"
          placeholder="Mem"
          register={register}
          options={getMachineAttr(MACHINE_ATTR.MEM)}
        />
        {/* Machine Type */}
        <Select
          label="MachineType"
          fieldName="controlPlaneMachineType"
          placeholder="machine type"
          register={register}
          options={getMachineAttr(MACHINE_ATTR.NAME, controlPlaneMachineTypeFilter)}
        />
      </div>

      <div className={styles.modalWorkersContainer}>
        <div className={styles.workersTitleContainer}>
          <a className={styles.modalCreateClusterTitle}>workers</a>
        </div>
        <div className={styles.modalWorkersBlock}>
          <div className={styles.infraNodeBlock}>
            <div className={classNames(styles.switchContainer, styles.justifyLeft)}>
              <a className={styles.createClusterLabel}>infraNode</a>
              <label className={styles.switch} htmlFor="workersInfraNodeSwitch">
                <input
                  type="checkbox"
                  id="workersInfraNodeSwitch"
                  name="workersInfraNodeSwitch"
                  onChange={handleWorkerConfigChange}
                />
                <span className={classNames(styles.slider, styles.round)}></span>
              </label>
            </div>
          </div>
          <div className={styles.workersInputRow}>
            <TextInput
              type="number"
              inputSize="sm"
              label="Replicas"
              placeholder="Replicas"
              fieldName="workersReplicas"
              defaultValue={3}
              min={3}
              onChange={handleWorkerConfigChange}
            />
            <Select
              label="CPU"
              inputSize="sm"
              fieldName="workersCPU"
              placeholder="CPU"
              options={getMachineAttr(MACHINE_ATTR.CPU)}
              onChange={handleWorkerConfigChange}
            />
            <Select
              label="Mem"
              inputSize="sm"
              fieldName="workersMem"
              placeholder="Mem"
              options={getMachineAttr(MACHINE_ATTR.MEM)}
              onChange={handleWorkerConfigChange}
            />
            <Select
              label="MachineType"
              fieldName="workersMachineType"
              placeholder="Machine type"
              options={getMachineAttr(MACHINE_ATTR.NAME, workersMachineTypeFilter)}
              onChange={handleWorkerConfigChange}
            />
          </div>

          <div className={classNames(styles.addButtonBlock, styles.justifyRight)}>
            <button onClick={addWorker} className={styles.solidMdButtonDefault}>
              Add
            </button>
          </div>

          <div className={styles.modalTableContainer}>
            <table className={styles.modalWorkersTable} id="wizardWorkersTable">
              <tbody>
                {workers.map((worker, index) => (
                  <Worker key={`worker-${index}`} index={index} worker={worker} removeWorker={removeWorker} />
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </>
  )
}

export { ControlPlane }
