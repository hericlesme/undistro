import { useEffect, useState } from 'react'
import { uid } from 'uid'
import Button from '@components/button'
import Input from '@components/input'
import Select from '@components/select'
import Toggle from '@components/toggle'
import { Grid, List } from '../common'

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

export type Data = {
  // #TODO: Os outros dados devem ser adicionados aqui. Além disso, é provável que esses nomes precisem ser mudados.
  controlPlaneCpu?: string
  controlPlaneMachineType: string
  controlPlaneMem?: string
  controlPlaneReplicas: number
  controlPlaneSubnet?: string
  internalLb: boolean
}

type FormData = Omit<
  Data,
  'labelsKey' | 'labelsValue' | 'providerTagsKey' | 'providerTagsValue' | 'taintsEffect' | 'taintsKey' | 'taintsValue'
> & {
  labels: Label[]
  providerTags: Label[]
  taints: Taint[]
}

type Props = {
  data: Data
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

export const InfranodeDetails = ({ data, onCancel: handleCancel, onSave: handleSave }: Props) => {
  // Opções.
  const [controlPlaneCpuOptions, setControlPlaneCpuOptions] = useState<string[]>([])
  const [controlPlaneMachineTypeOptions, setControlPlaneMachineTypeOptions] = useState<string[]>([])
  const [controlPlaneMemOptions, setControlPlaneMemOptions] = useState<string[]>([])
  const [taintsEffectOptions, setTaintsEffectOptions] = useState<string[]>([])

  // Dados editáveis.
  const [controlPlaneCpu, setControlPlaneCpu] = useState(data.controlPlaneCpu)
  const [controlPlaneMachineType, setControlPlaneMachineType] = useState(data.controlPlaneMachineType)
  const [controlPlaneMem, setControlPlaneMem] = useState(data.controlPlaneMem)
  const [controlPlaneReplicas, setControlPlaneReplicas] = useState(data.controlPlaneReplicas)
  const [controlPlaneSubnet, setControlPlaneSubnet] = useState(data.controlPlaneSubnet)
  const [internalLb, setInternalLb] = useState(data.internalLb)

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

  useEffect(() => {
    // Fetch options.
    setControlPlaneCpuOptions([])
    setControlPlaneMachineTypeOptions([])
    setControlPlaneMemOptions([])
    setTaintsEffectOptions([])

    // Fetch lists.
    setLabels([])
    setProviderTags([])
    setTaints([])
  }, [])

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

  return (
    <div className="details-container">
      <h2 className="details-heading">Control Plane</h2>
      <Toggle label="internal LB" value={internalLb} onChange={() => setInternalLb(!internalLb)} />
      <Grid columns="8">
        <Input
          className="col-span-1"
          label="replicas"
          placeholder="n#"
          value={controlPlaneReplicas}
          onChange={e => {
            if (!e.currentTarget.value) {
              setControlPlaneReplicas(1)

              return
            }

            const value = e.currentTarget.value.replace(/\D/g, '')

            setControlPlaneReplicas(+value) // Conversão para number
          }}
        />
        {controlPlaneSubnet && <Input
          className="col-span-2"
          type='text'
          label="subnet"
          placeholder="subnet"
          value={controlPlaneSubnet}
          onChange={(e) => setControlPlaneSubnet(e.currentTarget.value)}
        />}
        <Select
          className="col-span-1"
          label="CPU"
          options={controlPlaneCpuOptions}
          placeholder="CPU"
          value={controlPlaneCpu}
          onChange={(cpu: string) => setControlPlaneCpu(cpu)}
        />
        <Select
          className="col-span-1"
          label="mem"
          options={controlPlaneMemOptions}
          placeholder="Mem"
          value={controlPlaneMem}
          onChange={(mem: string) => setControlPlaneMem(mem)}
        />
        <Select
          className="col-span-3"
          label="machineType"
          options={controlPlaneMachineTypeOptions}
          placeholder="Machine type"
          value={controlPlaneMachineType}
          onChange={(machineType: string) => setControlPlaneMachineType(machineType)}
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
            <button onClick={() => setTaints(taints.filter(({ id }) => id === taint.id))} >
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
            <button onClick={() => setLabels(labels.filter(({ id }) => id === label.id))}>
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
            <button onClick={() => setProviderTags(providerTags.filter(({ id }) => id === providerTag.id))} >
              <XIcon />
            </button>
          </li>
        ))}
      </List>

      <footer className="details-footer">
        <Button children="Cancel" size="large" variant="black" onClick={handleCancel} />
        <Button
          children="Save Changes"
          size="large"
          variant="primary"
          onClick={() => {
            handleSave({
              controlPlaneCpu,
              controlPlaneMachineType,
              controlPlaneMem,
              controlPlaneReplicas,
              controlPlaneSubnet,
              internalLb,
              labels,
              providerTags,
              taints
            })

            alert('#TODO: Voltar para listagem')
          }}
        />
      </footer>
    </div>
  )
}
