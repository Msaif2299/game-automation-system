import React from "react";
import { Quit, WindowMinimise } from '../../wailsjs/runtime/runtime';
import '../css/menu.css'

class MaximizeButton extends React.Component {
    render(): React.ReactNode {
        return <button 
            id="maximize-button" 
            className="menu-button" 
            type="button"
            onClick={() => {}}
            disabled={true}
        >&#128913;</button>
    }
}

class MinimizeButton extends React.Component {
    render(): React.ReactNode {
        return <button 
            id="minimize-button" 
            className="menu-button" 
            type="button"
            onClick={WindowMinimise}
        >&#9644;</button>
    }
}

class CloseButton extends React.Component {
    render(): React.ReactNode {
        return <button 
        id="close-button" 
        className="menu-button" 
        type="button"
        onClick={Quit}
        >X</button>
    }
}

function Menu() {
    return (
        <div id="menu" data-augmented-ui="bl-2-clip-y">
            {/* <div data-augmented-ui-reset> */}
                <MaximizeButton />
                <MinimizeButton />
                <CloseButton />
            {/* </div> */}
        </div>
    )
}
export default Menu