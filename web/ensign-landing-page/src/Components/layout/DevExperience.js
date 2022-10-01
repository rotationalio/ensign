import React from 'react'

export default function DevExperience() {
    return (
        <section class="bg-[#ECF6FF] min-w-min">
            <div class="w-9/12 mx-auto">
                <h2 class="sm:leading-loose py-5 text-left">The Developer Experience</h2>
                <p class="text-left pb-5">Ensign is an advanced event data store designed with application developers, data scientists, and product managers in mind. Ensign combines fast transactional services with decoupled processing and rich, insight-driven online analysis, without the need for additional infrastructure or a PhD. in Kafka. Ensign makes event-driven microservices accessible to everyday developers, data scientists, and product managers.</p>
                <ul class="list-disc list-inside text-left">
                    <li>Create an account</li>
                    <li>Connect your data sources via our secure API</li>
                    <li>Set up publishers and consumers</li>
                    <li>Write sets of “rules” for the Ensign Event Broker to route, store, and/or transform data while in motion via our SDK</li>
                    <li>Integrate with your app, model, or dashboard</li>
                </ul>
            </div>
        </section>
    )
}

