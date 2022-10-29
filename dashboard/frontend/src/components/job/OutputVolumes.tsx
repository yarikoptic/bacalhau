import React, { FC} from 'react'
import { SxProps } from '@mui/system'
import Box from '@mui/material/Box'
import {
  StorageSpec,
  RunCommandResult,
} from '../../types'

const OutputVolumes: FC<{
  outputVolumes: StorageSpec[],
  publishedResults?: StorageSpec,
  sx?: SxProps,
}> = ({
  outputVolumes = [],
  publishedResults,
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
        publishedResults && (
          <li>
            <span
              style={{
                fontSize: '0.8em',
                color: '#333',
              }}
            >
              <a href={ `https://ipfs.io/ipfs/${publishedResults.CID}` }>
                all
              </a>
            </span>
          </li>
        )
      }
      {
        outputVolumes.map((storageSpec) => {
          return (
            <li key={storageSpec.Name}>
              <span
                style={{
                  fontSize: '0.8em',
                  color: '#333',
                }}
              >
                {
                  publishedResults ? (
                    <a href={ `https://ipfs.io/ipfs/${publishedResults.CID}${storageSpec.path}` }>
                      { storageSpec.Name }:{ storageSpec.path }
                    </a>
                  ) : (
                    <span>
                      { storageSpec.Name }:{ storageSpec.path }
                    </span>
                  )
                }
                
              </span>
            </li>
          )
        })
      }
    </Box>
  )
}

export default OutputVolumes
