import React, {FunctionComponentElement, useEffect, useState} from 'react'
import './App.css'
import 'rsuite/dist/styles/rsuite-default.css';
import Splash from "./components/Splash";
import Landing from "./components/Landing";

async function HasSMAPI(): Promise<boolean> {
    return true
}

function App() {
    let [showSplash, setShowSplash] = useState<boolean>(true)
    let [splashMessage, setSplashMessage] = useState<FunctionComponentElement<any>>(<h3>Loading...</h3>)

    useEffect(() => {
        HasSMAPI().then(hasSmapi => {
            if (hasSmapi) {
                setTimeout(() => setShowSplash(false), 1000)
            } else {
                setSplashMessage(<h3>"Please <a
                    href={"https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI"}>install
                    SMAPI</a> and restart the application"</h3>)
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
