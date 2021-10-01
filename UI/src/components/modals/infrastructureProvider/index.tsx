import React, { FC } from 'react'
import Select from '@components/select'
import Input from '@components/input'
import { TypeInfra } from '../../../types/cluster'

const InfrastructureProvider: FC<TypeInfra> = ({
  flavor,
  setFlavor,
  flavorOptions,
  k8sVersion,
  setK8sVersion,
  k8sOptions,
  sshKey,
  setSshKey,
  cidr,
  setCidr,
  id,
  setId,
}) => {

  const formId = (e: React.FormEvent<HTMLInputElement>) => {
    setId?.(e.currentTarget.value)
  }

  const formCidr = (e: React.FormEvent<HTMLInputElement>) => {
    setCidr?.(e.currentTarget.value)
  }

  const formSshKey = (e: React.FormEvent<HTMLInputElement>) => {
    setSshKey(e.currentTarget.value)
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
        <Input type='text' value={id} onChange={formId} label='ID' />
        <Input type='text' label='CIDR block' value={cidr} onChange={formCidr} />
        <Select
          value={flavor}
          onChange={(v: string) => {
            formK8s('')
            formFlavor(v)
          }}
          options={flavorOptions}
          label="flavor"
        />
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
