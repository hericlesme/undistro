import './index.scss'
import { useState, useEffect } from 'react'
import Modals from 'util/modals'

const RbacTable = (props) => {

  const showModal = () => {
    Modals.show('create-role', {
      title: 'Create',
      ndTitle: 'Role',
      width: '774',
      height: '771'
    })
  }

  return (
    <table className="rbac-table">
      <thead>
        <tr>
          {props.header.map(elm => (
            <td key={elm.name}>
              <ColumnHeader data={elm} />
            </td>
          ))}
        </tr>
      </thead>
      <tbody>
        {props.data.map((elm, i) => {
          return (
            <>
              <Row
                onChange={props.onChange}
                header={props.header}
                key={i}
                data={elm}
                showModal={showModal}
              />
            </>
          )
        })}
      </tbody>
    </table>
  )
}

const ColumnHeader = props => {
  return (
    <div className="header">
      <label>{props.data.name}</label>
    </div>
  )
}

const Row = props => {
  const [keys, setKeys] = useState([])

  useEffect(() => {
    const headers = props.header.map(elm => elm.field)
    setKeys(Object.keys(props.data).filter(elm => headers.includes(elm)))
  }, [props.data, props.header])

  return (
    <tr>
      {keys.map((key) => {
        return (
          <td className={key} key={key}>
            <div>
              <span>{props.data[key]}</span>
            </div>
          </td>
        )
      })}
      <div className='add-role-text'>
        <i className='icon-add-circle' />
        <p onClick={() => props.showModal()}>Add roles</p>
      </div>
    </tr>
  )
}

export default RbacTable
