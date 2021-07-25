import React, { useEffect, useState } from 'react'
import Table from '@components/table'
import Api from 'util/api'
import moment from 'moment'

import './index.scss'

const headers = [
	{ name: 'Name', field: 'name'},
	{ name: 'Provider', field: 'provider'},
	{ name: 'Flavor', field: 'flavor'},
	{ name: 'Version', field: 'version'},
	{ name: 'Cluster Groups', field: 'clusterGroups'},
	{ name: 'Machines', field: 'machines'},
	{ name: 'Age', field: 'age'},
	{ name: 'Status', field: 'status'},
]

type TypeCluster = {
	name: string,
	provider: string,
	flavor: string,
	version: string,
	clusterGroups: string,
	machines: number,
	age: Date,
	status: string
}

export default function HomePage () {
	const [clusters, setClusters] = useState<TypeCluster[]>([])
	const [pause, setPause] = useState<boolean>(false)
	const name = (clusters || []).map(elm => elm.name).toString()
	const namespace = (clusters || []).map(elm => elm.clusterGroups).toString()

	moment.updateLocale('en', {
    relativeTime : {
        past:   "%s",
        s  : 's',
        ss : '%ds',
        m:  "m",
        mm: "%dm",
        h:  "h",
        hh: "%dh",
        d:  "d",
        dd: "%dd",
        M:  "m",
        MM: "%dm",
        y:  "y",
        yy: "%dy"
    }
})

	const getClusters = () => {
		Api.Cluster.list('undistro-system')
			.then((clusters) => {
				setClusters(clusters.items.map((elm: any) => {
					let message: string = elm.status.conditions[0].message.toLowerCase()
					let status = ''
					if (message.includes('wait cluster')) status = 'Provisioning'
					else if (message.includes('error')) status = 'Error'
					else if (message.includes('paused')) status = 'Paused'
					else status = 'Ready'

					return {
						name: elm.metadata.name,
						provider: elm.spec.infrastructureProvider.name,
						flavor: elm.spec.infrastructureProvider.flavor,
						version: elm.spec.kubernetesVersion,
						clusterGroups: elm.metadata.namespace,
						machines: elm.status.controlPlane.replicas + elm.status.totalWorkerReplicas,
						age: moment(elm.metadata.creationTimestamp).startOf('day').fromNow(),
						status: status
					}
				}))
			})
	}

	const pauseCluster = () => {
		setPause(!pause)	
		const payload = {
			"spec": {
				"paused": pause
			}
		}

		Api.Cluster.put(payload, namespace, name)
			.then(_ => {
				console.log('success')
			})
	}

	const deleteCluster = () => {
		Api.Cluster.delete(namespace, name)
			.then(_ => {
				console.log('success')
			})
	}


	useEffect(() => {
		getClusters()
	}, [])

	console.log(pause)

	return (
		<div className='home-page-route'>
			<Table data={(clusters || [])} icon={pause} delete={() => deleteCluster()} pause={() => pauseCluster()} header={headers}/>	
		</div>
	)
}