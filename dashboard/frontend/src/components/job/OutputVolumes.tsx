import React, { FC} from 'react'
import { SxProps } from '@mui/system'
import Box from '@mui/material/Box'
import {
  StorageSpec,
} from '../../types'

const OutputVolumes: FC<{
  storageSpecs: StorageSpec[],
  sx?: SxProps,
}> = ({
  storageSpecs,
  sx = {},
}) => {
  return (
    <Box
      component="div"
      sx={{
        width: '100%',
        ...sx
      }}
    >
      {
        storageSpecs.map((storageSpec) => {
          return (
            <li key={storageSpec.Name}>
              <a
                href={ "" }
                target="_blank"
                rel="noreferrer"
                style={{
                  fontSize: '0.8em',
                  color: '#333',
                }}
              >
                { storageSpec.Name }:{ storageSpec.path }
              </a>
            </li>
          )
        })
      }
    </Box>
  )
}

export default OutputVolumes
