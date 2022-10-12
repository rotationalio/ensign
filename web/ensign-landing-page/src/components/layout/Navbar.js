import React from 'react';
import toggleResponsiveClass from './ResponsiveNav'
import logo from '../layout/img/logo.png'

export default function Navbar() {
    return (
        <nav className="relative max-w-7xl flex items-center justify-between text-white">
          <a href="https://rotational.io" target="_blank" rel="noreferrer" className="pt-3">
            <img 
              src={logo}
              alt="Rotational Labs logo"
              className="pl-10 h-14 w-auto sm:pl-20 h-14" 
            />
          </a>
          <ul className="topnav" id="myTopnav">
            <li><a href="https://rotational.app">Ensign</a></li>
            <li><a href="https://rotational.io/services/" target="_blank" rel="noreferrer">Services</a></li>
            <li><a href="https://rotational.io/about" target="_blank" rel="noreferrer">About</a></li>
            <li><a href="javascript:void(0);" class="icon" onClick={toggleResponsiveClass}>
                <i class="fa fa-bars"></i>
            </a></li>
          </ul>
        </nav>
    )
}
