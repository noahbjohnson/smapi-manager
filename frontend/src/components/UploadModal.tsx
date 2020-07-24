import * as React from 'react'
import {Button, Modal, Uploader} from "rsuite";
import {useRef} from "react";

const uploaderComponent = <Uploader draggable accept='.zip'
                                    multiple
                                    autoUpload={false}
                                    renderFileInfo={(file, fileElement) => {
                                        return (
                                            <div>
                                                <p>File Name: {file.name}</p>
                                                <p>File status: {file.status}</p>
                                            </div>
                                        );
                                    }}>
    <div>Click or Drag files to this area to upload</div>
</Uploader>

export default function (props: { show: boolean, hide: ()=>void}) {
    const uploader = useRef(uploaderComponent)

    return (
      <Modal show={props.show} onHide={props.hide}>
          <Modal.Header>
              <Modal.Title>Modal Title</Modal.Title>
          </Modal.Header>
          <Modal.Body>
              {uploader.current}
          </Modal.Body>
          <Modal.Footer>
              <Button onClick={()=>""} appearance="primary">
                  Ok
              </Button>
              <Button onClick={()=>""} appearance="subtle">
                  Cancel
              </Button>
          </Modal.Footer>
      </Modal>
  )
}