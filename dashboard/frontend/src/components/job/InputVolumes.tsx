import React, { FC } from 'react'
import { SxProps } from '@mui/system'
import Box from '@mui/material/Box'
import {
  getShortId,
} from '../../utils/job'
import {
  StorageSpec,
} from '../../types'

const InputVolumes: FC<{
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
          let useUrl = ''
          if(storageSpec.URL) {
            const parts = storageSpec.URL.split(':')
            parts.pop()
            useUrl = parts.join(':')
          }
          else if(storageSpec.CID) {
            useUrl = `http://ipfs.io/ipfs/${storageSpec.CID}` 
          }
          return (
            <li key={storageSpec.CID || storageSpec.URL}>
              <a
                href={ useUrl }
                target="_blank"
                rel="noreferrer"
                style={{
                  fontSize: '0.8em',
                  color: '#333',
                }}
              >
                { getShortId(storageSpec.CID || storageSpec.URL || '', 16) }:{ storageSpec.path }
              </a>
            </li>
          )
        })
      }
    </Box>
  )
}

export default InputVolumes
