/* eslint-disable react-hooks/exhaustive-deps */
import { InfranodeDetails } from '@components/details'
import { useEffect, useState } from 'react'
import Api from 'util/api'

export default function ControlPlanePage() {
  const [data, setData] = useState<any>()

  const getData = () => {
    Api.Cluster.get('undistro-system', 'wizard')
      .then(elm => setData(elm))
  }
  useEffect(() => {
    getData()
  }, [])

  return data?(
    <div className="home-page-route">
      <InfranodeDetails
        data={{
          controlPlaneMachineType: data.spec.controlPlane.machineType,
          controlPlaneReplicas: data.spec.controlPlane.replicas,
          controlPlaneSubnet: data.spec.controlPlane.subnet,
          internalLb: data.spec.controlPlane.internalLB,
        }}
        onCancel={() => alert('#TODO: Voltar para página de listagem!')}
        onSave={data => {
          console.log(data)
        }}
      />
    </div>
  ) : null
}
