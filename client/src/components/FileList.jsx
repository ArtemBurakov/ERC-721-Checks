import React, { useRef, useState } from 'react'

import { ListGroup } from 'react-bootstrap'

import CheckIcon from './CheckIcon'
import ClipboardIcon from './ClipboardIcon'

const CHECK_ICON_SHOW_TIME = 5000

function FileList({ uploadedFiles }) {
  const buttonRefs = useRef([])
  const [copiedIndexes, setCopiedIndexes] = useState([])

  const handleCopyToClipboard = async (cid, index) => {
    try {
      await navigator.clipboard.writeText(cid)
      setCopiedIndexes((prevIndexes) => [...prevIndexes, index])

      setTimeout(() => {
        setCopiedIndexes((prevIndexes) =>
          prevIndexes.filter((i) => i !== index)
        )
      }, CHECK_ICON_SHOW_TIME)
    } catch (error) {
      console.error('Error copying to clipboard:', error)
    }
  }

  return (
    <>
      {uploadedFiles.length > 0 ? (
        <>
          <h5 className="text-center">Your recently pinned files</h5>
          <ListGroup>
            {uploadedFiles.map((uploadedFile, index) => (
              <ListGroup.Item
                key={index}
                className="d-flex align-items-center my-2 rounded border shadow-sm"
                ref={(ref) => (buttonRefs.current[index] = ref)}
              >
                <div className="d-flex w-100 font-monospace align-items-center">
                  <span>{index + 1}.</span>
                  <div className="flex-grow-1 overflow-hidden ms-2">
                    <div className="w-100 fw-bold text-truncate">
                      {uploadedFile.name}
                    </div>
                    <div class="w-100 text-truncate">
                      ipfs://{uploadedFile.cid}
                    </div>
                  </div>
                  <button
                    className="btn btn-link flex-grow-2"
                    onClick={() =>
                      handleCopyToClipboard(uploadedFile.cid, index)
                    }
                  >
                    {copiedIndexes.includes(index) ? (
                      <CheckIcon />
                    ) : (
                      <ClipboardIcon />
                    )}
                  </button>
                </div>
              </ListGroup.Item>
            ))}
          </ListGroup>
        </>
      ) : (
        <h5 className="text-center">
          You list is empty, try to pin a new file
        </h5>
      )}
    </>
  )
}

export default FileList
