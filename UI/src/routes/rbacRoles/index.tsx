/* eslint-disable react-hooks/exhaustive-deps */
import Input from '@components/input'
import RbacTable from '@components/rbacTable'
import { Layout } from '@components/layout'
import { useServices } from 'providers/ServicesProvider'

import './index.scss'
import { useEffect, useState } from 'react'

const headers = [
  { name: 'User', field: 'user' },
  { name: 'Role', field: 'role' },
]

export default function RbacPage() {
  const [list, setList] = useState<string[]>([])
  const [array, setArray] = useState<any[]>([])
  const { Api } = useServices()

  const listUsers = () => {
    //boa sorte pra você que vai lidar com isso no futuro:
    const res: any = [];
    const cluster: any = [];

    Api.Rbac.listClusterRoleBindings()
      .then(elm => {
        elm.items.forEach((item: any) => {
          item.subjects?.forEach((subject: any) => {
            cluster.push({
              role: item.roleRef.name,
              name: subject.name,
              kind: subject.kind
            })
          })
        })
        const filteredCluster = cluster.filter((elm: any) => elm.kind === 'User').filter((i: any) => i.length !== 0)
        setArray(filteredCluster)
      })

    Api.Rbac.listRoleBindings()
      .then(elm => {
        elm.items.forEach((item: any) => {
          item.subjects.forEach((subject: any) => {
            res.push({
              role: item.roleRef.name,
              name: subject.name,
              kind: subject.kind
            })
          })
        })

        const filteredRes = res.filter((elm: any) => elm.kind === 'User').filter((i: any) => i.length !== 0)
        const arrays = filteredRes.concat(array)

        setList(arrays.map((elm: any) => {
          return {
            user: elm.name,
            role: elm.role
          }
        }))
      })
  }

  useEffect(() => {
    listUsers()
  }, [])
  return (
    <Layout>
      <div className='rbac-route'>
        <h3>Assign Role</h3>
        <div className='assign-role'>
          <Input placeholder='enter username' />
          <div className='add-role'>
            <i className='icon-add-circle' />
            <p>Add Role</p>
          </div>
        </div>

        <RbacTable data={list} header={headers} />
       </div>
    </Layout>
  )
}
