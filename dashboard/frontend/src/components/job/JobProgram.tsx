import React, { FC, useMemo } from 'react'
import { SxProps } from '@mui/system'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import {
  Job,
} from '../../types'

const JobProgram: FC<{
  job: Job,
  sx?: SxProps,
  imgSize?: number,
  fontSize?: string,
}> = ({
  job,
  sx = {},
  imgSize = 36,
  fontSize = '1em',
}) => {
  const engineLogo = useMemo(() => {
    if (job.Spec.Engine == "Docker") {
      return (
        <img
          style={{
            width: `${imgSize}px`,
            marginRight: '10px',
          }}
          src="/img/docker-logo.png" alt="Docker"
        />
      )
    } else if(job.Spec.Engine == "wasm") {
      return (
        <img
          style={{
            width: `${imgSize}px`,
            height: `${imgSize}px`,
          }}
          src="/img/wasm-logo.png" alt="WASM"
        />
      )
    }
  }, [
    job,
  ])

  const programDetails = useMemo(() => {
    if (job.Spec.Engine == "Docker") {
      const image = job.Spec.Docker?.Image || ''
      const entrypoint = job.Spec.Docker?.Entrypoint || []
      const details = `${image} ${(entrypoint || []).join(' ')}`
      return (
        <div>
          <div>
            <Typography variant="caption" style={{fontWeight: 'bold'}}>
              { image }
            </Typography>
          </div>
          <div>
            <Typography variant="caption" style={{color: '#666'}}>
              { (entrypoint || []).join(' ') }
            </Typography>
          </div>
        </div>
      )
    } else {
      return 'unknown'
    }
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
      <div>
        { engineLogo }
      </div>
      <div>
        { programDetails }
      </div>
      
      
    </Box>
  )
}

export default JobProgram
