import type { FC, ReactNode } from 'react'
import React from 'react'
import { ClusterInfo } from './ClusterInfo'
import { InfraProvider } from './InfraProvider'
import { AddOns } from './AddOns'
import { ControlPlane } from './ControlPlane'

type StepProps = {
  step: number
  currentStep: number
  children: ReactNode
}

const Step: FC<StepProps> = ({ step, currentStep, children }: StepProps) => {
  const stepStyles = step !== currentStep ? { display: 'none' } : {}
  return <div style={stepStyles}>{children}</div>
}

export { Step, ClusterInfo, InfraProvider, AddOns, ControlPlane }
