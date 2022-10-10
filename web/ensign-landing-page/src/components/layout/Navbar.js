import React from 'react';
import logo from '../layout/img/logo.png'

export default function Navbar() {

/* Toggle between adding and removing the "responsive" class to topnav when the user clicks on the icon */
function myFunction() {
    var x = document.getElementById("myTopnav");
    if (x.className === "topnav") {
      x.className += " responsive";
    } else {
      x.className = "topnav";
    }
  }

    return (
        <nav className="relative max-w-7xl flex items-center justify-between text-white">
          <a href="https://rotational.io" target="_blank" rel="noreferrer" className="pt-5">
            <img 
            src={logo}
            className="pl-20 h-14 w-auto sm:h-14" />
          </a>
          <ul className="topnav flex text-xl justify-end" id="myTopnav">
                <li><a href="https://rotational.app" className="sm:pl-2 md:pl-5 font-bold">Ensign</a></li>
                <li><a href="https://rotational.io/services/" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">Services</a></li>
                <li><a href="https://rotational.io/about" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">About</a></li>
                <li><a href="javascript:void(0);" class="icon" onclick={myFunction}>
                    <i class="fa fa-bars"></i>
                </a></li>
          </ul>
    </nav>
    )
}
