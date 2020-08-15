import React, {FunctionComponentElement, useEffect, useState} from 'react'
import './App.css'
import 'rsuite/dist/styles/rsuite-default.css';
import Splash from "./components/Splash";
import Landing from "./components/Landing";
import {getMods, Mod} from "./api";

interface bound {
    openSmapiInstall: () => Promise<void>
    hasSmapi: () => Promise<boolean>
    loadMods: () => Promise<Mod[]>
}

/**
 * Call this to have a (probably) type safe way to call bound methods
 */
export function BoundFunction() {
    // @ts-ignore
    return window.backend as unknown as bound
}


function App() {
    let [showSplash, setShowSplash] = useState<boolean>(true)
    let [splashMessage, setSplashMessage] = useState<FunctionComponentElement<any>>(<h3>Loading...</h3>)
    let [mods, setMods] = useState<Mod[]>([])

    function installMessage(openedBrowser = false) {
        setSplashMessage(<h3>Please {
            openedBrowser
                ? <a href={"https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI"}>install
                    SMAPI</a>
                : <span> switch to your browser or <a
                    href={"https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI"}>click
                    here</a> to install SMAPI</span>
        }
            and restart the application</h3>)
    }

    useEffect(() => {
        BoundFunction().hasSmapi().then(smapiStatus => {
            if (smapiStatus) {
                getMods().then(m => setMods(m))
                setTimeout(() => setShowSplash(false), 1000)
            } else {
                BoundFunction().openSmapiInstall().catch(r => {
                    console.error(r)
                    installMessage()
                }).then(() => {
                    installMessage(true)
                })
            }
        })
    }, [])


    return (
        <div id="app" className="App">
            <div className="App-body">
                {
                    showSplash
                        ? <Splash message={splashMessage}/>
                        : <Landing mods={mods}/>
                }

            </div>
        </div>
    )
}

export default App
