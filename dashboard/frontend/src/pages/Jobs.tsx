import React, { FC, useState, useEffect, useMemo, useCallback } from 'react'
import { A, navigate } from 'hookrouter'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Box from '@mui/material/Box'
import IconButton from '@mui/material/IconButton'
import Tooltip from '@mui/material/Tooltip'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import useApi from '../hooks/useApi'
import {
  getShortId,
  getJobStateTitle,
} from '../utils/job'
import {
  Job,
} from '../types'

import RefreshIcon from '@mui/icons-material/Refresh'
import InfoIcon from '@mui/icons-material/Info';
import InputVolumes from '../components/job/InputVolumes'
import OutputVolumes from '../components/job/OutputVolumes'
import JobState from '../components/job/JobState'
import JobProgram from '../components/job/JobProgram'
import useLoadingErrorHandler from '../hooks/useLoadingErrorHandler'

const columns: GridColDef[] = [
  {
    field: 'id',
    headerName: 'ID',
    width: 100,
    renderCell: (params: any) => {
      return (
        <span style={{
          fontSize: '0.8em'
        }}>
          <A href={`/jobs/${params.row.job.ID}`}>{ getShortId(params.row.job.ID) }</A>
        </span>
      )
    },
  },
  {
    field: 'date',
    headerName: 'Date',
    width: 120,
    renderCell: (params: any) => {
      return (
        <span style={{
          fontSize: '0.8em'
        }}>{ params.row.date }</span>
      )
    },
  },
  {
    field: 'inputs',
    headerName: 'Inputs',
    width: 200,
    renderCell: (params: any) => {
      return (
        <InputVolumes
          storageSpecs={ params.row.inputs }
        />
      )
    },
  },
  {
    field: 'program',
    headerName: 'Program',
    width: 500,
    renderCell: (params: any) => {
      return (
        <JobProgram
          job={ params.row.job }
        />
      )
    },
  },
  {
    field: 'outputs',
    headerName: 'Outputs',
    width: 200,
    renderCell: (params: any) => {
      return (
        <A href={`/jobs/${params.row.job.ID}`} style={{color: '#333'}}>
          <OutputVolumes
            outputVolumes={ params.row.outputs }
          />
        </A>
      )
    },
  },
  {
    field: 'state',
    headerName: 'State',
    width: 140,
    renderCell: (params: any) => {
      return (
        <JobState
          job={ params.row.job }
        />
      )
    },
  },
  {
    field: 'actions',
    headerName: 'Actions',
    flex: 1,
    renderCell: (params: any) => {
      return (
        <Box
          sx={{
            display: 'flex',   
            justifyContent: 'flex-start',
            alignItems: 'center',
            width: '100%',
          }}
          component="div"
        >
          <IconButton
            component="label"
            onClick={ () => navigate(`/jobs/${params.row.job.ID}`) }
          >
            <InfoIcon />
          </IconButton>
        </Box>
      )
    },
  },
]

const Jobs: FC = () => {
  const [ findJobID, setFindJobID ] = useState('')
  const [ jobs, setJobs ] = useState<Job[]>([])
  const api = useApi()
  const loadingErrorHandler = useLoadingErrorHandler()

  const rows = useMemo(() => {
    return jobs.map(job => {
      const {
        inputs = [],
        outputs = [],
      } = job.Spec
      return {
        job,
        id: getShortId(job.ID),
        date: new Date(job.CreatedAt).toLocaleDateString() + ' ' + new Date(job.CreatedAt).toLocaleTimeString(),
        inputs,
        outputs,
        shardState: getJobStateTitle(job),
      }
    })
  }, [
    jobs,
  ])

  const loadJobs = useCallback(async () => {
    const handler = loadingErrorHandler(async () => {
      const jobs = await api.post('/api/jobs', {
        maxJobs: 100,
        returnAll: true,
      })
      jobs.sort((a: any, b: any) => {
        if(a.CreatedAt > b.CreatedAt) {
          return -1
        }
        if(a.CreatedAt < b.CreatedAt) {
          return 1
        }
        return 0
      })
      setJobs(jobs)
    })
    await handler()
  }, [])

  const findJob = useCallback(async () => {
    const handler = loadingErrorHandler(async () => {
      if(!findJobID) throw new Error(`please enter a job id`)
      try {
        const info = await api.post('/api/jobinfo', {
          id: findJobID,
        })
        navigate(`/jobs/${info.job.ID}`)
      } catch(err) {
        throw new Error(`could not load job with id ${findJobID}`)
      }
    })
    await handler()
  }, [
    findJobID,
  ])

  useEffect(() => {
    loadJobs()
  }, [])

  return (
    <Container maxWidth={ 'xl' } sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        <Grid item xs={4}>
          <TextField
            fullWidth
            size="small"
            label="Find Job by ID"
            value={ findJobID }
            onChange={ (e) => setFindJobID(e.target.value) }
          />
        </Grid>
        <Grid item xs={2}>
          <Button
            size="small"
            variant="contained"
            onClick={ findJob }
          >
            Find Job
          </Button>
        </Grid>
        <Grid item xs={6} sx={{
          display: 'flex',
          justifyContent: 'flex-end',
        }}>
          <Tooltip title="Refresh">
            <IconButton aria-label="delete" color="primary" onClick={ loadJobs }>
              <RefreshIcon />
            </IconButton>
          </Tooltip>
        </Grid>
        <Grid item xs={12}>
          <div style={{ height: 800, width: '100%' }}>
            <DataGrid
              rows={rows}
              columns={columns}
              pageSize={25}
              rowsPerPageOptions={[10, 25, 100]}
            />
          </div>
        </Grid>
      </Grid>
    </Container>
  )
}

export default Jobs