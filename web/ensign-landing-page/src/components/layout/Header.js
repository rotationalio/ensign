import React from 'react';
import Navbar from './Navbar';
import hero from '../../assets/hero.png'
import PageTitle from '../content/PageTitle';

export default function Header() {
    return (
        <header>
            <div>
                <Navbar />
                <PageTitle />
                <img
                src={hero}
                alt="An illustration with a sky blue backgroundm, 2 white clouds, 3 birds flying and, a red and white lighthouse in the corner with 3 sea otters standing at the top. " />

            </div>
        </header>
    )
}