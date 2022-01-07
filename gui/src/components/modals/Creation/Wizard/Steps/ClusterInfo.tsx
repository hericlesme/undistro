import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import { useFetch } from '@/hooks/query'

import { Provider } from '@/types/cluster'
import { useWatch } from 'react-hook-form'
import { TextInput, Select } from '@/components/forms'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'

const ClusterInfo: VFC<FormActions> = ({ register, control }: FormActions) => {
  const selectedProvider = useWatch({
    control: control,
    name: 'clusterProvider'
  })

  const { data: providers } = useFetch<Provider[]>('/api/metadata/providers')

  const getRegionOptions = () => {
    const provider = selectedProvider && providers.find(p => p.metadata.name === selectedProvider)
    if (!provider || !provider.status.regionNames) return []

    return Array.from(new Set(provider.status.regionNames))
  }

  const getProviderOptions = () => {
    if (!providers) return []
    return providers.map(provider => provider.metadata.name)
  }

  return (
    <>
      <TextInput
        type="text"
        label="Cluster name"
        placeholder="choose a cool name for this cluster"
        fieldName="clusterName"
        register={register}
      />
      <TextInput
        type="text"
        label="Namespace"
        placeholder="namespace"
        fieldName="clusterNamespace"
        register={register}
      />
      <div className={styles.inputRow}>
        <Select
          label="Provider"
          fieldName="clusterProvider"
          placeholder="Select provider"
          register={register}
          options={getProviderOptions()}
        />
        <Select
          label="Region"
          fieldName="clusterDefaultRegion"
          placeholder="Select region"
          register={register}
          options={getRegionOptions()}
        />
      </div>
    </>
  )
}

export { ClusterInfo }
