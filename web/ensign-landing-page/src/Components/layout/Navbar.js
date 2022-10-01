import React from 'react';

export default function Navbar() {
    return (
        <nav class="absolute pt-10 top-2 max-w-7xl flex text-white px-4 sm:px-6">
            <div>
                <h1 class="font-extralight text-white">Rotational Labs</h1>
            </div>
            <p class="md:pl-20 font-bold">Ensign</p>
            <div class="pl-40">
                <a href="https://rotational.io/#services" target="_blank" class="sm:pl-2 md:pl-5">Services</a>
                <a href="https://rotational.io/#blog" target="_blank" class="sm:pl-2 md:pl-5">Blog</a>
                <a href="https://rotational.io/#about" target="_blank" class="sm:pl-2 md:pl-5">About</a>
                <a href="https://rotational.io/#contact" target="_blank" class="sm:pl-2 md:pl-5">Contact Us</a>
            </div>
    </nav>
    )
}
