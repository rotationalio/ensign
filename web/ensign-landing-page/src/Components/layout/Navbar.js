import React from 'react';

export default function Navbar() {
    return (
        <nav className="absolute pt-10 top-2 max-w-7xl flex text-white px-4 sm:px-6">
            <div>
                <h1 className="font-extralight text-white">Rotational Labs</h1>
            </div>
            <p className="md:pl-20 font-bold">Ensign</p>
            <div className="pl-40">
                <a href="https://rotational.io/#services" target="_blank" className="sm:pl-2 md:pl-5">Services</a>
                <a href="https://rotational.io/#blog" target="_blank" className="sm:pl-2 md:pl-5">Blog</a>
                <a href="https://rotational.io/#about" target="_blank" className="sm:pl-2 md:pl-5">About</a>
                <a href="https://rotational.io/#contact" target="_blank" className="sm:pl-2 md:pl-5">Contact Us</a>
            </div>
    </nav>
    )
}
