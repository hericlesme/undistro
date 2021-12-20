import classNames from 'classnames'
import styles from '@/components/overviews/Clusters/Creation/ClusterCreation.module.css'

const Wizard = () => {
  return (
    <div className={styles.createClusterWizContainer}>
      <div className={styles.modalTitleContainer}>
        <a className={styles.modalCreateClusterTitle}>cluster</a>
      </div>

      <div className={styles.modalContentContainer}>
        <div className={styles.modalInputArea}>
          <form className={styles.modalForm} id="wizardClusterForm">
            <div className={styles.inputBlock}>
              <label className={styles.createClusterLabel} htmlFor="clusterName">
                Cluster name
              </label>
              <input
                className={classNames(styles.createClusterTextInput, 'input100')}
                placeholder="choose a cool name for this cluster"
                type="text"
                id="clusterName"
                name="clusterName"
              />
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>

            <div className={styles.inputBlock}>
              <label className={styles.createClusterLabel} htmlFor="clusterNamespace">
                Namespace
              </label>
              <input
                className={classNames(styles.createClusterTextInput, 'input100')}
                placeholder="namespace"
                type="text"
                id="clusterNamespace"
                name="clusterNamespace"
              />
              <a className={styles.assistiveTextDefault}>Assistive text default color</a>
            </div>

            <div className={styles.inputRow}>
              <div className={styles.inputBlock}>
                <label className={styles.createClusterLabel} htmlFor="clusterProvider">
                  Provider
                </label>
                <select
                  className={classNames(styles.createClusterTextSelect, 'input100')}
                  id="clusterProvider"
                  name="clusterProvider"
                >
                  <option value="" disabled selected hidden>
                    Select provider
                  </option>
                  <option value="option1">option1</option>
                  <option value="option2">option2</option>
                  <option value="option3">option3</option>
                </select>
                <a className={styles.assistiveTextDefault}>Assistive text default color</a>
              </div>
              <div className={styles.inputBlock}>
                <label className={styles.createClusterLabel} htmlFor="clusterDefaultRegion">
                  Default region
                </label>
                <select
                  className={classNames(styles.createClusterTextSelect, 'input100')}
                  id="clusterDefaultRegion"
                  name="clusterDefaultRegion"
                >
                  <option value="" disabled selected hidden>
                    Select region
                  </option>
                  <option value="option1">option1</option>
                  <option value="option2">option2</option>
                  <option value="option3">option3</option>
                </select>
                <a className={styles.assistiveTextDefault}>Assistive text default color</a>
              </div>
            </div>
          </form>
        </div>
      </div>

      <div className={styles.modalDialogButtonsContainer}>
        <div className={styles.leftButtonContainer}>
          <button className={styles.borderButtonDefault}>
            <a>back</a>
          </button>
        </div>
        <div className={styles.rightButtonContainer}>
          <button className={styles.borderButtonSuccess}>
            <a>next</a>
          </button>
        </div>
      </div>
    </div>
  )
}

export { Wizard }
