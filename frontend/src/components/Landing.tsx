import * as React from 'react'
import {useState} from 'react'
import {Button} from "rsuite";
import UploadModal from "./UploadModal";

export default function () {
    let [showUploadModal, setShowUploadModal] = useState<boolean>(false)
    return (
        <div className="App">
            <h1>Loaded</h1>
            <Button onClick={() => setShowUploadModal(true)}>Add Mod Files</Button>
            <UploadModal show={showUploadModal} hide={() => setShowUploadModal(false)}/>
        </div>
    )
}


