import { useMemo } from 'react'
import { useRoutes } from 'hookrouter'

import {
  IRouteObject,
  routes,
} from '../routes'

export const useRoute = (): IRouteObject => {
  const routeResult = useRoutes(routes) as IRouteObject | undefined
  const route = useMemo<IRouteObject>(() => {
    return routeResult || {
      title: 'Not Found',
      render: () => (
        <div>Page not found</div>
      ),
      params: {}, 
    }
  }, [routeResult]) 
  return route
}

export default useRoute