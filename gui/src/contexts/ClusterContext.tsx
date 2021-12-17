import { createContext, useContext } from 'react';

export type ClusterContextType = {
    clusters: string[];
    setClusters: (clusters: string[]) => void;
}

export const ClusterContext = createContext<ClusterContextType>({ clusters: [], setClusters: clusters => { } });
export const useClusters = () => useContext(ClusterContext);