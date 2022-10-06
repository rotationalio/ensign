import React from 'react';
import diagram from './img/ensign-diagram.png'

export default function Diagram() {
    return (
        <section className="mx-auto container">
          <h2 className="">How It Works</h2>
          <img className="mx-auto" src={diagram} alt="" />
        </section>
    )
}