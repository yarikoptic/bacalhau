import Dashboard from './pages/Dashboard'

export type IRouteObject = {
  // show a different page title
  title?: string,
  render: {
    (): JSX.Element,
  },
  params: Record<string, any>,
}

export type IRouteFactory = (props: Record<string, any>) => IRouteObject

export const routes: Record<string, IRouteFactory> = {
  '/': () => ({
    title: 'Home',
    render: () => <Dashboard />,
    params: {},
  }),
}

export default routes