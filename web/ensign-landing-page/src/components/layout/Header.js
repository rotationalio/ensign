import React from 'react';
import Navbar from './Navbar';
import PageTitle from '../content/PageTitle';

export default function Header() {
    return (
        <div className="bg-wave-pattern bg-repeat-x bg-right min-h-[480px] w-screen">
          <header className="bg-hero bg-no-repeat bg-right bg-[length:882.5px_480px] min-h-[480px] w-screen">
            <Navbar />
            <PageTitle />
          </header>
        </div>
    )
}