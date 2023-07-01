import React from 'react'

import { ListGroup } from 'react-bootstrap'

import ListItem from './ListItem'

function FileList({ uploadedFiles }) {
  return (
    <>
      {uploadedFiles.length > 0 ? (
        <>
          <h5 className="text-center">Your recently pinned files</h5>
          <ListGroup>
            {uploadedFiles.map((uploadedFile, index) => (
              <ListItem key={index} uploadedFile={uploadedFile} index={index} />
            ))}
          </ListGroup>
        </>
      ) : (
        <h5 className="text-center">
          Your list is empty, try to pin a new file
        </h5>
      )}
    </>
  )
}

export default FileList
