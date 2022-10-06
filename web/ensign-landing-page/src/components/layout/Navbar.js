import React from 'react';

const shift = {
  marginTop: '-35px',
  marginRight: '180px',
}

export default function Navbar() {
    return (
        <nav className="pt-10 text-white">
            <div>
                <h1 className="pl-20 flex font-bold text-white text-3xl">Rotational Labs</h1>
            </div>
            <div className="flex justify-end text-xl" style={shift}>
                <a href="https://rotational.app" className="sm:pl-2 md:pl-5 font-bold">Ensign</a>
                <a href="https://rotational.io/services" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">Services</a>
                <a href="https://rotational.io/about" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">About</a>
                <a href="https://rotational.io/contact" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">Contact Us</a>
            </div>
    </nav>
    )
}
