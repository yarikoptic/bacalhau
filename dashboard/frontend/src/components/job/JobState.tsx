import React, { FC, useMemo } from 'react'
import { SxProps } from '@mui/system'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import {
  Job,
} from '../../types'

import FilPlus from './FilPlus'
import ShardState from './ShardState'

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
  withFilecoinPlus?: boolean,
  sx?: SxProps,
}> = ({
  job,
  withFilecoinPlus = true,
  sx = {},
}) => {
  const shardState = useMemo(() => {
    return getShardStateTitle(getJobShardState(job))
  }, [
    job,
  ])

  const isFilecoinPlus = useMemo(() => {
    if(!withFilecoinPlus) return false
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
    withFilecoinPlus,
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
        <ShardState
          state={ shardState }
        />
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
