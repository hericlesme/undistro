import type { VFC } from 'react'

import Link from 'next/link'
import classNames from 'classnames'

import styles from '@/components/ContentNotFound/ContentNotFound.module.css'

const ContentNotFound: VFC = () => {
  const contentNotFoundStyles = {
    mainTextFirstLine: classNames(styles.ContentNotFoundMainTextLine1, 'upperCase'),
    mainTextSecondLine: classNames(styles.ContentNotFoundMainTextLine2, 'upperCase'),
    secondaryTextFirstLine: classNames(styles.ContentNotFoundSecondaryTextLine1, 'upperCase'),
    secondaryTextSecondLine: classNames(styles.ContentNotFoundSecondaryTextLine2, 'upperCase')
  }

  return (
    <div className={styles.ContentNotFoundContainer}>
      <div className={styles.ContentNotFoundMonitorMessage}></div>
      <div className={contentNotFoundStyles.mainTextFirstLine}>it seems that one of our</div>
      <div className={contentNotFoundStyles.mainTextSecondLine}>trainees screwed up again..</div>
      <div className={contentNotFoundStyles.secondaryTextFirstLine}>
        you can go to the{' '}
        <Link href="/">
          <a>home page</a>
        </Link>{' '}
        while
      </div>
      <div className={contentNotFoundStyles.secondaryTextSecondLine}>we look for someone to blame</div>
    </div>
  )
}

export { ContentNotFound }
