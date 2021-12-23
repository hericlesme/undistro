import type { MouseEventHandler, VFC } from 'react'
import type { Worker as IWorker } from '@/types/cluster'
import type { FormActions } from '@/types/utils'

import { useState, useEffect } from 'react'
import classNames from 'classnames'

import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'

type WorkerProps = {
  worker: IWorker
  index: number
  removeWorker: (index: number) => MouseEventHandler<HTMLButtonElement>
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

const ControlPlane: VFC<FormActions> = ({ register, setValue }: FormActions) => {
  const [workers, setWorkers] = useState([])

  const [workerConfig, setWorkerConfig] = useState({
    workersInfraNodeSwitch: false,
    workersReplicas: '',
    workersCPU: '',
    workersMem: '',
    workersMachineType: ''
  })

  useEffect(() => {
    console.log(workers)
    setValue('workers', workers)
  }, [workers])

  const addWorker = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()

    setWorkers([...workers, { name: `worker-config-${Math.random() * 100}`, ...workerConfig }])
  }

  const removeWorker = (index: number) => (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    console.log(index)
    console.log(workers)
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
      <div className={styles.controlPlaneInputRow}>
        <div className={styles.inputBlockSmall}>
          <label className={styles.createClusterLabel} htmlFor="controlPlaneReplicas">
            replicas
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="controlPlaneReplicas"
            name="controlPlaneReplicas"
            {...register('controlPlaneReplicas')}
          >
            <option value="" disabled selected hidden>
              n#
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlockSmall}>
          <label className={styles.createClusterLabel} htmlFor="controlPlaneCPU">
            CPU
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="controlPlaneCPU"
            name="controlPlaneCPU"
            {...register('controlPlaneCPU')}
          >
            <option value="" disabled selected hidden>
              CPU
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlockSmall}>
          <label className={styles.createClusterLabel} htmlFor="controlPlaneMem">
            mem
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="controlPlaneMem"
            name="controlPlaneMem"
            {...register('controlPlaneMem')}
          >
            <option value="" disabled selected hidden>
              mem
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="controlPlaneMachineType">
            machineType
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="controlPlaneMachineType"
            name="controlPlaneMachineType"
            {...register('controlPlaneMachineType')}
          >
            <option value="" disabled selected hidden>
              machine type
            </option>
            <option value="option1">option1</option>
            <option value="option2">option2</option>
            <option value="option3">option3</option>
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
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
            <div className={styles.inputBlockSmall}>
              <label className={styles.createClusterLabel} htmlFor="workersReplicas">
                replicas
              </label>

              <select
                className={classNames(styles.createClusterTextSelect, styles.input100)}
                id="workersReplicas"
                name="workersReplicas"
                onChange={handleWorkerConfigChange}
              >
                <option value="" disabled selected hidden>
                  n#
                </option>
                <option value="option1">option1</option>
                <option value="option2">option2</option>
                <option value="option3">option3</option>
              </select>
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
            <div className={styles.inputBlockSmall}>
              <label className={styles.createClusterLabel} htmlFor="workersCPU">
                CPU
              </label>

              <select
                className={classNames(styles.createClusterTextSelect, styles.input100)}
                id="workersCPU"
                name="workersCPU"
                onChange={handleWorkerConfigChange}
              >
                <option value="" disabled selected hidden>
                  CPU
                </option>
                <option value="option1">option1</option>
                <option value="option2">option2</option>
                <option value="option3">option3</option>
              </select>
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
            <div className={styles.inputBlockSmall}>
              <label className={styles.createClusterLabel} htmlFor="workersMem">
                mem
              </label>

              <select
                className={classNames(styles.createClusterTextSelect, styles.input100)}
                id="workersMem"
                name="workersMem"
                onChange={handleWorkerConfigChange}
              >
                <option value="" disabled selected hidden>
                  mem
                </option>
                <option value="option1">option1</option>
                <option value="option2">option2</option>
                <option value="option3">option3</option>
              </select>
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
            <div className={styles.inputBlock}>
              <label className={styles.createClusterLabel} htmlFor="workersMachineType">
                machineType
              </label>

              <select
                className={classNames(styles.createClusterTextSelect, styles.input100)}
                id="workersMachineType"
                name="workersMachineType"
                onChange={handleWorkerConfigChange}
              >
                <option value="" disabled selected hidden>
                  machine type
                </option>
                <option value="option1">option1</option>
                <option value="option2">option2</option>
                <option value="option3">option3</option>
              </select>
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>
          </div>

          <div className={classNames(styles.addButtonBlock, styles.justifyRight)}>
            <button onClick={addWorker} className={styles.solidMdButtonDefault}>
              add
            </button>
          </div>

          <div className={styles.modalTableContainer}>
            <table className={styles.modalWorkersTable} id="wizardWorkersTable">
              <tbody>
                {workers.map((worker, index) => (
                  <Worker index={index} worker={worker} removeWorker={removeWorker} />
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
