import React, { FC, useState, useEffect, useMemo } from 'react'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'
import { DataGrid, GridColDef, GridValueGetterParams } from '@mui/x-data-grid'
import useApi from '../hooks/useApi'

const FILECOIN_PLUS_CIDS = [
  'Qmd9CBYpdgCLuCKRtKRRggu24H72ZUrGax5A9EYvrbC72j',
  'QmeZRGhe4PmjctYVSVHuEiA9oSXnqmYa4kQubSHgWbjv72',
]

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', width: 200 },
  { field: 'date', headerName: 'Date', width: 200 },
  {
    field: 'status',
    headerName: 'Status',
    width: 100,
    renderCell: (params: any) => {
      let hasFilecoinPlus = false
      const inputCids = params.row.inputCids || []
      inputCids.forEach((inputCid: any) => {
        if (FILECOIN_PLUS_CIDS.includes(inputCid)) {
          hasFilecoinPlus = true
        }
      })
      if (!hasFilecoinPlus) return ''
      return (
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
          }}
        >
          <img
            style={{
              width: '30px',
              height: '30px',
            }}
            src="/img/filecoin-logo.png" alt="Filecoin Plus"
          />
          <span style={{fontSize: '2em', marginTop: '3px', marginLeft: '3px'}}>
           +
          </span>
        </div>
        
      )
    },
  },
  { field: 'inputs', headerName: 'Inputs', width: 200 },
  { field: 'program', headerName: 'Program', width: 500 },
  { field: 'outputs', headerName: 'Outputs', width: 200 },
  
]

const Dashboard: FC = () => {
  const [ jobs, setJobs ] = useState<any[]>([])
  const api = useApi()

  const rows = useMemo(() => {
    return jobs.map(job => {
      const {
        Inputs = [],
        Outputs = [],
        Docker,
      } = job.Spec
      const inputCids = Inputs.map((input: any) => input.cid)
      const outputPaths = Outputs.map((output: any) => output.path)
      const dockerSpec = Docker as any
      return {
        id: job.ID,
        date: job.CreatedAt,
        inputs: inputCids.join(', '),
        inputCids,
        program: `${dockerSpec.Image} ${(dockerSpec.Entrypoint || []).join(' ')}`,
        outputs: outputPaths.join(', '),
        status: 'X'
      }
    })
  }, [
    jobs,
  ])

  useEffect(() => {
    const doAsync = async () => {
      const jobs = await api.post('/api/jobs', {
        maxJobs: 1,
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
              checkboxSelection
            />
          </div>
        </Grid>
      </Grid>
    </Container>
  )
}

export default Dashboard