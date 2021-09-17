import { useState, ReactNode } from 'react'
import './index.scss'
import Classnames from 'classnames'
import { useHistory } from 'react-router'
import Modal from 'util/modals'
import { useClusters } from 'providers/ClustersProvider'
import { useDisclosure } from 'hooks/useDisclosure'
import { PauseClusterAlert } from '@components/clusterActionAlert'
import { useLocation } from 'react-router-dom'
import Api from 'util/api'

type TypeSubItem = {
  name: string
  link?: string
  handleAction?: () => void | Promise<void>
}

type TypeItem = {
  name: string
  icon: ReactNode
  subItens: TypeSubItem[]
  url: string
}

const showModal = () => {
  Modal.show('create-cluster', {
    title: 'Create',
    ndTitle: 'Cluster',
    width: '720',
    height: '600'
  })
}

const MenuSideBar = () => {
  const { clusters } = useClusters()
  const [show, setShow] = useState<any>(false)
  const history = useHistory<any>()
  const [isPauseClusterAlertOpen, closePauseClusterAlert, openPauseClusterAlert] = useDisclosure()
  const location = useLocation()

  const handlePause = async () => {
    await Promise.all(
      clusters.map(async cluster => {
        const payload = {
          spec: {
            paused: !cluster.paused
          }
        }

        await Api.Cluster.put(payload, cluster.namespace, cluster.name)
      })
    )

    closePauseClusterAlert()
    window.location.reload()
  }

  const SubItens: TypeSubItem[] = [
    { name: 'Create', handleAction: () => showModal() },
    { name: 'Pause', handleAction: openPauseClusterAlert },
    { name: 'Update', link: '/nodepools' },
    { name: 'Settings', link: '/nodepools' }
  ]

  const SubItensMachines: TypeSubItem[] = [
    { name: 'Nodepools', link: '/nodepools' },
    { name: 'Pause', link: '/nodepools' }
  ]

  const items: TypeItem[] = [
    {
      name: 'Cluster',
      icon: (
        <svg
          width="27"
          height="27"
          style={{ marginBottom: '4px', marginRight: '12px', marginLeft: '8px' }}
          viewBox="0 0 27 27"
          fill="currentColor"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M16.9328 23.9502C17.2666 24.6314 17.6002 25.3126 17.9337 25.9941C18.0557 25.9364 18.1769 25.8787 18.2981 25.821L18.299 25.8206C18.4205 25.7627 18.5419 25.7049 18.6642 25.6471C18.3309 24.9668 17.9975 24.2863 17.6639 23.6058C20.0532 22.3593 21.8628 20.1824 22.5923 17.5666C23.3346 17.7331 24.0768 17.8997 24.8189 18.0664C24.8491 17.9374 24.8793 17.8076 24.9094 17.6779L24.9105 17.6733C24.9403 17.5451 24.9702 17.4169 25 17.2894C24.2579 17.1228 23.5161 16.9563 22.7743 16.79C22.886 16.2041 22.9444 15.5997 22.9444 14.9819C22.9444 12.8279 22.234 10.8368 21.0309 9.22266C21.6226 8.759 22.2141 8.29525 22.8056 7.83134L22.2985 7.20851C21.7074 7.6721 21.1163 8.13554 20.525 8.59888C18.7838 6.6225 16.2482 5.34036 13.406 5.23185V3H12.594V5.23292C9.75682 5.34879 7.22707 6.63425 5.49109 8.61148C4.89446 8.14394 4.29792 7.67631 3.70155 7.20851L3.19444 7.83134C3.79176 8.30002 4.38892 8.76837 4.98598 9.23651C3.78909 10.8482 3.08276 12.834 3.08276 14.9819C3.08276 15.5986 3.14097 16.2018 3.25229 16.7867C2.50138 16.9551 1.75061 17.1237 1 17.2923C1.03018 17.4213 1.06037 17.5511 1.09055 17.6809C1.1207 17.8104 1.15096 17.9405 1.18111 18.0694C1.93226 17.9007 2.68325 17.7323 3.43414 17.5639C4.16108 20.1736 5.96314 22.3469 8.34346 23.5954C8.00746 24.2811 7.67155 24.967 7.33581 25.653C7.45798 25.7108 7.5794 25.7686 7.70082 25.8264C7.82232 25.8842 7.94403 25.9422 8.06628 26C8.40193 25.3135 8.73758 24.6273 9.07323 23.9413C10.2812 24.4549 11.6136 24.7395 13.0136 24.7395C14.4054 24.7395 15.7304 24.4582 16.9328 23.9502ZM15.3529 15.9429C16.9301 16.2973 18.508 16.651 20.0859 17.0046L20.0902 17.0056C20.2857 17.0494 20.4812 17.0932 20.6767 17.1371C20.084 19.1757 18.6713 20.8744 16.8116 21.8669C16.7301 21.7006 16.6486 21.5343 16.5671 21.3681L16.5576 21.3487C15.8549 19.9153 15.1523 18.4818 14.4504 17.0491C14.8609 16.7733 15.1718 16.3878 15.3529 15.9429ZM16.0803 22.2109C15.8595 21.7605 15.6386 21.3102 15.4178 20.8598C14.8516 19.7051 14.2855 18.5506 13.7199 17.3961C13.249 17.5355 12.748 17.5355 12.2771 17.3961C11.8093 18.3519 11.3415 19.3074 10.8738 20.2628C10.5574 20.9091 10.241 21.5554 9.92457 22.2017C10.8745 22.5946 11.9182 22.8117 13.0136 22.8117C14.1003 22.8117 15.1363 22.598 16.0803 22.2109ZM9.19587 21.8564C9.41455 21.4103 9.63323 20.9644 9.85191 20.5184C10.418 19.364 10.9841 18.2093 11.5496 17.0551C11.1391 16.7763 10.8312 16.3907 10.6501 15.9458C9.07393 16.2997 7.49851 16.6528 5.92309 17.0059L5.91566 17.0076C5.727 17.0499 5.53835 17.0922 5.34969 17.1345C5.9399 19.167 7.34508 20.862 9.19587 21.8564ZM5.16733 16.3574C5.51754 16.2789 5.86778 16.2004 6.21803 16.122C7.63484 15.8045 9.05192 15.4869 10.469 15.1688C10.463 15.1154 10.4599 15.062 10.4599 15.0086C10.4599 14.5697 10.5777 14.1515 10.789 13.7838L9.68994 12.9227C8.6336 12.0952 7.57763 11.2679 6.5216 10.4402C5.59177 11.7211 5.04477 13.2888 5.04477 14.9819C5.04477 15.4512 5.08679 15.9109 5.16733 16.3574ZM7.02657 9.81449C7.18344 9.93738 7.3403 10.0603 7.49717 10.1831C8.76271 11.1745 10.0283 12.1659 11.293 13.158C11.4259 13.0393 11.5707 12.9326 11.7307 12.8436C12.0024 12.6894 12.2952 12.5885 12.594 12.544V7.16283C10.3776 7.27579 8.40084 8.27855 7.02657 9.81449ZM13.406 7.16149V12.544C13.7048 12.5885 13.9976 12.6894 14.2693 12.8436C14.4293 12.9355 14.5741 13.0393 14.707 13.158C15.9716 12.166 17.237 11.1747 18.5024 10.1834L18.5037 10.1824C18.6656 10.0556 18.8275 9.92881 18.9894 9.80199C17.61 8.26685 15.6274 7.26732 13.406 7.16149ZM19.4951 10.4259C19.3333 10.5527 19.1714 10.6795 19.0096 10.8063C17.7441 11.7975 16.4787 12.7888 15.2141 13.7808C15.4254 14.1486 15.5401 14.5697 15.5401 15.0057C15.5401 15.0591 15.537 15.1124 15.531 15.1658C17.1071 15.5197 18.6825 15.8728 20.2579 16.2259L20.2653 16.2275C20.4633 16.2719 20.6613 16.3163 20.8592 16.3607C20.9402 15.9131 20.9824 15.4524 20.9824 14.9819C20.9824 13.2824 20.4313 11.7094 19.4951 10.4259ZM14.7643 14.9819C14.7643 14.0319 13.9805 13.2618 13.0136 13.2618C12.0467 13.2618 11.2629 14.0319 11.2629 14.9819C11.2629 15.932 12.0467 16.7021 13.0136 16.7021C13.9805 16.7021 14.7643 15.932 14.7643 14.9819Z"
            fill="cuurentColor"
          />
        </svg>
      ),
      subItens: SubItens,
      url: '/'
    },
    {
      name: 'Machines',
      icon: (
        <svg
          width="24"
          height="24"
          style={{ height: '1.5rem', width: '1.5rem', marginBottom: '4px', marginRight: '14px', marginLeft: '10px' }}
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M10.893 22.4565V1.5435L16.959 0V24L10.893 22.4565ZM12.355 19.9775C12.355 20.7475 12.9545 21.5035 13.7175 21.6715C14.508 21.846 15.1675 21.3265 15.1675 20.5055C15.1675 19.6845 14.508 18.9165 13.7175 18.7945C12.955 18.6775 12.355 19.2075 12.355 19.9775ZM11.2245 16.2805L16.53 16.8305V14.824L11.2245 14.5035V16.2805ZM11.2245 12.527L16.53 12.5945V10.5895L11.2245 10.75V12.527ZM11.2245 8.772L16.53 8.3585V6.35L11.2245 6.993V8.772ZM11.2245 5.0165L16.53 4.1225V2.1155L11.2245 3.238V5.0165Z"
            fill="currentColor"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M17.326 24V0L24 2.4225V21.5775L17.326 24ZM22.2215 19.6055L23.536 19.274V4.726L22.2215 4.3945V19.6055Z"
            fill="currentColor"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M0 4.171V19.809L4.53599 20.963V3.017L0 4.171ZM2.112 19.221C1.5415 19.0955 1.093 18.53 1.093 17.9545C1.093 17.379 1.543 16.9815 2.112 17.0695C2.7035 17.162 3.2 17.7365 3.2 18.35C3.2 18.9635 2.7035 19.3525 2.112 19.222V19.221ZM4.215 15.6025L0.25 15.1925V13.862L4.215 14.1V15.6025ZM4.215 12.4345L0.25 12.384V11.054L4.21699 10.9345L4.215 12.4345ZM4.215 9.267L0.25 9.576V8.246L4.21699 7.7665L4.215 9.267ZM4.215 6.0995L0.25 6.768V5.438L4.215 4.6V6.0995Z"
            fill="currentColor"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M4.81049 20.963V3.017L9.75549 4.8285V19.15L4.81049 20.963ZM8.47099 17.6795L9.45399 17.4295V6.553L8.47099 6.303V17.6795Z"
            fill="currentColor"
          />
        </svg>
      ),
      subItens: SubItensMachines,
      url: '/nodepools'
    },
    {
      name: 'Rbac Roles',
      icon: (
        <svg
          width="24"
          height="24"
          style={{ height: '1.5rem', width: '1.5rem', marginBottom: '4px', marginRight: '14px', marginLeft: '10px' }}
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M10.893 22.4565V1.5435L16.959 0V24L10.893 22.4565ZM12.355 19.9775C12.355 20.7475 12.9545 21.5035 13.7175 21.6715C14.508 21.846 15.1675 21.3265 15.1675 20.5055C15.1675 19.6845 14.508 18.9165 13.7175 18.7945C12.955 18.6775 12.355 19.2075 12.355 19.9775ZM11.2245 16.2805L16.53 16.8305V14.824L11.2245 14.5035V16.2805ZM11.2245 12.527L16.53 12.5945V10.5895L11.2245 10.75V12.527ZM11.2245 8.772L16.53 8.3585V6.35L11.2245 6.993V8.772ZM11.2245 5.0165L16.53 4.1225V2.1155L11.2245 3.238V5.0165Z"
            fill="currentColor"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M17.326 24V0L24 2.4225V21.5775L17.326 24ZM22.2215 19.6055L23.536 19.274V4.726L22.2215 4.3945V19.6055Z"
            fill="currentColor"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M0 4.171V19.809L4.53599 20.963V3.017L0 4.171ZM2.112 19.221C1.5415 19.0955 1.093 18.53 1.093 17.9545C1.093 17.379 1.543 16.9815 2.112 17.0695C2.7035 17.162 3.2 17.7365 3.2 18.35C3.2 18.9635 2.7035 19.3525 2.112 19.222V19.221ZM4.215 15.6025L0.25 15.1925V13.862L4.215 14.1V15.6025ZM4.215 12.4345L0.25 12.384V11.054L4.21699 10.9345L4.215 12.4345ZM4.215 9.267L0.25 9.576V8.246L4.21699 7.7665L4.215 9.267ZM4.215 6.0995L0.25 6.768V5.438L4.215 4.6V6.0995Z"
            fill="currentColor"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M4.81049 20.963V3.017L9.75549 4.8285V19.15L4.81049 20.963ZM8.47099 17.6795L9.45399 17.4295V6.553L8.47099 6.303V17.6795Z"
            fill="currentColor"
          />
        </svg>
      ),
      subItens: SubItens,
      url: '/roles'
    }
  ]

  const toggle = (i: any) => {
    show === i ? setShow(null) : setShow(i)
  }

  const handleRedirect = (link: string) => {
    history.push(link)
  }

  return (
    <>
      <PauseClusterAlert
        heading="Multiple clusters selected"
        isOpen={isPauseClusterAlertOpen}
        onActionConfirm={handlePause}
        onDismiss={closePauseClusterAlert}
      />
      <div className="menu-side-container">
        <ul className="side-itens">
          {location.pathname !== '/' && ( 
            <li onClick={() => history.push('/')} className="side-item-default">
              <svg
                width="27"
                height="27"
                style={{ marginBottom: '4px', marginRight: '12px', marginLeft: '8px' }}
                viewBox="0 0 27 27"
                fill="currentColor"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  fillRule="evenodd"
                  clipRule="evenodd"
                  d="M16.9328 23.9502C17.2666 24.6314 17.6002 25.3126 17.9337 25.9941C18.0557 25.9364 18.1769 25.8787 18.2981 25.821L18.299 25.8206C18.4205 25.7627 18.5419 25.7049 18.6642 25.6471C18.3309 24.9668 17.9975 24.2863 17.6639 23.6058C20.0532 22.3593 21.8628 20.1824 22.5923 17.5666C23.3346 17.7331 24.0768 17.8997 24.8189 18.0664C24.8491 17.9374 24.8793 17.8076 24.9094 17.6779L24.9105 17.6733C24.9403 17.5451 24.9702 17.4169 25 17.2894C24.2579 17.1228 23.5161 16.9563 22.7743 16.79C22.886 16.2041 22.9444 15.5997 22.9444 14.9819C22.9444 12.8279 22.234 10.8368 21.0309 9.22266C21.6226 8.759 22.2141 8.29525 22.8056 7.83134L22.2985 7.20851C21.7074 7.6721 21.1163 8.13554 20.525 8.59888C18.7838 6.6225 16.2482 5.34036 13.406 5.23185V3H12.594V5.23292C9.75682 5.34879 7.22707 6.63425 5.49109 8.61148C4.89446 8.14394 4.29792 7.67631 3.70155 7.20851L3.19444 7.83134C3.79176 8.30002 4.38892 8.76837 4.98598 9.23651C3.78909 10.8482 3.08276 12.834 3.08276 14.9819C3.08276 15.5986 3.14097 16.2018 3.25229 16.7867C2.50138 16.9551 1.75061 17.1237 1 17.2923C1.03018 17.4213 1.06037 17.5511 1.09055 17.6809C1.1207 17.8104 1.15096 17.9405 1.18111 18.0694C1.93226 17.9007 2.68325 17.7323 3.43414 17.5639C4.16108 20.1736 5.96314 22.3469 8.34346 23.5954C8.00746 24.2811 7.67155 24.967 7.33581 25.653C7.45798 25.7108 7.5794 25.7686 7.70082 25.8264C7.82232 25.8842 7.94403 25.9422 8.06628 26C8.40193 25.3135 8.73758 24.6273 9.07323 23.9413C10.2812 24.4549 11.6136 24.7395 13.0136 24.7395C14.4054 24.7395 15.7304 24.4582 16.9328 23.9502ZM15.3529 15.9429C16.9301 16.2973 18.508 16.651 20.0859 17.0046L20.0902 17.0056C20.2857 17.0494 20.4812 17.0932 20.6767 17.1371C20.084 19.1757 18.6713 20.8744 16.8116 21.8669C16.7301 21.7006 16.6486 21.5343 16.5671 21.3681L16.5576 21.3487C15.8549 19.9153 15.1523 18.4818 14.4504 17.0491C14.8609 16.7733 15.1718 16.3878 15.3529 15.9429ZM16.0803 22.2109C15.8595 21.7605 15.6386 21.3102 15.4178 20.8598C14.8516 19.7051 14.2855 18.5506 13.7199 17.3961C13.249 17.5355 12.748 17.5355 12.2771 17.3961C11.8093 18.3519 11.3415 19.3074 10.8738 20.2628C10.5574 20.9091 10.241 21.5554 9.92457 22.2017C10.8745 22.5946 11.9182 22.8117 13.0136 22.8117C14.1003 22.8117 15.1363 22.598 16.0803 22.2109ZM9.19587 21.8564C9.41455 21.4103 9.63323 20.9644 9.85191 20.5184C10.418 19.364 10.9841 18.2093 11.5496 17.0551C11.1391 16.7763 10.8312 16.3907 10.6501 15.9458C9.07393 16.2997 7.49851 16.6528 5.92309 17.0059L5.91566 17.0076C5.727 17.0499 5.53835 17.0922 5.34969 17.1345C5.9399 19.167 7.34508 20.862 9.19587 21.8564ZM5.16733 16.3574C5.51754 16.2789 5.86778 16.2004 6.21803 16.122C7.63484 15.8045 9.05192 15.4869 10.469 15.1688C10.463 15.1154 10.4599 15.062 10.4599 15.0086C10.4599 14.5697 10.5777 14.1515 10.789 13.7838L9.68994 12.9227C8.6336 12.0952 7.57763 11.2679 6.5216 10.4402C5.59177 11.7211 5.04477 13.2888 5.04477 14.9819C5.04477 15.4512 5.08679 15.9109 5.16733 16.3574ZM7.02657 9.81449C7.18344 9.93738 7.3403 10.0603 7.49717 10.1831C8.76271 11.1745 10.0283 12.1659 11.293 13.158C11.4259 13.0393 11.5707 12.9326 11.7307 12.8436C12.0024 12.6894 12.2952 12.5885 12.594 12.544V7.16283C10.3776 7.27579 8.40084 8.27855 7.02657 9.81449ZM13.406 7.16149V12.544C13.7048 12.5885 13.9976 12.6894 14.2693 12.8436C14.4293 12.9355 14.5741 13.0393 14.707 13.158C15.9716 12.166 17.237 11.1747 18.5024 10.1834L18.5037 10.1824C18.6656 10.0556 18.8275 9.92881 18.9894 9.80199C17.61 8.26685 15.6274 7.26732 13.406 7.16149ZM19.4951 10.4259C19.3333 10.5527 19.1714 10.6795 19.0096 10.8063C17.7441 11.7975 16.4787 12.7888 15.2141 13.7808C15.4254 14.1486 15.5401 14.5697 15.5401 15.0057C15.5401 15.0591 15.537 15.1124 15.531 15.1658C17.1071 15.5197 18.6825 15.8728 20.2579 16.2259L20.2653 16.2275C20.4633 16.2719 20.6613 16.3163 20.8592 16.3607C20.9402 15.9131 20.9824 15.4524 20.9824 14.9819C20.9824 13.2824 20.4313 11.7094 19.4951 10.4259ZM14.7643 14.9819C14.7643 14.0319 13.9805 13.2618 13.0136 13.2618C12.0467 13.2618 11.2629 14.0319 11.2629 14.9819C11.2629 15.932 12.0467 16.7021 13.0136 16.7021C13.9805 16.7021 14.7643 15.932 14.7643 14.9819Z"
                  fill="cuurentColor"
                />
              </svg>
            </li>
          )}
          {items.map((elm: any, i = 0) => {
            return (
              <>
                <li key={i} onClick={() => toggle(i)} className={Classnames('side-item', { active: show === i })}>
                  {typeof elm.icon === 'string' ? <i className={elm.icon} /> : elm.icon}
                  <p>{elm.name}</p>
                  {show === i ? <i className="icon-arrow-up" /> : <i className="icon-arrow-down" />}
                </li>
                {show === i && (
                  <div className="item-menu">
                    {elm.subItens.map((elm: TypeSubItem) => (
                      <p onClick={() => elm.handleAction?.() || handleRedirect(elm.link!)}>{elm.name}</p>
                    ))}
                  </div>
                )}
              </>
            )
          })}
        </ul>
      </div>
    </>
  )
}

export default MenuSideBar
