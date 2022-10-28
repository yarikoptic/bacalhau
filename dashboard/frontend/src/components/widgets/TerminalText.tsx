import React, { FC } from 'react'
import Box from '@mui/material/Box'

const TerminalText: FC<{
  data: any,
}> = ({
  data,
}) => {
  return (
    <Box
      component="div"
      sx={{
        width: '100%',
        padding: 2,
        margin: 0,
        backgroundColor: '#000000',
        overflow: 'auto',
      }}
    >
      <Box
        component="pre"
        sx={{
          padding: 1,
          margin: 0,
          color: '#ffffff',
          font: 'Courier',
          fontSize: '12px',
        }}
      >
        { typeof(data) === 'string' ? data : JSON.stringify(data, null, 4) }
      </Box>
    </Box>
  )
}

export default TerminalText
