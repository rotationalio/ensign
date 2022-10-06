import React from "react";
import Access from "../content/Access";

export default function Main() {
    return(
        <main className="pb-20 md:ml-20">
            <section className="pb-10">
                <h2 className="leading-10">Data Engineering Simplified.</h2>
                <h2>Navigate Your Data to Where It's Valued Most.</h2>
            </section>
            <section className="sm:flex ">
                <section className="text-left md:w-1/2">
                    <p className="pb-5">Introducing <span className="text-[#1D65A6] font-bold">Ensign</span>, our <span className="font-bold">intelligent event-data platform</span> for real-time apps and analytics. Combine different data sources and deliver value in real-time with no additional infrastructure or admin burden.</p>
                    <p className="pb-5">Ensign is the first eventing platform that focuses on developer experience, with an API that makes eventing click for everyday tech folks. Ensign makes it easy to:</p>
                    <ul className="list-disc list-inside pb-5">
                        <li>Customize your data pipelines</li>
                        <li>Quickly build or integrate events into new or existing applications</li>
                        <li>Provide fast, consistent, and personalized digital experiences across time and space</li>
                        <li>Accelerate time-to-insight in business intelligence and data analytics</li>
                    </ul>
                    <p>Even better, Ensign grows with you with built-in geo-scaling, data compliance, and diasaster recovery controls.</p>
                </section>
                <section className="sm:mx-auto ml-20 max-h-screen bg-[#DED6C5]">
                    <Access />
                </section>
            </section>
        </main>
    )
}

