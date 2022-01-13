import type { FunctionComponent } from 'react'

import { ThemeProvider } from 'next-themes'
import { QueryCache, QueryClient, QueryClientProvider } from 'react-query'

import { ModalProvider } from '@/contexts/ModalContext'
import { ClusterProvider } from '@/contexts/ClusterContext'
import api from '@/lib/axios'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      queryFn: async ({ queryKey }) => {
        const { data } = await api(queryKey[0])
        return data
      },
      refetchInterval: 5 * 60 * 1000
    }
  },
  queryCache: new QueryCache()
})

const AppProviders: FunctionComponent = ({ children }) => (
  <QueryClientProvider client={queryClient}>
    <ModalProvider>
      <ClusterProvider>
        <ThemeProvider defaultTheme="dark-mode" themes={['dark-mode']}>
          {children}
        </ThemeProvider>
      </ClusterProvider>
    </ModalProvider>
  </QueryClientProvider>
)

export { AppProviders }
