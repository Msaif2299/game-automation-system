import React, { createContext, Dispatch, SetStateAction, useState } from "react";
import PlayButton from "./PlayButton";
import ScriptSelector, { ScriptSelectorLabel } from "./ScriptSelector";
import ClassSelector, { ClassSelectorLabel } from "./ClassSelector";

type ControlContextType = {
    selectedScript: string
    setSelectedScript: Dispatch<SetStateAction<string>>
    selectedClass: string
    setSelectedClass: Dispatch<SetStateAction<string>>
}

export const ControlContext = createContext<ControlContextType>({
    selectedScript: "",
    selectedClass: "",
    setSelectedScript: () => {console.log("Unimplemented setSelectedScript error")},
    setSelectedClass: () => {console.log("Unimplemented setSelectedClass error")}
}) 

export default function Controls() {
    const [selectedScript, setSelectedScript] = useState("");
    const [selectedClass, setSelectedClass] = useState("");
    return <div id="control-panel" 
    style={{display: 'flex', justifyContent: 'center', alignItems: 'center', height:'100%'}}>
        <ControlContext.Provider value={{selectedScript, setSelectedScript, selectedClass, setSelectedClass}}>
            <div id="label-box" style={{float: "left", margin: "5px", justifyContent: "center", alignItems: "center"}}>
                <ScriptSelectorLabel /><br />
                <ClassSelectorLabel />
            </div>
        <div style={{float: 'left', justifyContent: "center", alignItems: "center", marginRight: "20px"}}>
            <ScriptSelector /> <ClassSelector />
        </div><div><PlayButton /></div>
        </ControlContext.Provider>
    </div>
}
