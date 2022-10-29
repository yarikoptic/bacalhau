import React, { FC, useState, useEffect, useCallback } from 'react'
import bluebird from 'bluebird'
import { A, navigate } from 'hookrouter'
import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import Typography from '@mui/material/Typography'
import Paper from '@mui/material/Paper'
import IconButton from '@mui/material/IconButton'
import Tooltip from '@mui/material/Tooltip'
import Divider from '@mui/material/Divider'
import RefreshIcon from '@mui/icons-material/Refresh'
import useApi from '../hooks/useApi'
import {
  JobInfo,
} from '../types'
import {
  getShortId,
  getJobStateTitle,
} from '../utils/job'
import InputVolumes from '../components/job/InputVolumes'
import OutputVolumes from '../components/job/OutputVolumes'
import JobState from '../components/job/JobState'
import ShardState from '../components/job/ShardState'
import JobProgram from '../components/job/JobProgram'
import {
  SmallText,
  SmallLink,
  TinyText,
  BoldSectionTitle,
  RequesterNode,
} from '../components/widgets/GeneralText'
import TerminalWindow from '../components/widgets/TerminalWindow'
import useLoadingErrorHandler from '../hooks/useLoadingErrorHandler'

type JSONWindowConfig = {
  title: string,
  data: any,
}

const InfoRow: FC<{
  title: string,
  rightAlign?: boolean,
  withDivider?: boolean,
}> = ({
  title,
  rightAlign = false,
  withDivider = false,
  children,
}) => {
  return (
    <>
      <Grid item xs={3}>
        <Typography variant="caption">
          { title }:
        </Typography>
      </Grid>
      <Grid item xs={9} sx={{
        pl: 8,
        display: 'flex',
        alignItems: 'center',
        justifyContent: rightAlign ? 'flex-end' : 'flex-start',
      }}>
        { children }
      </Grid>
      {
        withDivider && (
          <Grid item xs={12}>
            <Divider sx={{
              mt: 1,
              mb: 1,
            }} />
          </Grid>
        )
      }
    </>
  )
}

const JobPage: FC<{
  id: string,
}> = ({
  id,
}) => {
  const [ jobInfo, setJobInfo ] = useState<JobInfo>()
  const [ jsonWindow, setJsonWindow ] = useState<JSONWindowConfig>()
  const api = useApi()
  const loadingErrorHandler = useLoadingErrorHandler()

  const isRequesterNodeID = useCallback((id: string): boolean => {
    if(!jobInfo) return false
    return jobInfo.job.RequesterNodeID == id
  }, [
    jobInfo,
  ])

  const loadInfo = useCallback(async () => {
    const handler = loadingErrorHandler(async () => {
      const info = await api.post('/api/jobinfo', {
        id,
      })
      setJobInfo(info)
    })
    await handler()
  }, [])

  useEffect(() => {
    loadInfo()
  }, [])

  if(!jobInfo) return null

  return (
    <Container maxWidth={ 'xl' } sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        <Grid item xs={4}>
          <Paper
            sx={{
              p: 2,
            }}
          >
            <Grid container spacing={1}>
              <Grid item xs={6}>
                <BoldSectionTitle>
                  Job Details
                </BoldSectionTitle>
              </Grid>
              <Grid item xs={6} sx={{
                display: 'flex',
                justifyContent: 'flex-end',
              }}>
                <Tooltip title="Refresh">
                  <IconButton aria-label="delete" color="primary" onClick={ loadInfo }>
                    <RefreshIcon />
                  </IconButton>
                </Tooltip>
              </Grid>
              <InfoRow title="ID">
                <SmallText>
                  { jobInfo.job.ID }
                </SmallText>
              </InfoRow>
              <InfoRow title="Date">
                <SmallText>
                  { new Date(jobInfo.job.CreatedAt).toLocaleDateString() + ' ' + new Date(jobInfo.job.CreatedAt).toLocaleTimeString()}
                </SmallText>
              </InfoRow>
              <InfoRow title="Concurrency">
                <SmallText>
                  { jobInfo.job.Deal.Concurrency }
                </SmallText>
              </InfoRow>
              <InfoRow title="Shards">
                <SmallText>
                { jobInfo.job.ExecutionPlan.ShardsTotal }
                </SmallText>
              </InfoRow>
              <InfoRow title="State" withDivider>
                <JobState
                  job={ jobInfo.job }
                />
              </InfoRow>
              <InfoRow title="Inputs" withDivider>
                <InputVolumes
                  storageSpecs={ jobInfo.job.Spec.inputs || [] }
                />
              </InfoRow>
              <Grid item xs={12} sx={{
                direction: 'column',
                display: 'flex',
                justifyContent: 'center',
              }}>
                <Box
                  sx={{
                    cursor: 'pointer',
                  }}
                  onClick={() => setJsonWindow({
                    title: 'Program',
                    data: jobInfo.job.Spec,
                  })}
                >
                  <JobProgram
                    job={ jobInfo.job }
                  />
                </Box>
                <br />
                
              </Grid>
              <Grid item xs={12} sx={{
                direction: 'column',
                display: 'flex',
                justifyContent: 'center',
              }}>
                <SmallLink
                  onClick={() => setJsonWindow({
                    title: 'Program',
                    data: jobInfo.job.Spec,
                  })}
                >
                  view info
                </SmallLink>
              </Grid>
              <Grid item xs={12}>
                <Divider sx={{
                  mt: 1,
                  mb: 1,
                }} />
              </Grid>
              <InfoRow title="Outputs" withDivider>
                <OutputVolumes
                  outputVolumes={ jobInfo.job.Spec.outputs || [] }
                />
              </InfoRow>
            </Grid>
          </Paper>
        </Grid>
        <Grid item xs={4}>
          <Paper
            sx={{
              p: 2,
              mb: 2,
            }}
          >
            <Grid container spacing={1}>
              <Grid item xs={12}>
                <BoldSectionTitle>
                  Nodes
                </BoldSectionTitle>
              </Grid>
              <Grid item xs={3}>
                <Typography variant="caption">
                  Requester Node:
                </Typography>
              </Grid>
              <Grid item xs={9}>
                <SmallText>
                  <RequesterNode>
                    { getShortId(jobInfo.job.RequesterNodeID) }
                  </RequesterNode>
                </SmallText>
              </Grid>
            </Grid>
          </Paper>
          {
            Object.keys(jobInfo.state.Nodes).map(nodeID => {
              const nodeState = jobInfo.state.Nodes[nodeID]
              return (
                <Paper
                  key={ nodeID }
                  sx={{
                    p: 2,
                    mb: 2,
                  }}
                >
                  <Grid container spacing={0.5}>
                    <Grid item xs={12}>
                      <BoldSectionTitle>
                        <A href="/network">
                          { getShortId(nodeID) }
                        </A>
                      </BoldSectionTitle>
                    </Grid>
                    {
                      Object.keys(nodeState.Shards).map((shardIndex, i) => {
                        const shardState = nodeState.Shards[shardIndex as unknown as number]
                        return (
                          <React.Fragment key={ shardIndex }>
                            <InfoRow title="Shard Index">
                              <SmallText>
                                { shardIndex }
                              </SmallText>
                            </InfoRow>
                            <InfoRow title="State">
                              <SmallText>
                                <ShardState state={ shardState.State } />
                              </SmallText>
                            </InfoRow>
                            <InfoRow title="Status">
                              <TinyText>
                                exitCode: { shardState.RunOutput.exitCode } &nbsp;
                                <span style={{color:'#999'}}>{ shardState.Status }</span>
                              </TinyText>
                            </InfoRow>
                            {
                              shardState.RunOutput.stdout && (
                                <InfoRow title="stdout">
                                  <TinyText>
                                    <span style={{color:'#999'}}>{ shardState.RunOutput.stdout }</span>
                                  </TinyText>
                                </InfoRow>
                              )
                            }
                            {
                              shardState.RunOutput.stderr && (
                                <InfoRow title="stderr">
                                  <TinyText>
                                    <span style={{color:'#999'}}>{ shardState.RunOutput.stderr }</span>
                                  </TinyText>
                                </InfoRow>
                              )
                            }
                            <InfoRow title="Outputs" withDivider={ i < Object.keys(nodeState.Shards).length - 1 }>
                              <OutputVolumes
                                outputVolumes={ jobInfo.job.Spec.outputs || [] }
                                publishedResults={ shardState.PublishedResults }
                              />
                            </InfoRow>
                          </React.Fragment>
                        )
                      })
                    }
                  </Grid>
                </Paper>
              )
            })
          }
        </Grid>
        <Grid item xs={4}>
        <Paper
            sx={{
              p: 2,
            }}
          >
            <Grid container spacing={0.5}>
              <Grid item xs={8}>
                <BoldSectionTitle>
                  Events
                </BoldSectionTitle>
              </Grid>
              <Grid item xs={4} sx={{
                display: 'flex',
                justifyContent: 'flex-end',
              }}>
                <SmallLink
                  onClick={() => setJsonWindow({
                    title: 'Events',
                    data: jobInfo.events,
                  })}
                >
                  view all
                </SmallLink>
              </Grid>
              <Grid item xs={4}>
                <SmallText>
                  <strong>Node</strong>
                </SmallText>
              </Grid>
              <Grid item xs={4}>
              <SmallText>
                  <strong>Event</strong>
                </SmallText>
              </Grid>
              <Grid item xs={4}>
                <SmallText>
                  <strong>Date</strong>
                </SmallText>
              </Grid>
              {
                jobInfo.events.map((event, i) => {
                  return (
                    <React.Fragment key={ i }>
                      <Grid item xs={4}>
                        <SmallText>
                          {
                            isRequesterNodeID(event.SourceNodeID) && (event.TargetNodeID || event.EventName == 'Created') ? (
                              <RequesterNode>
                                { getShortId(event.SourceNodeID) }
                              </RequesterNode>
                            ) : getShortId(event.SourceNodeID)
                          }
                        </SmallText>
                      </Grid>
                      <Grid item xs={4}>
                        <SmallLink
                          onClick={() => setJsonWindow({
                            title: 'Event',
                            data: event,
                          })}
                        >
                          { event.EventName }
                        </SmallLink>
                      </Grid>
                      <Grid item xs={4}>
                        <TinyText>
                          { new Date(event.EventTime).toLocaleDateString() + ' ' + new Date(event.EventTime).toLocaleTimeString()}
                        </TinyText>
                      </Grid>
                      
                    </React.Fragment>
                  )
                })
              }
            </Grid>
          </Paper>
        </Grid>
      </Grid>
      {
        jsonWindow && (
          <TerminalWindow
            open
            title={ jsonWindow.title }
            backgroundColor="#fff"
            color="#000"
            data={ jsonWindow.data }
            onClose={ () => setJsonWindow(undefined) }
          />
        )
      }
    </Container>
  )
}

export default JobPage