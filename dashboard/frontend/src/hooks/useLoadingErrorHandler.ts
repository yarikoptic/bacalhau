import useLoading from './useLoading'
import useSnackbar from './useSnackbar'

type asyncFunction = {
  (): Promise<void>,
}

export const useLoadingErrorHandler = ({
  withSnackbar = true,
  withLoading = true,
}: {
  withSnackbar?: boolean,
  withLoading?: boolean,
} = {}) => {
  const loading = useLoading()
  const snackbar = useSnackbar()
  return (handler: asyncFunction): asyncFunction => {
    return async () => {
      if(withLoading) loading.setLoading(true)
      try {
        await handler()
      } catch(e: any) {
        console.error(e)
        if(withSnackbar) snackbar.error(e.toString())
      }
      if(withLoading) loading.setLoading(false)
    }
  }
}

export default useLoadingErrorHandler