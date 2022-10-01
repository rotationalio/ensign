import React from 'react';
import Header from '../layout/Header';
import Footer from '../layout/Footer'
import seaotter from '../layout/img/dunking-seaotter.png'
import finger from '../layout/img/finger.png'
import twitter from '../layout/img/twitter.png'

export default function Confirmation() {
    return (
        <>
            <Header />
            <section className="flex">
                <section className="ml-20">
                    <h1 className="pb-5">Success! Thank you for your interest in Ensign.</h1>
                    <p className="pb-5 font-bold">We'll be in touch soon.</p>
                    <p className="pb-3">Now you might be thinking: Why sea otters?</p>
                    <div className="flex">
                        <p className="pb-10">Oh, don't get us started!</p>
                        <img
                            src={finger}
                            alt="A pointer finger outlined in blue."
                            className="pl-5" />
                    </div>

                    <p className="pb-5 font-bold">What next?</p>
                    <ul className="pb-5">
                        <li className="pb-3">Expect a confirmation email from us.</li>
                        <div className="flex">
                            <li className="pb-3">Twitter about us (or sea otters)?</li>
                            <img
                            src={twitter}
                            alt="Twitter logo, a white bird with a blue background."
                            class />
                        </div>
                        <li className="pb-3">Teach your kids (or friends) about streaming</li>
                        <li className="pb-3">Dream about your first event stream.</li>
                    </ul>
                    <p>Probably you'll just wait to hear from us.</p>
                </section>
                <section className="ml-20">
                    <img
                    src={seaotter}
                    alt="A seaotter in a pool playing basketball."/>
                </section>
            </section>
            <Footer />
        </>
    )
}