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
  setRegion,
  accessKey,
  setAccesskey,
  secret,
  setSecret,
  session,
  setSession
}) => {

  const formCluster = (e: React.FormEvent<HTMLInputElement>) => {
    setClusterName(e.currentTarget.value)
  }

  const formNamespace = (e: React.FormEvent<HTMLInputElement>) => {
    setNamespace(e.currentTarget.value)
  }

  const formAccessKey = (e: React.FormEvent<HTMLInputElement>) => {
    setAccesskey(e.currentTarget.value)
  }

  const formSecret = (e: React.FormEvent<HTMLInputElement>) => {
    setSecret(e.currentTarget.value)
  }

  const formProvider = (value: string) => {
    setProvider(value)
  }

  const formRegion = (value: string) => {
    setRegion(value)
  }

  const formSession = (e: React.FormEvent<HTMLInputElement>) => {
    setSession(e.currentTarget.value)
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
        <Input disabled type='text' label='secret access ID' value={accessKey} onChange={formAccessKey} />
        <Input disabled type='text' label='secret access key' value={secret} onChange={formSecret} />
        <Input disabled type='text' label='session token' value={session} onChange={formSession} />
      </form>
    </>
  )
}

export default Cluster