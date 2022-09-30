import React from 'react';

export default function Navbar() {
    return (
        <nav class="relative max-w-7xl mx-auto flex items-center justify-between text-white px-4 sm:px-6"
        aria-label="Global">
        <a href="" class="pt-3 sm:pt-0">
            <img src="img/logo.png" alt="Rotational" class="h-14 w-auto sm:h-14" />
        </a>
        <div class="topnav" id="myTopnav">
            <a href="#" class="active">Services</a>
            <a href="#">Open Source</a>
            <a href="#">Blog</a>
            <a href="#">About</a>
            <a href="#">Contact Us</a>
            <a href="javascript:void(0);" class="icon" onclick="myFunction()">
                <i class="fa fa-bars"></i>
            </a>
        </div>
    </nav>
    )
}