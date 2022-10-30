import React, { FC, useState, useEffect, useCallback } from 'react'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import useApi from '../hooks/useApi'
import useLoadingErrorHandler from '../hooks/useLoadingErrorHandler'

import {
  ClusterMapResult,
  DebugResponse,
} from '../types'

import ForceGraph from '../components/network/ForceGraph'

const Network: FC = () => {
  const [ mapData, setMapData ] = useState<ClusterMapResult>()
  const [ nodeData, setNodeData ] = useState<DebugResponse>()
  const api = useApi()
  const loadingErrorHandler = useLoadingErrorHandler()

  const loadMapData = useCallback(async () => {
    const handler = loadingErrorHandler(async () => {
      const mapData = await api.post('/api/map', {})
      setMapData(mapData)
    })
    await handler()
  }, [])

  const loadNodeData = useCallback(async () => {
    const handler = loadingErrorHandler(async () => {
      const nodeData = await api.post('/api/nodes', {})
      setNodeData(nodeData)
    })
    await handler()
  }, [])

  useEffect(() => {
    loadMapData()
    loadNodeData()
  }, [])

  return (
    <Container maxWidth={ 'xl' } sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        <Grid item xs={6}>
          Nodes
        </Grid>
        <Grid item xs={6}>
          {
            mapData && (
              <ForceGraph data={ mapData } />
            )
          }
          {/* <iframe width="600px" height="600px" src="/html/viz.html" frameBorder="none" style={{
            border: '1px solid #999'
          }}></iframe> */}
        </Grid>
      </Grid>
    </Container>
  )
}

export default Network