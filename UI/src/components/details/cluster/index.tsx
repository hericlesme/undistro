import { useEffect, useState } from 'react'
import { uid } from 'uid'
import Button from '@components/button'
import Input from '@components/input'
import Select from '@components/select'
import Toggle from '@components/toggle'
import { Grid, List } from '../common'
import { DataCell, EditButton } from '../common'
import '../index.scss'
 
type Block = {
  id: string
  cidr: string
  public: boolean
  zone: string
}
 
type Data = {
  // General.
  generalAccessKeyId: string
  generalClusterName: string
  generalDefaultRegion: string
  generalNamespace: string
  generalProvider: string
  generalSecretAccessKey: string
  generalSessionToken?: string
 
  // Infrastructure Provider.
  infraFlavor: string
  infraK8sVersion: string
  infraProvider: string
  infraRegion: string
  infraSshKey?: string
 
  // Infra Network - VPC.
  infraNetworkCidrBlock: string
  infraNetworkId: string
  infraNetworkZone: string
 
  // Bastion.
  bastionDisableIngressRules: boolean
  bastionEnabled: boolean
  bastionUserDefaultBlocksCidr?: string
 
  // Kubernetes Network.
  k8sApiServerPort?: string
  k8sNetworkMultiZone: boolean
  k8sPodsRanges?: string
  k8sServiceDomain?: string
  k8sServiceRanges?: string
}
 
export type FormData = {
  bastionDisableIngressRules: boolean;
  bastionEnabled: boolean;
  blocks: Block[]
  k8sNetworkMultiZone: boolean;
}
 
type Props = {
  data: Data
  onCancel: () => void
  onSave: (data: FormData) => void
}
 
const XIcon = () => {
  return (
    <svg fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
      <path
        clipRule="evenodd"
        d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
        fillRule="evenodd"
      />
    </svg>
  )
}
 
export const ClusterDetails = ({ data, onCancel: handleCancel, onSave: handleSave }: Props) => {
  // Permissões para editar dados.
  const [canEditBastionDisableIngressRules, setCanEditBastionDisableIngressRules] = useState(false)
  const [canEditBastionEnabled, setCanEditBastionEnabled] = useState(false)
  const [canEditK8sNetworkMultiZone, setCanEditK8sNetworkMultiZone] = useState(false)
 
  // Dados editáveis.
  const [bastionDisableIngressRules, setBastionDisableIngressRules] = useState(data.bastionDisableIngressRules)
  const [bastionEnabled, setBastionEnabled] = useState(data.bastionEnabled)
  const [k8sNetworkMultiZone, setK8sNetworkMultiZone] = useState(data.k8sNetworkMultiZone)
 
  // Formulário.
  const [subnetId, setSubnetId] = useState('')
  const [subnetIsPublic, setSubnetIsPublic] = useState(false)
  const [subnetZone, setSubnetZone] = useState('')
  const [subnetCidrBlock, setSubnetCidrBlock] = useState('')
 
  // Opções.
  const [subnetIdOptions, setSubnetIdOptions] = useState<string[]>([])
 
  // Listagens.
  const [blocks, setBlocks] = useState<Block[]>([])
 
  useEffect(() => {
    // Fetch options.
    setSubnetIdOptions([])
 
    // Fetch lists.
    setBlocks([])
  }, [])
 
  const addBlock = () => {
    setBlocks([...blocks, { id: uid(), cidr: subnetCidrBlock, public: subnetIsPublic, zone: subnetZone }])
    setSubnetId('')
    setSubnetIsPublic(false)
    setSubnetZone('')
    setSubnetCidrBlock('')
  }
 
  return (
    <div className="details-container">
      <h2 className="details-heading">General</h2>
      <DataCell label="Cluster Name">{data.generalClusterName}</DataCell>
      <DataCell label="Namespace">{data.generalNamespace}</DataCell>
      <DataCell label="Provider">{data.generalProvider}</DataCell>
      <DataCell label="Default Region">{data.generalDefaultRegion}</DataCell>
      <DataCell label="Access Key ID">{data.generalAccessKeyId}</DataCell>
      <DataCell label="Secret Access Key ID">{data.generalSecretAccessKey}</DataCell>
      {data.generalSessionToken && <DataCell label="Session Token">{data.generalSessionToken}</DataCell>}
 
      <h2 className="details-heading">Infrastructure Provider</h2>
      <DataCell label="Provider">{data.infraProvider}</DataCell>
      <DataCell label="Region">{data.infraRegion}</DataCell>
      <DataCell label="Kubernetes Version">{data.infraK8sVersion}</DataCell>
      <DataCell label="Flavor">{data.infraFlavor}</DataCell>
      {data.infraSshKey && <DataCell label="sshKey">{data.infraSshKey}</DataCell>}
 
      <h2 className="details-heading">Infra Network - VPC</h2>
      <DataCell label="ID">{data.infraNetworkId}</DataCell>
      <DataCell label="CIDR Block">{data.infraNetworkCidrBlock}</DataCell>
      <div className="bordered-container">
        <div className="center heading">Subnet</div>
        <div className="inner-container">
          <Toggle
            label="is public"
            value={subnetIsPublic}
            onChange={() => {
              setSubnetIsPublic(!subnetIsPublic)
            }}
          />
          <Grid columns="3">
            <Select
              label="ID"
              options={subnetIdOptions}
              placeholder="default"
              value={subnetId}
              onChange={(id: string) => setSubnetId(id)}
            />
            <Input
              label="zone"
              placeholder="zone"
              value={subnetZone}
              onChange={e => setSubnetZone(e.currentTarget.value)}
            />
            <Input
              label="CIDR block"
              placeholder="optional"
              value={subnetCidrBlock}
              onChange={e => setSubnetCidrBlock(e.currentTarget.value)}
            />
          </Grid>
          <Button disabled={!subnetZone} size="small" style={{ marginLeft: 'auto' }} variant="gray" onClick={addBlock}>
            Add
          </Button>
          <List>
            {blocks.map(block => (
              <li key={block.id}>
                <span>
                  {block.zone}-{block.cidr}
                </span>
                <button
                  onClick={() => {
                    alert('#TODO: Delete block!')
                    console.table(block)
 
                    // #TODO: Aproveitar trecho de código abaixo para atualizar a listagem após completar a requisição de exclusão.
                    //.then(() => setBlocks(blocks.filter(({ id }) => id === block.id)))
                  }}
                >
                  <XIcon />
                </button>
              </li>
            ))}
          </List>
        </div>
      </div>
 
      <h2 className="details-heading">Kubernetes Network</h2>
      {data.k8sApiServerPort && <DataCell label="API Server Port">{data.k8sApiServerPort}</DataCell>}
      {data.k8sServiceDomain && <DataCell label="Service Domain">{data.k8sServiceDomain}</DataCell>}
      {data.k8sPodsRanges && <DataCell label="Pods Ranges">{data.k8sPodsRanges}</DataCell>}
      {data.k8sServiceRanges && <DataCell label="Service Ranges">{data.k8sServiceRanges}</DataCell>}
      <DataCell label="Multi-Zone" borderless>
        <Toggle
          value={k8sNetworkMultiZone}
          onChange={() => {
            if (canEditK8sNetworkMultiZone) setK8sNetworkMultiZone(!k8sNetworkMultiZone)
          }}
        />
        <EditButton
          editing={canEditK8sNetworkMultiZone}
          onClick={() => setCanEditK8sNetworkMultiZone(!canEditK8sNetworkMultiZone)}
        />
      </DataCell>
 
      <h2 className="details-heading">Bastion</h2>
      <DataCell label="Enabled" borderless>
        <Toggle
          value={bastionEnabled}
          onChange={() => {
            if (canEditBastionEnabled) setBastionEnabled(!bastionEnabled)
          }}
        />
        <EditButton editing={canEditBastionEnabled} onClick={() => setCanEditBastionEnabled(!canEditBastionEnabled)} />
      </DataCell>
      <DataCell label="Disable Ingress Rules" borderless>
        <Toggle
          value={bastionDisableIngressRules}
          onChange={() => {
            if (canEditBastionDisableIngressRules) setBastionDisableIngressRules(!bastionDisableIngressRules)
          }}
        />
        <EditButton
          editing={canEditBastionDisableIngressRules}
          onClick={() => setCanEditBastionDisableIngressRules(!canEditBastionDisableIngressRules)}
        />
      </DataCell>
      {data.bastionUserDefaultBlocksCidr && <DataCell label="User Default Blocks CIDR">{data.bastionUserDefaultBlocksCidr}</DataCell>}
 
      <footer className="details-footer">
        <Button children="Cancel" size="large" variant="black" onClick={handleCancel} />
        <Button
          children="Save Changes"
          size="large"
          variant="primary"
          onClick={() =>
            handleSave({
              bastionDisableIngressRules,
              bastionEnabled,
              blocks,
              k8sNetworkMultiZone
            })
          }
        />
      </footer>
    </div>
  )
}