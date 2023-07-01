import React from 'react'

import { Toast, ToastContainer } from 'react-bootstrap'

export default function ToastError({
  show,
  toggleShow,
  position,
  background,
  message,
}) {
  return (
    <ToastContainer className="p-3" position={position} style={{ zIndex: 1 }}>
      <Toast
        show={show}
        onClose={toggleShow}
        bg={background}
        delay={8000}
        autohide
      >
        <Toast.Header>
          <strong className="me-auto">ERC-721-Checks</strong>
        </Toast.Header>
        <Toast.Body>{message}</Toast.Body>
      </Toast>
    </ToastContainer>
  )
}
