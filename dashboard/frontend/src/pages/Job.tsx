import React, { FC, useState, useEffect } from 'react'
import bluebird from 'bluebird'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import Typography from '@mui/material/Typography'
import Paper from '@mui/material/Paper'
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
  BoldSectionTitle,
} from '../components/widgets/GeneralText'
import useLoadingErrorHandler from '../hooks/useLoadingErrorHandler'

const JobPage: FC<{
  id: string,
}> = ({
  id,
}) => {
  const [ jobInfo, setJobInfo ] = useState<JobInfo>()
  const api = useApi()
  const loadingErrorHandler = useLoadingErrorHandler()

  useEffect(() => {
    const doAsync = loadingErrorHandler(async () => {
      const info = await api.post('/api/jobinfo', {
        id,
      })
      await bluebird.delay(2000)
      
      setJobInfo(info)

      throw new Error('test')
    })
    doAsync()
  }, [])

  console.dir(jobInfo)

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
                  Inputs:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <InputVolumes
                  storageSpecs={ jobInfo.job.Spec.inputs || [] }
                />
              </Grid>
              <Grid item xs={2}>
                <Typography variant="caption">
                  Program:
                </Typography>
              </Grid>
              <Grid item xs={10}>
                <JobProgram
                  job={ jobInfo.job }
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
    </Container>
  )
}

export default JobPage