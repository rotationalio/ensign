import React from 'react';
import Header from '../layout/Header';
import Footer from '../layout/Footer';
import seaotter from '../layout/img/dunking-seaotter.gif';
import finger from '../layout/img/finger.png';
import withTracker from '../../lib/analytics';

const Confirmation = () => {
    return (
        <>
            <Header />
            <section className="mx-auto container pt-8 pb-10">
              <div className="grid grid-cols-3 gap-8">
                <section className="col-span-2">
                    <h2 className="leading-10">Success! Thank you for your interest in Ensign.</h2>
                    <h3 className="pb-5 font-bold">We'll be in touch soon.</h3>
                    <p className="pb-3">Now you might be thinking: Why sea otters?</p>
                    <div className="flex">
                        <p className="pb-10">Oh, don't get us started!</p>
                        <img
                            src={finger}
                            alt="A pointer finger outlined in blue."
                            className="pl-5 flex h-8"
                        />
                    </div>
                    <h3 className="pb-5 font-bold">What next?</h3>
                    <ul className="pb-5 list-disc list-outside">
                        <li className="pb-1">Expect a confirmation email from us.</li>
                        <li className="pb-1">Tweet about <a href="https://twitter.com/rotationalio" className="text-[#1D66A6]">us</a> (or <a href="https://twitter.com/in_otter_news2" className="text-[#1D66A6]">sea otters</a>)?</li>
                        <li className="pb-1"><a href="https://www.gentlydownthe.stream/" className="text-[#1D66A6]">Teach your kids (or friends) about streaming</a></li>
                        <li className="pb-1">Dream about your first event stream.</li>
                    </ul>
                    <p>Probably you'll just wait to hear from us.</p>
                </section>
                <section>
                    <img
                    src={seaotter}
                    alt="A sea otter in a pool playing basketball."/>
                </section>
              </div>
            </section>
            <Footer />
        </>
    )
}

export default withTracker(Confirmation);