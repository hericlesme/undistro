import type { ReactNode } from 'react'
import type { GetServerSideProps, NextPage } from 'next'

import { ClustersOverview } from '@/components/overviews/Clusters'
import { Workspace } from '@/components/Workspace'

type HomePageProps = {
  selectedClusters?: string[]
  children?: ReactNode
  page?: string
}

export const getServerSideProps: GetServerSideProps<HomePageProps> = async context => {
  const { query } = context

  let selectedClusters: string[] = []
  const { cluster, page } = query

  let pageVar = '1'
  if (page != undefined) {
    pageVar = page as string
  }

  if (cluster != undefined) {
    if (!(cluster instanceof Array)) {
      selectedClusters = [cluster]
    } else {
      selectedClusters = cluster
    }
  }

  return {
    props: {
      selectedClusters: selectedClusters,
      page: pageVar
    }
  }
}

const Home: NextPage<HomePageProps> = ({ selectedClusters, page }: HomePageProps) => {
  return (
    <Workspace selectedClusters={selectedClusters || []}>
      <ClustersOverview page={page!} />
    </Workspace>
  )
}

export default Home
