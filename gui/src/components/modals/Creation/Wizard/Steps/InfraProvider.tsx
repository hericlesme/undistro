import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import classNames from 'classnames'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import { TextInput } from '@/components/forms/TextInput'
import { useFetch } from '@/hooks/query'
import { Flavor } from '@/types/cluster'
import { useWatch } from 'react-hook-form'
import { Select } from '@/components/forms/Select'

const InfraProvider: VFC<FormActions> = ({ register, control }: FormActions) => {
  const { data: flavors } = useFetch<Flavor[]>('/api/metadata/flavors')

  const selectedFlavor = useWatch({
    control: control,
    name: 'infraProviderFlavor'
  })

  const getFlavorOptions = () => {
    if (!flavors) return []
    return flavors.map(flavor => flavor.name)
  }

  const getK8SVersionOptions = () => {
    if (!selectedFlavor) return []
    return flavors.find(f => f.name === selectedFlavor).supportedVersions
  }

  return (
    <>
      <div className={styles.inputRow}>
        <TextInput type="text" label="ID" placeholder="ID" fieldName="infraProviderID" register={register} />
        <TextInput
          type="text"
          label="CIDR block"
          placeholder="CIDR block"
          fieldName="infraProviderCIDR"
          register={register}
        />
      </div>
      <div className={styles.inputRow}>
        <Select
          label="Flavor"
          fieldName="infraProviderFlavor"
          placeholder="Select flavor"
          register={register}
          options={getFlavorOptions()}
        />
        <Select
          label="Kubernetes version"
          fieldName="infraProviderK8sVersion"
          placeholder="Select K8s"
          register={register}
          options={getK8SVersionOptions()}
        />
      </div>
      <div className={styles.inputBlockTextArea}>
        <label className={styles.createClusterLabel} htmlFor="infraProviderSshKey">
          sshKey
        </label>
        <br />
        <textarea
          className={classNames(styles.createClusterTextInput, styles.input100, styles.textAreaInput)}
          placeholder="ssh Key"
          id="infraProviderSshKey"
          name="infraProviderSshKey"
          {...register('infraProviderSshKey', { required: true })}
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>

      <div className={styles.inputBlock}>
        <label className={styles.createClusterLabel} htmlFor="infraProviderUploadFile">
          upload config file
        </label>
        <div className={styles.inputUploadFileBlock}>
          <input
            className={classNames(styles.createClusterTextInput, styles.input100)}
            placeholder="clouds.yaml"
            type="text"
            id="infraProviderUploadFile"
            name="infraProviderUploadFile"
            {...register('infraProviderUploadFile', { required: false })}
          />
          <div className={styles.uploadFileButtonContainer}>
            <button className={styles.uploadFileIcon}></button>
          </div>
        </div>
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>
    </>
  )
}

export { InfraProvider }
