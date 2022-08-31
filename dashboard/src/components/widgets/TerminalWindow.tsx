import React, { FC } from 'react'
import Window, { WindowProps } from './Window'
import TerminalText from './TerminalText'

interface TerminalWindowProps extends WindowProps {
  data: any,
  title?: string,
  onClose: {
    (): void,
  }
}

const TerminalWindow: FC<TerminalWindowProps> = ({
  data,
  title = 'kubectl get all',
  onClose,
  ...windowProps
}) => {
  return (
    <Window
      withCancel
      compact
      title={ title }
      onCancel={ onClose }
      cancelTitle="Close"
      {...windowProps}
    >
      <TerminalText
        data={ data }
      />
    </Window>
  )
}

export default TerminalWindow
