/* eslint-disable react-hooks/exhaustive-deps */
import React, { FC, useEffect, useState } from 'react'
import Input from '@components/input'
import Toggle from '@components/toggle'
import Select from '@components/select'
import { TypeBastion, TypeOption } from '../../../types/cluster'
import { useServices } from 'providers/ServicesProvider'

import './index.scss'

const Bastion: FC<TypeBastion> = ({
  enabled,
  setEnabled,
  ingress,
  setIngress,
  cpu,
  setCpu,
  getCpu,
  getMem,
  memory,
  setMemory,
  machineTypes,
  setMachineTypes,
  getMachineTypes,
  cidr,
  setCidr,
  cidrs,
  handleEvent,
  deleteCidr
}) => {
  const { Api } = useServices()
  const [ip, setIp] = useState<string>('')
  const formCidr = (e: React.FormEvent<HTMLInputElement>) => {
    setCidr(e.currentTarget.value)
  }

  const formCpu = (option: TypeOption | null) => {
    setCpu(option)
  }

  const formMem = (option: TypeOption | null) => {
    setMemory(option)
  }

  const formMachineTypes = (option: TypeOption | null) => {
    setMachineTypes(option)
  }

  useEffect(() => {
    Api.Provider.getUserIp().then(res => {
      setIp(res.ip)
    })
  }, [])

  return (
    <>
      <h3 className="title-box">Bastion</h3>
      <div className="bastion">
        <Toggle label="enabled" value={enabled} onChange={() => setEnabled(!enabled)} />
        <Toggle label="disable ingress rules" value={ingress} onChange={() => setIngress(!ingress)} />
        <div className="flex-text">
          <p>user default blocks CIDR:</p>
          <span>{ip}/32</span>
        </div>

        <div className="input-container">
          <Select value={cpu} onChange={formCpu} options={getCpu} label="CPU" />
          <Select value={memory} onChange={formMem} options={getMem} label="mem" />
          <Select value={machineTypes} onChange={formMachineTypes} options={getMachineTypes} label="machineType" />
        </div>

        <div className="cidrs-container">
          <Input
            type="text"
            addButton
            handleEvent={handleEvent}
            label="allowed blocks CIDR"
            value={cidr}
            onChange={formCidr}
          />

          <ul>
            {cidrs.map((elm, i) => {
              return (
                <li>
                  <p>allowedBlock-{i}</p>
                  <i onClick={() => deleteCidr?.()} className="icon-close" />
                </li>
              )
            })}
          </ul>
        </div>
      </div>
    </>
  )
}

export default Bastion
