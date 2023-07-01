import React, { useRef, useState } from 'react'

import { ListGroup } from 'react-bootstrap'

import CheckIcon from '../icons/CheckIcon'
import ClipboardIcon from '../icons/ClipboardIcon'

const CHECK_ICON_SHOW_TIME = 5000

function ListItem({ uploadedFile, index }) {
  const buttonRef = useRef(null)
  const [copied, setCopied] = useState(false)

  const handleCopyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText('ipfs://' + uploadedFile.cid)
      setCopied(true)

      setTimeout(() => {
        setCopied(false)
      }, CHECK_ICON_SHOW_TIME)
    } catch (error) {
      console.error('Error copying to clipboard:', error)
    }
  }

  return (
    <ListGroup.Item
      className="d-flex align-items-center my-2 rounded border shadow-sm"
      ref={buttonRef}
    >
      <div className="d-flex w-100 font-monospace align-items-center">
        <span>{index + 1}.</span>
        <div className="flex-grow-1 overflow-hidden ms-2">
          <div className="w-100 fw-bold text-truncate">{uploadedFile.name}</div>
          <div className="w-100 text-truncate">ipfs://{uploadedFile.cid}</div>
        </div>
        <button
          className="btn btn-link flex-grow-2"
          onClick={handleCopyToClipboard}
        >
          {copied ? <CheckIcon /> : <ClipboardIcon />}
        </button>
      </div>
    </ListGroup.Item>
  )
}

export default ListItem
