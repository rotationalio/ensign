import React from 'react';
import hero from './img/hero.png'
import PageTitle from '../content/PageTitle';

export default function Header() {
    return (
        <header class="pb-20">
            <div class="relative">
                <img
                src={hero}
                alt="An illustration with a sky blue backgroundm, 2 white clouds, 3 birds flying and, a red and white lighthouse in the corner with 3 sea otters standing at the top. " />
                <PageTitle />
            </div>
            
        </header>
    )
}