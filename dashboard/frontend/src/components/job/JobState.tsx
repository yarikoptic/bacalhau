import React, { FC, useMemo } from 'react'
import { SxProps } from '@mui/system'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import {
  Job,
} from '../../types'

import FilPlus from './FilPlus'

import {
  getJobShardState,
  getShardStateTitle,
} from '../../utils/job'

const FILECOIN_PLUS_CIDS = [
  'Qmd9CBYpdgCLuCKRtKRRggu24H72ZUrGax5A9EYvrbC72j',
  'QmeZRGhe4PmjctYVSVHuEiA9oSXnqmYa4kQubSHgWbjv72',
]

const JobState: FC<{
  job: Job,
  sx?: SxProps,
}> = ({
  job,
  sx = {},
}) => {
  const shardState = useMemo(() => {
    const title = getShardStateTitle(getJobShardState(job))
    let color = '#666'
    if(title == 'Error') {
      color = '#990000'
    } else if(title == 'Completed') {
      color = '#009900'
    }
    return (
      <Typography variant="caption" style={{color}}>
        { title }
      </Typography>
    )
  }, [
    job,
  ])

  const isFilecoinPlus = useMemo(() => {
    const {
      inputs = [],
    } = job.Spec

    return Math.random() > 0.5
    let hasFilecoinPlus = true
    inputs.forEach((input) => {
      if(input.CID && FILECOIN_PLUS_CIDS.includes(input.CID)) {
        hasFilecoinPlus = true
      }
    })
    return hasFilecoinPlus
  }, [
    job,
  ])

  return (
    <Box
      component="div"
      sx={{
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'flex-start',
        ...sx
      }}
    >
      <div
        style={{
          minWidth: '70px',
        }}
      >
        { shardState }
      </div>
      <div
        style={{
          minWidth: '50px',
        }}
      >
        {
          isFilecoinPlus && (
            <FilPlus />
          )
        }
      </div>
      
    </Box>
  )
}

export default JobState
