import React, { useContext, useState } from "react";
import {LoadClasses} from "../../wailsjs/go/bot/Bot"
import Selector, { SelectorLabel } from "./Selector";
import { ControlContext } from "./Controls";

export function ClassSelectorLabel() {
    return <SelectorLabel labelText="Class" />
}

export default function ClassSelector() {
    const controlContext = useContext(ControlContext);
    const name = "class-selector";
    const [scripts, setScripts] = useState<string[]>([]);
    const loadScripts = async () => {
        try {
            const loadedClasses = await LoadClasses()
            setScripts(loadedClasses);
            if (loadedClasses.length > 0) {
                controlContext.setSelectedClass(loadedClasses[0])
            }
        } catch (error) {
            console.error("Error while loading scripts: ", error)
        }
    };
    const handleOnChange = (selectedClass: string) => {
        if(selectedClass !== "") {
            controlContext.setSelectedClass(selectedClass);
        }
    }
    return <Selector name={name} scripts={scripts} loadOptions={loadScripts} handleOnChange={handleOnChange} />;
}