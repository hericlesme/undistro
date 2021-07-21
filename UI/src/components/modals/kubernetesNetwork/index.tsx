import React, { FC } from 'react'
import Input from '@components/input'
import { TypeK8sNetwork } from '../../../types/cluster'
import './index.scss'

const k8sNetwork: FC<TypeK8sNetwork> = ({
  serverPort,
  setServerPort,
  serviceDomain,
  setServiceDomain,
  podsRanges,
  setPodsRanges,
  serviceRanges,
  setServiceRanges
}) => {

  const formServer = (e: React.FormEvent<HTMLInputElement>) => {
    setServerPort(e.currentTarget.value)
  }

  const formServiceDomain = (e: React.FormEvent<HTMLInputElement>) => {
    setServiceDomain(e.currentTarget.value)
  }

  const formPods = (e: React.FormEvent<HTMLInputElement>) => {
    setPodsRanges(e.currentTarget.value)
  }

  const formServiceRanges = (e: React.FormEvent<HTMLInputElement>) => {
    setServiceRanges(e.currentTarget.value)
  }

  return (
    <>
    <h3 className="title-box">Kubernetes network</h3>
      <div className='kubernetes-network'>
        <div className='single-input'>
          <Input type='text' label='API server port' value={serverPort} onChange={formServer} />
        </div>
        <Input type='text' label='service domain' value={serviceDomain} onChange={formServiceDomain} />
        <div className='input-flex'>
          <Input type='text' label='pods ranges' value={podsRanges} onChange={formPods} />
          <Input type='text' label='service ranges' value={serviceRanges} onChange={formServiceRanges} />
        </div>
      </div>
    </>
  )
}

export default k8sNetwork
