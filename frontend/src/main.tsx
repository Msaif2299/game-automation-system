//@ts-check
import React from 'react';
import ReactDOM from 'react-dom/client'; // For React 18
import Menu from './components/Menu'; // Your main App component
import './css/style.css'; // Any global styles
import Controls from './components/Controls';
import './css/glitch.scss';

// Correctly render the React app using StrictMode as a value
ReactDOM.createRoot(document.getElementById('app')!).render(
  <React.StrictMode>
    <Menu />
    <div id="wrapper" data-augmented-ui="tr-clip-x bl-2-clip-y inlay">
      <Controls />
      <div id="animator"></div>
    </div>
    <div id="drawer"></div>
  </React.StrictMode>
);
