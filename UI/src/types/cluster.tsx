export type TypeWorker = {
  id: string
  machineType: string
  replicas: number
  infraNode: boolean
}

export type TypeSubnet = {
  id: string
  cidrBlock: string
  zone: string
  isPublic: boolean
}

export type TypeOption = {
  value: string
  label: string
}

export type TypeSelectOptions = {
  [instanceType: string]: {
    selectOptions: TypeOption[];
  };
};

export type TypeAsyncSelect = {
  label?: string
  onChange: (option: TypeOption | null) => void
  loadOptions: any
  value: TypeOption | null
}

export type TypeCluster = {
  clusterName: string
  setClusterName: Function
  namespace: string
  setNamespace: Function
  provider: string
  setProvider: Function,
  providerOptions: any
  regionOptions: TypeOption[]
  region: string
  setRegion: Function
  accessKey: string
  setAccesskey: Function
  secret: string
  setSecret: Function
  session: string
  setSession: Function
}

export type TypeInfra = {
  provider: string
  setProvider: Function
  providerOptions: any
  regionOptions: TypeOption[]
  region: string
  setRegion: Function
  flavor: string
  setFlavor: Function
  flavorOptions: TypeOption[]
  k8sVersion: string
  setK8sVersion: Function
  sshKey: string
  setSshKey: Function
  k8sOptions: any
}

export type TypeTaints = {
  key: string
  value: string
  effect: string
}

export type TypeControlPlane = {
  replicas: number
  setReplicas: Function
  cpu: string
  setCpu: Function
  getCpu: TypeOption[] | undefined
  getMem: TypeOption[] | undefined
  memory: string
  setMemory: Function
  machineTypes: string
  setMachineTypes: Function
  getMachineTypes: TypeOption[] | undefined
  infraNode?: boolean | any 
  setInfraNode?: Function
  replicasWorkers?: number | any
  setReplicasWorkers?: Function
  cpuWorkers?: string
  setCpuWorkers?: Function
  memoryWorkers?: string
  setMemoryWorkers?: Function
  machineTypesWorkers?: string
  setMachineTypesWorkers?: Function
  createWorkers?: () => void
  workers?: TypeWorker[]
  deleteWorkers?: Function
  clusterName?: string
  isAdvanced?: boolean
  keyTaint?: string
  setKeyTaint?: Function
  valueTaint?: string
  setValueTaint?: Function
  effectValue?: string
  setEffectValue?: Function
  effect?: TypeOption[]
  keyLabel?: string
  setKeyLabel?: Function
  valueLabel?: string
  setValueLabel?: Function
  keyProv?: string
  setKeyProv?: Function
  valueProv?: string
  setValueProv?: Function
  taints?: TypeTaints[]
  labels?: {}[]
  providers?: {}[]
  deleteTaints?: Function
  deleteLabels?: Function
  deleteProviders?: Function
  handleActionTaints?: Function
  handleActionLabel?: Function
  handleActionProv?: Function
  internalLB?: boolean
  setInternalLB?: Function
}

export type TypeWorkersAdvanced = {
  replicas: number
  setReplicas: Function
  cpu: string
  setCpu: Function
  getCpu: TypeOption[] | undefined
  getMem: TypeOption[] | undefined
  memory: string
  setMemory: Function
  machineTypes: string
  setMachineTypes: Function
  getMachineTypes: TypeOption[] | undefined
  autoScale: boolean
  setAutoScale: Function
  maxSize: number
  setMaxSize: Function
  minSize: number
  setMinSize: Function
  keyTaint?: string
  setKeyTaint?: Function
  valueTaint?: string
  setValueTaint?: Function
  effectValue?: string
  setEffectValue?: Function
  effect?: TypeOption[]
  keyLabel?: string
  setKeyLabel?: Function
  valueLabel?: string
  setValueLabel?: Function
  keyProv?: string
  setKeyProv?: Function
  valueProv?: string
  setValueProv?: Function
  taints?: TypeTaints[]
  labels?: {}[]
  providers?: {}[]
  deleteTaints?: Function
  deleteLabels?: Function
  deleteProviders?: Function
  handleActionTaints?: Function
  handleActionLabel?: Function
  handleActionProv?: Function
  groupIdOptions: TypeOption[]
  groupId: string
  setGroupId: Function
  handleAction: Function
  subnet: string
  setSubnet: Function
}

export type TypeNetwork = {
  id: string
  setId: Function
  idSubnet: string
  setIdSubnet: Function
  isPublic: boolean
  setIsPublic: Function
  zone: string
  setZone: Function
  cidr: string
  setCidr: Function
  cidrSubnet: string
  setCidrSubnet: Function
  addSubnet: Function
  subnets: TypeSubnet[]
  deleteSubnet: Function
}

export type TypeK8sNetwork = {
  serverPort: string
  setServerPort: Function
  serviceDomain: string
  setServiceDomain: Function
  podsRanges: string
  setPodsRanges: Function
  serviceRanges: string
  setServiceRanges: Function
}

export type TypeBastion = {
  enabled: boolean
  setEnabled: Function
  ingress: boolean
  setIngress: Function
  replicas: number
  setReplicas: Function
  cpu: string
  setCpu: Function
  getCpu: TypeOption[] | undefined
  getMem: TypeOption[] | undefined
  memory: string
  setMemory: Function
  machineTypes: string
  setMachineTypes: Function
  getMachineTypes: TypeOption[] | undefined
  cidr: string
  setCidr: Function
  cidrs: string[]
  handleEvent: Function
  deleteCidr: Function
}

