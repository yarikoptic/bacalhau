import React, { FC, useState, useEffect } from 'react'
import bluebird from 'bluebird'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import Typography from '@mui/material/Typography'
import Paper from '@mui/material/Paper'
import Divider from '@mui/material/Divider'
import useApi from '../hooks/useApi'
import {
  JobInfo,
} from '../types'
import InputVolumes from '../components/job/InputVolumes'
import OutputVolumes from '../components/job/OutputVolumes'
import JobState from '../components/job/JobState'
import JobProgram from '../components/job/JobProgram'
import {
  SmallText,
  SmallLink,
  BoldSectionTitle,
} from '../components/widgets/GeneralText'
import TerminalWindow from '../components/widgets/TerminalWindow'
import useLoadingErrorHandler from '../hooks/useLoadingErrorHandler'

const JobPage: FC<{
  id: string,
}> = ({
  id,
}) => {
  const [ jobInfo, setJobInfo ] = useState<JobInfo>()
  const [ jobSpecOpen, setJobSpecOpen ] = useState(false)
  const api = useApi()
  const loadingErrorHandler = useLoadingErrorHandler()

  useEffect(() => {
    const doAsync = loadingErrorHandler(async () => {
      const info = await api.post('/api/jobinfo', {
        id,
      })
      setJobInfo(info)
    })
    doAsync()
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
              <Grid item xs={12}>
                <BoldSectionTitle>
                  Job Details
                </BoldSectionTitle>
              </Grid>
              <Grid item xs={2}>
                <Typography variant="caption">
                  ID:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <SmallText>
                  { jobInfo.job.ID }
                </SmallText>
              </Grid>
              <Grid item xs={2}>
                <Typography variant="caption">
                  Date:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <SmallText>
                  { new Date(jobInfo.job.CreatedAt).toLocaleDateString() + ' ' + new Date(jobInfo.job.CreatedAt).toLocaleTimeString()}
                </SmallText>
              </Grid>
              <Grid item xs={2}>
                <Typography variant="caption">
                  Status:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <JobState
                  job={ jobInfo.job }
                />
              </Grid>
              <Grid item xs={12}>
                <Divider sx={{
                  mt: 1,
                  mb: 1,
                }} />
              </Grid>
              <Grid item xs={2}>
                <Typography variant="caption">
                  Inputs:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <InputVolumes
                  storageSpecs={ jobInfo.job.Spec.inputs || [] }
                />
              </Grid>
              <Grid item xs={12}>
                <Divider sx={{
                  mt: 1,
                  mb: 1,
                }} />
              </Grid>
              <Grid item xs={2}>
                <Typography variant="caption">
                  Program:
                </Typography>
              </Grid>
              <Grid
                item
                xs={10}
                sx={{
                  cursor: 'pointer',
                }}
                onClick={() => setJobSpecOpen(true)}
              >
                <JobProgram
                  job={ jobInfo.job }
                />
              </Grid>
              <Grid item xs={12} sx={{
                display: 'flex',
                justifyContent: 'center',
              }}>
                <SmallLink
                  onClick={() => setJobSpecOpen(true)}
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
              <Grid item xs={2}>
                <Typography variant="caption">
                  Outputs:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <OutputVolumes
                  storageSpecs={ jobInfo.job.Spec.outputs || [] }
                />
              </Grid>
            </Grid>
          </Paper>
        </Grid>
        <Grid item xs={4}>
        <Paper
            sx={{
              p: 2,
            }}
          >
            <Typography variant="subtitle1">
              Nodes
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={4}>
        <Paper
            sx={{
              p: 2,
            }}
          >
            <Typography variant="subtitle1">
              Events
            </Typography>
          </Paper>
        </Grid>
      </Grid>
      {
        jobSpecOpen && (
          <TerminalWindow
            open
            title="Job Spec"
            backgroundColor="#fff"
            color="#000"
            data={ jobInfo.job.Spec }
            onClose={ () => setJobSpecOpen(false) }
          />
        )
      }
    </Container>
  )
}

export default JobPage