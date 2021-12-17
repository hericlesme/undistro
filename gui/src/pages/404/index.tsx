import Link from 'next/link'
import classNames from 'classnames'

import styles from './404.module.css'

const Page404 = () => (
  <div className={styles.page404Container}>
    <div className={styles.page404logoContainer}>
      <div className={styles.page404logo}></div>
    </div>
    <div className={styles.page404messageContainer}>
      <div className={styles.page404MonitorMessage}></div>
      <div className={classNames(styles.page404MainTextLine1, 'upperCase')}>it seems that one of our</div>
      <div className={classNames(styles.page404MainTextLine2, 'upperCase')}>trainees screwed up again...</div>
      <div className={classNames(styles.page404SecondaryTextLine1, 'upperCase')}>
        you can go to the{' '}
        <Link href="/">
          <a>home page</a>
        </Link>{' '}
        while
      </div>
      <div className={classNames(styles.page404SecondaryTextLine2, 'upperCase')}>we look for someone to blame</div>
    </div>
  </div>
)

export default Page404
