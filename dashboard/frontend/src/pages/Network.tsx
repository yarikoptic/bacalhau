import React, { FC, useState, useEffect, useCallback } from 'react'
import prettyBytes from 'pretty-bytes'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import useApi from '../hooks/useApi'
import useLoadingErrorHandler from '../hooks/useLoadingErrorHandler'
import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import CardActions from '@mui/material/CardActions'
import CardContent from '@mui/material/CardContent'
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import {
  ClusterMapResult,
  DebugResponse,
} from '../types'
import {
  getShortId,
} from '../utils/job'
import ForceGraph from '../components/network/ForceGraph'

const Network: FC = () => {
  const [ mapData, setMapData ] = useState<ClusterMapResult>()
  const [ nodeData, setNodeData ] = useState<DebugResponse[]>()
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
          <Box sx={{
            display: 'inline-block',
          }}>
            {
              nodeData && nodeData.map((node, i) => {
                console.dir(node)
                return (
                  <Card sx={{ minWidth: 300, display: 'inline-block', m: 1 }} key={ i }>
                    <CardContent>
                      <Typography variant="h5" component="div">
                        { getShortId(node.ID) }
                      </Typography>
                      <Typography variant="body2">
                        <ul>
                          <li>
                            <div style={{minWidth: '140px', display: 'inline-block'}}>CPU:</div>
                            <strong>{ node.AvailableComputeCapacity.CPU }</strong>
                          </li>
                          <li>
                            <div style={{minWidth: '140px', display: 'inline-block'}}>Memory:</div>
                            <strong>{ prettyBytes(node.AvailableComputeCapacity.Memory || 0) }</strong>
                          </li>
                          <li>
                            <div style={{minWidth: '140px', display: 'inline-block'}}>Disk:</div>
                            <strong>{ prettyBytes(node.AvailableComputeCapacity.Disk || 0) }</strong>
                          </li>
                          <li>
                            <div style={{minWidth: '140px', display: 'inline-block'}}>GPU:</div>
                            <strong>{ node.AvailableComputeCapacity.GPU || 0 }</strong>
                          </li>
                          <li>
                            <div style={{minWidth: '140px', display: 'inline-block'}}>Requester Jobs:</div>
                            <strong>{ node.RequesterJobs.length || 0 }</strong>
                          </li>
                          <li>
                            <div style={{minWidth: '140px', display: 'inline-block'}}>Compute Jobs:</div>
                            <strong>{ node.ComputeJobs.length || 0 }</strong>
                          </li>
                        </ul>
                      </Typography>
                    </CardContent>
                    <CardActions>
                      <Button size="small">More Info</Button>
                    </CardActions>
                  </Card>
                )
              })
            }
          </Box>
        </Grid>
        <Grid item xs={6}>
          {
            mapData && (
              <ForceGraph data={ mapData } />
            )
          }
        </Grid>
      </Grid>
    </Container>
  )
}

export default Network