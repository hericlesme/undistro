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

  const currCPU = useWatch({
    control,
    name: 'controlPlaneCPU'
  })

  const [workers, setWorkers] = useState([])
  const { data: machineTypes } = useFetch<MachineType[]>('/api/metadata/machinetypes')

  const getMachineTypeOptions = () => {
    if (!machineTypes) return []
    return machineTypes.map(machineType => machineType.name)
  }

  const sortUniques = (arr: (string | number)[]): (string | number)[] => {
    const sortAlphaNum = (a, b) => a.localeCompare(b, 'en', { numeric: true })
    return Array.from(new Set(arr)).sort(sortAlphaNum)
  }

  const getMachineAttr = (attr: MACHINE_ATTR) => {
    if (!machineTypes) return []
    return sortUniques(machineTypes.map(machineType => machineType[attr]))
  }

  const workerDefaults = {
    workersInfraNodeSwitch: false,
    workersReplicas: 3,
    workersMachineType: ''
  }

  const [workerConfig, setWorkerConfig] = useState(workerDefaults)

  useEffect(() => {
    console.log(workers)
    let workersData = workers.map(w => ({
      infraNode: w.workersInfraNodeSwitch ? true : false,
      machineType: w.workersMachineType,
      replicas: w.workersReplicas
    }))

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
    console.log(event.target.id)
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
          options={getMachineAttr(MACHINE_ATTR.NAME)}
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
              options={getMachineAttr(MACHINE_ATTR.NAME)}
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
