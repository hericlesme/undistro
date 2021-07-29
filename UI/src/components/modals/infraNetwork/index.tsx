import React, { FC } from 'react'
import Input from '@components/input'
import Button from '@components/button'
import Toggle from '@components/toggle'
import { TypeNetwork } from '../../../types/cluster'
import './index.scss'

const InfraNetwork: FC<TypeNetwork> = ({
  id,
  setId,
  idSubnet,
  setIdSubnet,
  isPublic,
  setIsPublic,
  zone,
  setZone,
  cidr,
  setCidr,
  cidrSubnet,
  setCidrSubnet,
  zoneSubnet,
  setZoneSubnet,
  addSubnet,
  subnets,
  deleteSubnet
}) => {

  const formId = (e: React.FormEvent<HTMLInputElement>) => {
    setId(e.currentTarget.value)
  }
  
  const formIdSubnet = (e: React.FormEvent<HTMLInputElement>) => {
    setIdSubnet(e.currentTarget.value)
  }

  const formZone = (e: React.FormEvent<HTMLInputElement>) => {
    setZone(e.currentTarget.value)
  }

  const formZoneSubnet = (e: React.FormEvent<HTMLInputElement>) => {
    setZoneSubnet(e.currentTarget.value)
  }
  
  const formCidr = (e: React.FormEvent<HTMLInputElement>) => {
    setCidr(e.currentTarget.value)
  }

  const formCidrSubnet = (e: React.FormEvent<HTMLInputElement>) => {
    setCidrSubnet(e.currentTarget.value)
  }

  return (
    <>
      <h3 className="title-box">Infra network - VPC</h3>
      <div className='infra-network'>
        <div className='input-container'>
          <Input type='text' value={id} onChange={formId} label='ID' />
          <Input type='text' label='zone' value={zone} onChange={formZone} />
          <Input type='text' label='CIDR block' value={cidr} onChange={formCidr} />
        </div>

        <div className='subnet'>
          <h3 className="title-box">Subnet</h3>
          
          <Toggle label='Is public' value={isPublic} onChange={() => setIsPublic(!isPublic)} />
          <div className='subnet-inputs'>
            <Input type='text' value={idSubnet} onChange={formIdSubnet} label='ID' />
            <Input type='text' label='zone' value={zoneSubnet} onChange={formZoneSubnet} />
            <Input type='text' label='CIDR block' value={cidrSubnet} onChange={formCidrSubnet} />
            <div className='button-container'>
              <Button onClick={() => addSubnet()} type='gray' size='small' children='Add' />
            </div>
          </div>

          <ul>
            {(subnets || []).map((elm, i = 0) => {
              return (
                <li key={elm.id}>
                  <p>subnet-{i}</p>
                  <i className='icon-close' onClick={() => deleteSubnet(elm.id)} />
                </li>
              )
            })}
          </ul>
        </div>
      </div>
    </>
  )
}

export default InfraNetwork
