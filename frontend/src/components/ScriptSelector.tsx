import React, { useContext, useState } from "react";
import {LoadScripts} from "../../wailsjs/go/bot/Bot"
import Selector, { SelectorLabel } from "./Selector";
import { ControlContext } from "./Controls";

export function ScriptSelectorLabel() {
    return <SelectorLabel labelText="Script" />
}

export default function ScriptSelector() {
    const controlContext = useContext(ControlContext);
    const name = "script-selector";
    const [scripts, setScripts] = useState<string[]>([]);
    const loadScripts = async () => {
        try {
            const loadedScripts = await LoadScripts()
            setScripts(loadedScripts);
            if (loadedScripts.length > 0) {
                controlContext.setSelectedScript(loadedScripts[0])
            }
        } catch (error) {
            console.error("Error while loading scripts: ", error)
        }
    };
    const handleOnChange = (selectedScript: string) => {
        if(selectedScript !== "") {
            controlContext.setSelectedScript(selectedScript);
        }
    }
    return <Selector name={name} scripts={scripts} loadOptions={loadScripts} handleOnChange={handleOnChange} />;
}