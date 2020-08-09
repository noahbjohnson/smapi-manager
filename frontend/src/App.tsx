import React, {FunctionComponentElement, useEffect, useState} from 'react'
import './App.css'
import 'rsuite/dist/styles/rsuite-default.css';
import Splash from "./components/Splash";
import Landing from "./components/Landing";

interface bound {
    openSmapiInstall: () => Promise<void>
    hasSmapi: () => Promise<boolean>
}

export const BoundFunction = window as unknown as bound

function App() {
    let [showSplash, setShowSplash] = useState<boolean>(true)
    let [splashMessage, setSplashMessage] = useState<FunctionComponentElement<any>>(<h3>Loading...</h3>)

    useEffect(() => {
        BoundFunction.hasSmapi().then(smapiStatus => {
            if (smapiStatus) {
                setTimeout(() => setShowSplash(false), 1000)
            } else {
                BoundFunction.openSmapiInstall().catch(r => {
                    console.error(r)
                    setSplashMessage(<h3>Please <a
                        href={"https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI"}>install
                        SMAPI</a> and restart the application</h3>)
                }).then(() => {
                    setSplashMessage(<h3>Please switch to your browser or <a
                        href={"https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI"}>click
                        here</a> to install SMAPI, and restart the application</h3>)
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
                        : <Landing/>
                }

            </div>
        </div>
    )
}

export default App
