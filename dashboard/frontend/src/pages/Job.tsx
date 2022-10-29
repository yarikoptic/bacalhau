import React, { FC, useState, useEffect } from 'react'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import Box from '@mui/material/Box'
import useApi from '../hooks/useApi'
import {
  Job,
} from '../types'
import { RouterContext } from '../contexts/router'

const JobPage: FC<{
  id: string,
}> = ({
  id,
}) => {
  const [ job, setJob ] = useState<Job>()
  const api = useApi()

  useEffect(() => {
    const doAsync = async () => {
      const job = await api.post('/api/job', {
        id,
      })
      setJob(job)
    }
    doAsync()
  }, [])

  console.dir(job)

  return (
    <Container maxWidth={ 'xl' } sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          Job
        </Grid>
      </Grid>
    </Container>
  )
}

export default JobPage