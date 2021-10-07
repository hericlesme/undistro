/* eslint-disable react-hooks/exhaustive-deps */
import { Layout } from '@components/layout'
import { InfranodeDetails } from '@components/details'
import { useEffect, useState } from 'react'
import { useServices } from 'providers/ServicesProvider'

export default function ControlPlanePage() {
  const [data, setData] = useState<any>()
  const { Api } = useServices()

  const getData = () => {
    Api.Cluster.get('undistro-system', 'wizard').then(elm => setData(elm))
  }
  useEffect(() => {
    getData()
  }, [])

  return data ? (
    <Layout>
      <div className="home-page-route">
        <InfranodeDetails
          data={{
            controlPlaneMachineType: data.spec.controlPlane.machineType,
            controlPlaneReplicas: data.spec.controlPlane.replicas,
            controlPlaneSubnet: data.spec.controlPlane.subnet,
            internalLb: data.spec.controlPlane.internalLB
          }}
          onCancel={() => alert('#TODO: Voltar para página de listagem!')}
          onSave={data => {
            console.log(data)
          }}
        />
      </div>
    </Layout>
  ) : null
}
