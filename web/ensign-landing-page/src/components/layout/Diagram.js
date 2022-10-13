import React from 'react';
import diagram from './img/ensign-diagram.png'

export default function Diagram() {
    return (
        <section className="mx-auto container pb-10">
          <h2>Data Engineering Simplified</h2>
          <img 
            className="mx-auto" 
            src={diagram} 
            alt="Illustration of how Ensign works. Text in diagram reads: 1. Connect data sources. 2. Enrich, Combine, Move Data on Ensign. 3. Deliver real-time events and analytics seamlessly." 
            />
        </section>
    )
}