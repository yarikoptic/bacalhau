import { createTheme, ThemeProvider } from '@mui/material/styles'
import AllContextProvider from './contexts/all'
import Layout from './pages/Layout'

const mdTheme = createTheme({
  palette: {
    primary: {
      main: '#04206F'
    },
    secondary: {
      main: '#1FC3CD'
    }
  } 
})

export default function App() {
  return (
    <AllContextProvider>
      <ThemeProvider theme={mdTheme}>
        <Layout />
      </ThemeProvider>
    </AllContextProvider>
  )
}
