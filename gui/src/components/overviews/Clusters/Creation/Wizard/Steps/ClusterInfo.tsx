import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import { useFetch } from '@/hooks/query'

import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'
import { Provider } from '@/types/cluster'
import { useWatch } from 'react-hook-form'

const ClusterInfo: VFC<FormActions> = ({ register, control }: FormActions) => {
  const selectedProvider = useWatch({
    control: control,
    name: 'clusterProvider'
  })

  const { data: providers, isLoading } = useFetch<Provider[]>('/api/metadata/providers')

  const getRegionOptions = () => {
    const provider = selectedProvider && providers.find(p => p.metadata.name === selectedProvider)

    if (!provider || !provider.status.regionNames) return []

    return provider.status.regionNames.map(region => (
      <option key={region} value={region}>
        {region}
      </option>
    ))
  }

  return (
    <>
      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor="clusterName">
          Cluster name
        </label>
        <input
          className={classNames(styles.createClusterTextInput, styles.input100)}
          placeholder="choose a cool name for this cluster"
          type="text"
          id="clusterName"
          name="clusterName"
          {...register('clusterName', { required: true })}
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>
      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor="clusterNamespace">
          Namespace
        </label>
        <input
          className={classNames(styles.createClusterTextInput, styles.input100)}
          placeholder="namespace"
          type="text"
          id="clusterNamespace"
          name="clusterNamespace"
          {...register('clusterNamespace', { required: true })}
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>
      <div className={styles.inputRow}>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="clusterProvider">
            Provider
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="clusterProvider"
            name="clusterProvider"
            required
            {...register('clusterProvider', { required: true })}
          >
            <option value="" disabled selected hidden>
              Select provider
            </option>
            {!isLoading &&
              providers.map(provider => (
                <option key={provider.metadata.name} value={provider.metadata.name}>
                  {provider.metadata.name}
                </option>
              ))}
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={styles.inputBlock}>
          <label className={styles.createClusterLabel} htmlFor="clusterDefaultRegion">
            Default region
          </label>
          <select
            className={classNames(styles.createClusterTextSelect, styles.input100)}
            id="clusterDefaultRegion"
            name="clusterDefaultRegion"
            required
            {...register('clusterDefaultRegion', { required: true })}
          >
            <option value="" disabled selected hidden>
              Select region
            </option>
            {getRegionOptions()}
          </select>
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
    </>
  )
}

export { ClusterInfo }
