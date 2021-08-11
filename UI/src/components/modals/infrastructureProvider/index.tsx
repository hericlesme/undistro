import React, { FC } from 'react'
import Select from '@components/select'
import Input from '@components/input'
import { TypeInfra } from '../../../types/cluster'

const InfrastructureProvider: FC<TypeInfra> = ({
  provider,
  setProvider,
  providerOptions,
  flavor,
  setFlavor,
  flavorOptions,
  k8sVersion,
  setK8sVersion,
  k8sOptions,
  regionOptions,
  region,
  setRegion,
  sshKey,
  setSshKey
}) => {
  const formProvider = (value: string) => {
    setProvider(value)
  }

  const formSshKey = (e: React.FormEvent<HTMLInputElement>) => {
    setSshKey(e.currentTarget.value)
  }

  const formRegion = (value: string) => {
    setRegion(value)
  }

  const formFlavor = (value: string) => {
    setFlavor(value)
  }

  const formK8s = (value: string) => {
    setK8sVersion(value)
  }

  return (
    <>
      <h3 className="title-box">Infrastructure provider</h3>
      <form className="infra-provider">
        <Select value={provider} onChange={formProvider} options={providerOptions} label="provider" />
        <Select
          value={flavor}
          onChange={(v: string) => {
            formK8s('')
            formFlavor(v)
          }}
          options={flavorOptions}
          label="flavor"
        />
        <Select options={regionOptions} value={region} onChange={formRegion} label="region" />
        <Select
          value={k8sVersion}
          onChange={formK8s}
          options={k8sOptions?.[flavor]?.selectOptions ?? []}
          label="kubernetes version"
        />
        <Input type="text" value={sshKey} onChange={formSshKey} label="sshKey" />
      </form>
    </>
  )
}

export default InfrastructureProvider
