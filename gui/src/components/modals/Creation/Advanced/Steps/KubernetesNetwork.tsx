import type { VFC } from 'react'
import type { FormActions } from '@/types/utils'

import { useFetch } from '@/hooks/query'

import { Provider } from '@/types/cluster'
import { useWatch } from 'react-hook-form'
import { TextInput, Select } from '@/components/forms'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'
import classNames from 'classnames'

const KubernetesNetwork: VFC<FormActions> = ({ register }: FormActions) => {
  return (
    <>
      <div className={classNames(styles.inputBlock, styles.inputMedium)}>
        <label className={styles.createClusterLabel} htmlFor="k8sNetworkApiServerPort">
          API server port
        </label>
        <input
          className={classNames(styles.createClusterTextInput, styles.input100)}
          placeholder="port n#"
          type="text"
          id="k8sNetworkApiServerPort"
          name="k8sNetworkApiServerPort"
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>

      <div className={classNames(styles.inputBlock, styles.inputFit)}>
        <label className={styles.createClusterLabel} htmlFor="k8sNetworkServiceDomain">
          service domain
        </label>
        <input
          className={classNames(styles.createClusterTextInput, styles.input100)}
          placeholder="service domain"
          type="text"
          id="k8sNetworkServiceDomain"
          name="k8sNetworkServiceDomain"
        />
        <a className={styles.assistiveTextDefault}>Assistive text default color</a>
      </div>

      <div className={styles.inputRow}>
        <div className={classNames(styles.inputBlock, styles.inputFit)}>
          <label className={styles.createClusterLabel} htmlFor="k8sNetworkPodsRanges">
            pods ranges
          </label>
          <input
            className={classNames(styles.createClusterTextInput, styles.input100)}
            placeholder="000.000.000.000/0000"
            type="text"
            id="k8sNetworkPodsRanges"
            name="k8sNetworkPodsRanges"
          />
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
        <div className={classNames(styles.inputBlock, styles.inputFit)}>
          <label className={styles.createClusterLabel} htmlFor="k8sNetworkServiceRanges">
            service ranges
          </label>
          <input
            className={classNames(styles.createClusterTextInput, styles.input100)}
            placeholder="000.000.000.000/0000"
            type="text"
            id="k8sNetworkServiceRanges"
            name="k8sNetworkServiceRanges"
          />
          <a className={styles.assistiveTextDefault}>Assistive text default color</a>
        </div>
      </div>
      <div className={styles.inputRow}>
        <div className={classNames(styles.inputBlock, styles.inputFit)}>
          <div className={classNames(styles.switchContainer, styles.justifyLeft)}>
            <a className={styles.createClusterLabel}>multi-zone</a>
            <label className={styles.switch} htmlFor="k8sNetworkMultiZone">
              <input type="checkbox" id="k8sNetworkMultiZone" name="k8sNetworkMultiZone" />
              <span className={classNames(styles.slider, styles.round)}></span>
            </label>
          </div>
        </div>
      </div>
    </>
  )
}

export { KubernetesNetwork }
