import React from "react";
import AccessForm from "../content/AccessForm";

export default function Main() {
    return(
        <main className="container pt-8">
            <section className="xl:grid grid-cols-3 gap-8">
                <section className="col-span-2">
                    <h2 className="leading-10 pb-5">Powering Real-time Apps & Analytics</h2>
                    <p className="pb-5">The global real-time revolution is here.</p>
                    <p className="pb-5">
                        <span className="text-[#1D65A6] font-bold">Ensign</span> is a <span className="font-bold">next-generation distributed event store and stream-processing platform</span> for real-time apps and analytics that requires no additional investment in infrastructure or overhead. <span className="font-bold">Built for speed and simplicity</span>, Ensign offers advanced features for next-generation applications while significantly reducing barriers to building, maintaining, and scaling event-driven applications. With Ensign, you can:
                    </p>
                    <ul className="list-disc list-outside pb-5 pl-10">
                        <li>Customize your data pipelines</li>
                        <li>Quickly build or integrate events into new or existing applications</li>
                        <li>Provide fast, consistent, and personalized digital experiences across time and space</li>
                        <li>Accelerate time-to-insight in business intelligence and data analytics</li>
                    </ul>
                    <p className="pb-5">
                        Designed as a “low ops / no ops” cloud-agnostic <span className="font-bold">managed service</span>, Ensign is ideal for developers and organizations building <span className="font-bold">event-driven microservices</span> to power rich consumer digital experiences, streaming machine learning models, and real-time business intelligence dashboards.
                    </p>
                </section>
                <section className="pb-8">
                    <AccessForm />
                </section>
            </section>
        </main>
    )
}

