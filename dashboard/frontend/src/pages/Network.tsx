import React, { FC } from 'react'
import Grid from '@mui/material/Grid'
import Container from '@mui/material/Container'

const Network: FC = () => {

  return (
    <Container maxWidth={ 'xl' } sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          Network
        </Grid>
      </Grid>
    </Container>
  )
}

export default Network