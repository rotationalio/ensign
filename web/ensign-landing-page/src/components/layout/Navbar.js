import React from 'react';

const shift = {
  marginTop: '-35px',
  marginRight: '180px',
}

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
        <nav className="pt-10 text-white">
            <div>
                <h1 className="pl-20 flex font-bold text-white text-3xl">Rotational Labs</h1>
            </div>
            <div className="topnav flex justify-end text-xl" id="myTopnav" style={shift}>
                <a href="https://rotational.app" className="sm:pl-2 md:pl-5 font-bold">Ensign</a>
                <a href="https://rotational.io/services" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">Services</a>
                <a href="https://rotational.io/about" target="_blank" rel="noreferrer" className="sm:pl-2 md:pl-5">About</a>
                <a href="javascript:void(0);" class="icon" onclick={myFunction}>
                    <i class="fa fa-bars"></i>
                </a>
            </div>
    </nav>
    )
}
