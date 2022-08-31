import React, { FC, useState, useEffect, useMemo } from 'react'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import { DataGrid, GridColDef, GridValueGetterParams } from '@mui/x-data-grid'
import useApi from '../hooks/useApi'

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', width: 200 },
  { field: 'date', headerName: 'Date', width: 200 },
  { field: 'inputs', headerName: 'Inputs', width: 200 },
  { field: 'program', headerName: 'Program', width: 600 },
  { field: 'outputs', headerName: 'Outputs', width: 200 },
]

const Dashboard: FC = () => {
  const [ jobs, setJobs ] = useState<any[]>([])
  const api = useApi()

  const rows = useMemo(() => {
    console.dir(jobs)
    return jobs.map(job => {
      const {
        inputs,
        outputs,
        job_spec_docker,
      } = job.spec
      const inputCids = inputs.map((input: any) => input.cid)
      const outputPaths = outputs.map((output: any) => output.path)
      const dockerSpec = job_spec_docker as any
      return {
        id: job.id,
        date: job.created_at,
        inputs: inputCids.join(', '),
        program: `${job_spec_docker.image} ${(job_spec_docker.entrypoint || []).join(' ')}`,
        outputs: outputPaths.join(', '),
      }
    })
  }, [
    jobs,
  ])

  useEffect(() => {
    const doAsync = async () => {
      const result = await api.get('/api/jobs')
      const jobs = Object.values(result)
      jobs.sort((a: any, b: any) => {
        if(a.created_at > b.created_at) {
          return -1
        }
        if(a.created_at < b.created_at) {
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
              pageSize={20}
              rowsPerPageOptions={[5]}
              checkboxSelection
            />
          </div>
        </Grid>
      </Grid>
    </Container>
  )
}

export default Dashboard