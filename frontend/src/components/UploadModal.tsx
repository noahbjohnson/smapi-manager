import * as React from 'react'
import {useState} from 'react'
import {Button, Modal, Uploader} from "rsuite";
import {FileType} from "rsuite/es/Uploader";


export default function (props: { show: boolean, hide: () => void }) {
    const [uploadQueue, setUploadQueue] = useState<FileType[]>([])

    const dismiss = (): void => {
        props.hide()
        setTimeout(() => setUploadQueue([]), 500)
    }

    const uploadFile = async (file: File): Promise<void> => {
        let formData = new FormData()
        formData.append("zip", file)
        await fetch('http://localhost:53494/upload', {method: "POST", body: formData});
    }

    const uploadFiles = async (): Promise<void> => {
        for (const file of uploadQueue) {
            const blobFile = file.blobFile
            if (blobFile) {
                await uploadFile(blobFile)
            }
        }
        dismiss()
    }

    return (
        <Modal show={props.show} onHide={dismiss}>
            <Modal.Header>
                <Modal.Title>Modal Title</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Uploader draggable accept='.zip'
                          multiple
                          autoUpload={false}
                          fileList={uploadQueue}
                          onChange={q => setUploadQueue(q)}
                          renderFileInfo={(file, fileElement) => {
                              return (
                                  <div>
                                      <p>File Name: {file.name}</p>
                                  </div>
                              );
                          }}>
                    <div>Click or Drag files to this area to upload</div>
                </Uploader>
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={uploadFiles} appearance="primary">
                    Ok
                </Button>
                <Button onClick={() => ""} appearance="subtle">
                    Cancel
                </Button>
            </Modal.Footer>
        </Modal>
    )
}