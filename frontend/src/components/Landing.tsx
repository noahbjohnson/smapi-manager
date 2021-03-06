import * as React from 'react'
import { Button, ButtonToolbar, Container, Content, Drawer, Header } from 'rsuite'
import Dropzone from 'react-dropzone'
import { Mod, uploadFile } from '../api'

export default function (props: {mods: Mod[]}) {
  const [uploadQueue, setUploadQueue] = React.useState<File[]>([])
  const [showDrawer, setDrawer] = React.useState(false)

  const uploadFiles = async (): Promise<void> => {
    await Promise.all(uploadQueue.map(file => {
      return uploadFile(file).then(() => setUploadQueue(uploadQueue.filter((f: File): boolean => {
        return file.name !== f.name
      })))
    }))
    hide()
  }

  const addOne = (file: File): void => {
    if (uploadQueue.map(f => f.name).includes(file.name)) {
      console.warn('duplicate file detected')
    } else {
      uploadQueue.push(file)
      setUploadQueue(uploadQueue.map(f => f))
    }
  }

  const addToQueue = (files: File | File[]): void => {
    if (Array.isArray(files)) {
      files.forEach(addOne)
    } else {
      addOne(files)
    }
    setUploadQueue(uploadQueue)
  }

  const hide = () => {
    setTimeout(() => setUploadQueue([]), 500)
  }

  return (
    <Container>
      <Header>
        <ButtonToolbar>
          <Button onClick={() => setDrawer(!showDrawer)}>Add Mods</Button>
        </ButtonToolbar>
      </Header>
      <Content>
        <div>
          {props.mods.map((mod, n) => {
            return <p key={n}>{mod.metadata.Name}</p>
          })}
        </div>
      </Content>

      <Drawer
        show={showDrawer}
        onHide={() => setDrawer(false)}
        backdrop={'static'}
      >
        <Drawer.Header>
          <Drawer.Title>Add Mod Files</Drawer.Title>
        </Drawer.Header>
        <Dropzone
          onDrop={
            (acceptedFiles) => {
              addToQueue(acceptedFiles)
            }
          }
          accept={'.zip'}>
          {({ getRootProps, getInputProps }) => (
            <section>
              <div {...getRootProps()}>
                <input {...getInputProps()} />
                <h1>Drag 'n' drop some zip files here, or click to select.</h1>
              </div>
            </section>
          )}
        </Dropzone>
        {
          uploadQueue.map((file, key) => {
            return <p>{file.name}</p>
          })
        }
        <Drawer.Footer>
          <Button onClick={uploadFiles} appearance="primary">Add Selected Mods</Button>
          <Button onClick={() => setDrawer(false)} appearance="subtle">Cancel</Button>
        </Drawer.Footer>
      </Drawer>
    </Container>
  )
}
