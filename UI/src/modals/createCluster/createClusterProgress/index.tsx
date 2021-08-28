import './createClusterProgress.scss'

type CreateClusterProgressProps = {
  hasFinished: boolean
  messages: string[]
}

export const CreateClusterProgress = ({ hasFinished, messages }: CreateClusterProgressProps) => {
  return (
    <>
      <h2 className="title-box">process</h2>
      <div className="cluster-progress-bar-container">
        <div className={`cluster-progress-bar ${hasFinished ? 'finished' : 'progress'}`} />
      </div>
      <div className="cluster-progress-container">
        {messages
          .filter(m => !!m)
          .map((m, i) => (
            <p key={i} className="cluster-progress-message">
              {m}
            </p>
          ))}
      </div>
    </>
  )
}