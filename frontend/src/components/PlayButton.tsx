import React, { useContext, useState } from "react";
import "../css/controls.css";
import { ControlContext } from "./Controls";
import {StartBot, StopBot} from "../../wailsjs/go/bot/Bot"
const playIcon = <>&#9658;</>
const pauseIcon = <>&#10074;&#10074;</>
export default function PlayButton () {
    const controlContext = useContext(ControlContext)
    const [playing, setPlaying] = useState(false)
    const changeState = async () => {
        if (!playing) {
            console.log("PLAYING")
            try {
                console.log(`SELECTED CLASS: ${controlContext.selectedClass}, SELECTED SCRIPT: ${controlContext.selectedScript}`)
                await StartBot(controlContext.selectedClass, controlContext.selectedScript)
                setPlaying(true)
            } catch (error) {
                setPlaying(false)
                console.log("Backend error: ", error)
                alert(`Error encountered: ${error}`)
            }
            return
        }
        console.log("PAUSING")
        StopBot()
        setPlaying(false)
    }
    return <button id="play-button"
        onClick={changeState}>
            {playing ? pauseIcon : playIcon}
        </button>
}
