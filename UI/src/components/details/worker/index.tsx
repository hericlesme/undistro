import { useEffect, useState } from 'react'
import { uid } from 'uid'
import Button from '@components/button'
import Input from '@components/input'
import Select from '@components/select'
import Toggle from '@components/toggle'
import { Grid, List } from '../common'
import { TypeOption } from '../../../types/cluster'
import Api from 'util/api'

import '../index.scss'

type Label = {
  id: string
  key: string
  value: string
}

type ProviderTag = {
  id: string
  key: string
  value: string
}

type Taint = {
  id: string
  effect: string
  key: string
  value: string
}

type Group = {
  // #TODO: Os outros dados devem ser adicionados aqui. Além disso, é provável que esses nomes precisem ser mudados.
  infraNode: boolean
  maxSize?: number
  minSize?: number
  name?: string
  workerCpu?: string
  workerMachineType: string
  workerMem?: string
  workerReplicas: number
  workerSubnet?: string
}

type FormData = Omit<
  Group,
  'labelsKey' | 'labelsValue' | 'providerTagsKey' | 'providerTagsValue' | 'taintsEffect' | 'taintsKey' | 'taintsValue'
> & {
  labels: Label[]
  providerTags: Label[]
  taints: Taint[]
}

type Props = {
  groups: Group[]
  onCancel: () => void
  onSave: (data: FormData) => void
}

const XIcon = () => {
  return (
    <svg fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
      <path
        clipRule="evenodd"
        d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
        fillRule="evenodd"
      />
    </svg>
  )
}

export const WorkerDetails = ({ groups, onCancel: handleCancel, onSave: handleSave }: Props) => {
  const [groupName, setGroupName] = useState('')

  const data = groups.find(({ name }) => name === groupName) || ({} as Group)
  const groupNames = groups.map(({ name }) => ({ label: name, value: name }))

  // Opções.
  const [taintsEffectOptions, setTaintsEffectOptions] = useState<string[]>([])
  const [cpuOptions, setCpuOptions] = useState<TypeOption[]>()
  const [memOptions, setMemOptions] = useState<TypeOption[]>()
  const [MachineOptions, setMachineOptions] = useState<TypeOption[]>()

  // Dados editáveis.
  const [enableAutoScale, setEnableAutoScale] = useState(!!data.maxSize && !!data.minSize)
  const [infraNode, setInfraNode] = useState(data.infraNode || false)
  const [maxSize, setMaxSize] = useState(data.maxSize || 0)
  const [minSize, setMinSize] = useState(data.minSize || 0)

  const [workerCpu, setWorkerCpu] = useState(data.workerCpu || '')
  const [workerMachineType, setWorkerMachineType] = useState(data.workerMachineType || '')
  const [workerMem, setWorkerMem] = useState(data.workerMem || '')
  const [workerReplicas, setWorkerReplicas] = useState(data.workerReplicas || 0)
  const [workerSubnet, setWorkerSubnet] = useState(data.workerSubnet || '')

  // Formulário.
  const [labelsKey, setLabelsKey] = useState('')
  const [labelsValue, setLabelsValue] = useState('')
  const [providerTagsKey, setProviderTagsKey] = useState('')
  const [providerTagsValue, setProviderTagsValue] = useState('')
  const [taintsEffect, setTaintsEffect] = useState('')
  const [taintsKey, setTaintsKey] = useState('')
  const [taintsValue, setTaintsValue] = useState('')

  // Listagens.
  const [labels, setLabels] = useState<Label[]>([])
  const [providerTags, setProviderTags] = useState<ProviderTag[]>([])
  const [taints, setTaints] = useState<Taint[]>([])

  const getMachines = () => {
    Api.Provider.list('awsmachines').then(res => {
      const name = res.items.map((elm: any) => ({
        label: elm.metadata.name,
        value: elm.metadata.name
      }))
      const cpu = res.items.map((elm: any) => ({
        label: elm.spec.vcpus,
        value: elm.spec.vcpus
      }))
      const mem = res.items.map((elm: any) => ({
        label: elm.spec.memory,
        value: elm.spec.memory
      }))

      setMachineOptions(name)
      setMemOptions(mem)
      setCpuOptions(cpu)
    })
  }

  useEffect(() => {
    // Fetch options.
    setTaintsEffectOptions([])
    getMachines()
    // Fetch lists.
    setLabels([])
    setProviderTags([])
    setTaints([])
  }, [])

  useEffect(() => {
    setEnableAutoScale(!!data.maxSize && !!data.minSize)
    setInfraNode(data.infraNode || false)
    setMaxSize(data.maxSize || 0)
    setMinSize(data.minSize || 0)
    setWorkerCpu(data.workerCpu || '')
    setWorkerMachineType(data.workerMachineType || '')
    setWorkerMem(data.workerMem || '')
    setWorkerReplicas(data.workerReplicas || 0)
    setWorkerSubnet(data.workerSubnet || '')

    setLabelsKey('')
    setLabelsValue('')
    setProviderTagsKey('')
    setProviderTagsValue('')
    setTaintsEffect('')
    setTaintsKey('')
    setTaintsValue('')
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [groupName])

  const addLabel = () => {
    setLabels([...labels, { id: uid(), key: labelsKey, value: labelsValue }])
    setLabelsKey('')
    setLabelsValue('')
  }

  const addProviderTag = () => {
    setProviderTags([...providerTags, { id: uid(), key: providerTagsKey, value: providerTagsValue }])
    setProviderTagsKey('')
    setProviderTagsValue('')
  }

  const addTaint = () => {
    setTaints([...taints, { id: uid(), effect: taintsEffect, key: taintsKey, value: taintsValue }])
    setTaintsEffect('')
    setTaintsKey('')
    setTaintsValue('')
  }

  const clearForm = () => {
    setEnableAutoScale(false)
    setInfraNode(false)
    setLabelsKey('')
    setLabelsValue('')
    setMaxSize(0)
    setMinSize(0)
    setProviderTagsKey('')
    setProviderTagsValue('')
    setTaintsEffect('')
    setTaintsKey('')
    setTaintsValue('')
    setWorkerCpu('')
    setWorkerMachineType('')
    setWorkerMem('')
    setWorkerReplicas(0)
    setWorkerSubnet('')

    setGroupName('')
  }

  return (
    <div className="details-container">
      <div className="bordered-container">
        <Select
          className="center"
          label="group ID"
          options={groupNames}
          value={groupName}
          onChange={(groupName: string) => setGroupName(groupName)}
        />
        <h2 className="details-heading">Workers</h2>
        <Toggle label="infraNode" value={infraNode} onChange={() => setInfraNode(!infraNode)} />
        <Grid columns="8">
          <Input
            className="col-span-1"
            label="replicas"
            placeholder="n#"
            value={workerReplicas}
            onChange={e => {
              if (!e.currentTarget.value) {
                setWorkerReplicas(1)

                return
              }

              const value = e.currentTarget.value.replace(/\D/g, '')

              setWorkerReplicas(+value) // Conversão para number
            }}
          />
          <Input
            className="col-span-2"
            label="subnet"
            placeholder="subnet"
            value={workerSubnet}
            onChange={(e) => setWorkerSubnet(e.currentTarget.value)}
          />
          <Select
            className="col-span-1"
            label="CPU"
            options={cpuOptions}
            placeholder="CPU"
            value={workerCpu}
            onChange={(cpu: string) => setWorkerCpu(cpu)}
          />
          <Select
            className="col-span-1"
            label="mem"
            options={memOptions}
            placeholder="Mem"
            value={workerMem}
            onChange={(mem: string) => setWorkerMem(mem)}
          />
          <Select
            className="col-span-3"
            label="machineType"
            options={MachineOptions}
            placeholder="Machine type"
            value={workerMachineType}
            onChange={(machineType: string) => setWorkerMachineType(machineType)}
          />
        </Grid>
        <Toggle
          label="enable auto scale"
          value={enableAutoScale}
          onChange={() => {
            const enabled = !enableAutoScale

            if (!enabled) {
              setMaxSize(0)
              setMinSize(0)
            }

            setEnableAutoScale(enabled)
          }}
        />
        <Grid columns="12">
          <Input
            className="col-span-1"
            disabled={!enableAutoScale}
            label="min size"
            placeholder="n#"
            value={minSize}
            onChange={e => {
              if (!e.currentTarget.value) {
                setMinSize(1)

                return
              }

              const value = e.currentTarget.value.replace(/\D/g, '')

              setMinSize(+value) // Conversão para number
            }}
          />
          <Input
            className="col-span-1"
            disabled={!enableAutoScale}
            label="max size"
            placeholder="n#"
            value={maxSize}
            onChange={e => {
              if (!e.currentTarget.value) {
                setMaxSize(1)

                return
              }

              const value = e.currentTarget.value.replace(/\D/g, '')

              setMaxSize(+value) // Conversão para number
            }}
          />
        </Grid>

        <h2 className="details-heading">Taints</h2>
        <Grid columns="3">
          <Input label="key" placeholder="key" value={taintsKey} onChange={e => setTaintsKey(e.currentTarget.value)} />
          <Input
            label="value"
            placeholder="value"
            value={taintsValue}
            onChange={e => setTaintsValue(e.currentTarget.value)}
          />
          <Select
            label="taint effect"
            options={taintsEffectOptions}
            value={taintsEffect}
            onChange={(effect: string) => setTaintsEffect(effect)}
          />
        </Grid>
        <Button
          disabled={!taintsEffect && !taintsKey && !taintsValue}
          size="small"
          style={{ marginLeft: 'auto' }}
          variant="gray"
          onClick={addTaint}
        >
          Add
        </Button>
        <List>
          {taints.map(taint => (
            <li key={taint.id}>
              <span>
                {taint.key}-{taint.value}
              </span>
              <button
                onClick={() => {
                  alert('#TODO: Delete taint!')
                  console.table(taint)

                  // #TODO: Aproveitar trecho de código abaixo para atualizar a listagem após completar a requisição de exclusão.
                  //.then(() => setTaints(taints.filter(({ id }) => id === taint.id)))
                }}
              >
                <XIcon />
              </button>
            </li>
          ))}
        </List>

        <h2 className="details-heading">Labels</h2>
        <Grid columns="3">
          <Input label="key" placeholder="key" value={labelsKey} onChange={e => setLabelsKey(e.currentTarget.value)} />
          <Input
            label="value"
            placeholder="value"
            value={labelsValue}
            onChange={e => setLabelsValue(e.currentTarget.value)}
          />
        </Grid>
        <Grid columns="3" noMargin>
          <div />
          <Button
            disabled={!labelsKey && !labelsValue}
            size="small"
            style={{ marginLeft: 'auto' }}
            variant="gray"
            onClick={addLabel}
          >
            Add
          </Button>
        </Grid>
        <List>
          {labels.map(label => (
            <li key={label.id}>
              <span>
                {label.key}-{label.value}
              </span>
              <button
                onClick={() => {
                  alert('#TODO: Delete label!')
                  console.table(label)

                  // #TODO: Aproveitar trecho de código abaixo para atualizar a listagem após completar a requisição de exclusão.
                  //.then(() => setLabels(labels.filter(({ id }) => id === label.id)))
                }}
              >
                <XIcon />
              </button>
            </li>
          ))}
        </List>

        <h2 className="details-heading">Provider Tags</h2>
        <Grid columns="3">
          <Input
            label="key"
            placeholder="key"
            value={providerTagsKey}
            onChange={e => setProviderTagsKey(e.currentTarget.value)}
          />
          <Input
            label="value"
            placeholder="value"
            value={providerTagsValue}
            onChange={e => setProviderTagsValue(e.currentTarget.value)}
          />
        </Grid>
        <Grid columns="3" noMargin>
          <div />
          <Button
            disabled={!providerTagsKey && !providerTagsValue}
            size="small"
            style={{ marginLeft: 'auto' }}
            variant="gray"
            onClick={addProviderTag}
          >
            Add
          </Button>
        </Grid>
        <List>
          {providerTags.map(providerTag => (
            <li key={providerTag.id}>
              <span>
                {providerTag.key}-{providerTag.value}
              </span>
              <button
                onClick={() => {
                  alert('#TODO: Delete provider tag!')
                  console.table(providerTags)

                  // #TODO: Aproveitar trecho de código abaixo para atualizar a listagem após completar a requisição de exclusão.
                  //.then(() => setProviderTags(providerTags.filter(({ id }) => id === providerTag.id)))
                }}
              >
                <XIcon />
              </button>
            </li>
          ))}
        </List>
      </div>

      <footer className="details-footer">
        <Button children="Cancel" size="large" variant="black" onClick={handleCancel} />
        <Button
          children="Save Changes"
          size="large"
          variant="primary"
          onClick={() => {
            handleSave({
              infraNode,
              labels,
              maxSize,
              minSize,
              name: data.name, // Se não tiver um nome, será criação de novo grupo. Caso contrário, atualização.
              providerTags,
              taints,
              workerCpu,
              workerMachineType,
              workerMem,
              workerReplicas,
              workerSubnet
            })

            clearForm()
          }}
        />
      </footer>
    </div>
  )
}
