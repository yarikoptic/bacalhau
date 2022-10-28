import React, { FC, useState, useEffect, useMemo } from 'react'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import { DataGrid, GridColDef, GridValueGetterParams } from '@mui/x-data-grid'
import useApi from '../hooks/useApi'
import {
  getShortId,
  getJobStateTitle,
} from '../utils/job'
import {
  Job,
} from '../types'


import InputVolumes from '../components/job/InputVolumes'
import OutputVolumes from '../components/job/OutputVolumes'
import JobState from '../components/job/JobState'
import JobProgram from '../components/job/JobProgram'

const columns: GridColDef[] = [
  {
    field: 'id',
    headerName: 'ID',
    width: 100,
    renderCell: (params: any) => {
      return (
        <span style={{
          fontSize: '0.8em'
        }}>{ params.row.id }</span>
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
        <OutputVolumes
          storageSpecs={ params.row.outputs }
        />
      )
    },
  },
  {
    field: 'status',
    headerName: 'Status',
    width: 140,
    renderCell: (params: any) => {
      return (
        <JobState
          job={ params.row.job }
        />
      )
    },
  },
]

const Dashboard: FC = () => {
  const [ jobs, setJobs ] = useState<Job[]>([])
  const api = useApi()

  const rows = useMemo(() => {
    return jobs.map(job => {
      const {
        inputs = [],
        outputs = [],
        Docker,
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

  useEffect(() => {
    const doAsync = async () => {
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
    }
    doAsync()
  }, [])

  return (
    <Container maxWidth={ 'xl' } sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
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

export default Dashboard