import {BoundFunction} from "./App";

const API_ROOT = "http://localhost:53494/"
const UPLOAD_PATH = API_ROOT + 'upload'

export const uploadFile = async (file: File): Promise<void> => {
    const formData = new FormData()
    formData.append('zip', file)
    await fetch(UPLOAD_PATH, {method: 'POST', body: formData})
}

export const getMods = async (): Promise<Mod[]> => {
    return BoundFunction().loadMods()
}

export interface Mod {
    directory: string
    enabled: boolean
    managed: boolean
    metadata: ModMetadata
}

interface ModMetadata {
    Name: string
    Author: string
    Version: string
    Description: string
    UniqueID: string
    EntryDll: string
    ContentPackFor: ContentPackRef
    MinimumApiVersion: string
    Dependencies: DependencyRef[]
    UpdateKeys: string[]
}

interface DependencyRef {
    UniqueID: string
    MinimumVersion: string
    IsRequired: boolean
}

interface ContentPackRef {
    UniqueID: string
    MinimumVersion: string
}