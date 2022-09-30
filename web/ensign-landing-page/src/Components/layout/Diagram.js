import React from 'react';
import diagram from './img/ensign-diagram.jpeg'

export default function Diagram() {
    return (
        <section>
                <h2 class="ml-20">How It Works</h2>
                <img
                class="max-w-7xl mx-auto w-screen"
                src={diagram}
                alt="" />
        </section>
    )
}