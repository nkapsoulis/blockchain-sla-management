import React from 'react';
import './Loader.css';

function Loader() {
  return (
    <div className="centered">
      <div className="lds-ring hidden" id="loader">
        <div />
        <div />
        <div />
        <div />
      </div>
    </div>
  );
}

export default Loader;
