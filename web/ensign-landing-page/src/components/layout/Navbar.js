import React from 'react';

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
        <nav className="pt-10 flex items-center text-white">
          <h1 className="pl-20 font-bold text-white text-3xl">Rotational Labs</h1>
          <ul className="topnav flex text-xl" id="myTopnav">
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
