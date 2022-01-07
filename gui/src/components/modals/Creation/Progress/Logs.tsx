import { useState, VFC } from 'react'
import { useEffect } from 'react'

import styles from '@/components/modals/Creation/ClusterCreation.module.css'

type LogsProps = {
  namespace: string
  cluster: string
}

const Logs: VFC<LogsProps> = ({ namespace, cluster }: LogsProps) => {
  const [logs, setLogs] = useState<string[]>([])

  const addLogMessage = (message: string) => {
    setLogs(logs => [...logs, message])
    var elem = document.getElementById('progressDialog')
    elem.scrollTop = elem.scrollHeight
  }

  const streamUpdates = (controller: AbortController) => {
    fetch(`/api/events/${namespace}`, { signal: controller.signal })
      .then(response => {
        const stream = response.body.getReader()
        const utf8Decoder = new TextDecoder('utf-8')

        let buffer = ''

        return stream.read().then(function processText({ value }) {
          buffer += utf8Decoder.decode(value)

          buffer = onNewLine(buffer, chunk => {
            if (chunk.trim().length === 0) return

            try {
              const event = JSON.parse(chunk)
              const { object } = event
              const { involvedObject, message, reason } = object

              if (involvedObject.name.includes(cluster)) {
                const newMessage = `Reason: ${reason} Message: ${message}`

                addLogMessage(newMessage)
              }
            } catch (error) {
              console.log('Error while parsing', chunk, '\n', error)
            }
          })

          return stream.read().then(processText)
        })
      })
      .catch(err => {
        console.log('Error! Retrying in 5 seconds...')
        console.log(err)
        if (!controller.signal.aborted) setTimeout(() => streamUpdates(controller), 5000)
      })

    const onNewLine = (buffer, fn) => {
      const newLineIndex = buffer.indexOf('\n')

      if (newLineIndex === -1) return buffer

      const chunk = buffer.slice(0, buffer.indexOf('\n'))
      const newBuffer = buffer.slice(buffer.indexOf('\n') + 1)

      fn(chunk)

      return onNewLine(newBuffer, fn)
    }
  }

  useEffect(() => {
    const controller = new AbortController()
    streamUpdates(controller)

    return function cancel() {
      controller.abort()
    }
  }, [])

  return (
    <div id="progressDialog" className={styles.progressDialogContainer}>
      <ol className={styles.progressDialog}>
        {logs.map((log, index) => (
          <li key={index}>
            <a>{log}</a>
          </li>
        ))}
      </ol>
    </div>
  )
}

export { Logs }
