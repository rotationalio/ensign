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
            <section class="flex">
                <section class="ml-20">
                    <h1 class="pb-5">Success! Thank you for your interest in Ensign.</h1>
                    <p class="pb-5 font-bold">We'll be in touch soon.</p>
                    <p class="pb-3">Now you might be thinking: Why sea otters?</p>
                    <div class="flex">
                        <p class="pb-10">Oh, don't get us started!</p>
                        <img 
                            src={finger}
                            alt="A pointer finger outlined in blue."
                            class="pl-5" />
                    </div>
                    
                    <p class="pb-5 font-bold">What next?</p>
                    <ul class="pb-5">
                        <li class="pb-3">Expect a confirmation email from us.</li>
                        <div class="flex">
                            <li class="pb-3">Twitter about us (or sea otters)?</li>
                            <img 
                            src={twitter}
                            alt="Twitter logo, a white bird with a blue background."
                            class />
                        </div>
                        <li class="pb-3">Teach your kids (or friends) about streaming</li>
                        <li class="pb-3">Dream about your first event stream.</li>
                    </ul>
                    <p>Probably you'll just wait to hear from us.</p>
                </section>
                <section class="ml-20">
                    <img 
                    src={seaotter}
                    alt="A seaotter in a pool playing basketball."/>
                </section>
            </section>
            <Footer />
        </>
    )
}