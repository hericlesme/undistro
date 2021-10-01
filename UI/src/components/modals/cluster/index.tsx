import React, { FC } from 'react'
import Input from '@components/input'
import Select from '@components/select'
import { TypeCluster } from '../../../types/cluster'

import './index.scss'

const Cluster: FC<TypeCluster> = ({
  clusterName,
  setClusterName,
  namespace,
  setNamespace,
  provider,
  setProvider,
  providerOptions,
  regionOptions,
  region,
  setRegion
}) => {

  const formCluster = (e: React.FormEvent<HTMLInputElement>) => {
    setClusterName(e.currentTarget.value)
  }

  const formNamespace = (e: React.FormEvent<HTMLInputElement>) => {
    setNamespace(e.currentTarget.value)
  }

  const formProvider = (value: string) => {
    setProvider(value)
  }

  const formRegion = (value: string) => {
    setRegion(value)
  }

  return (
    <>
      <h3 className="title-box">Cluster</h3>
      <form className='create-cluster'>
        <Input value={clusterName} onChange={formCluster} type='text' label='cluster name' />
        <Input value={namespace} onChange={formNamespace} type='text' label='namespace' />
        <div className='select-flex'>
          <Select value={provider} onChange={formProvider} options={providerOptions} label='select provider' />
          <Select options={regionOptions} value={region} onChange={formRegion} label='default region' />
        </div>
      </form>
    </>
  )
}

export default Cluster