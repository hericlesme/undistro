import type { ReactNode } from 'react'
import type { GetServerSideProps, NextPage } from 'next'

import { Workspace } from '@/components/Workspace'
import { NodePoolsOverview } from '@/components/overviews/MachinePools'

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
      <NodePoolsOverview page={page!} />
    </Workspace>
  )
}

export default Home
