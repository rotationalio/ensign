import React from 'react';
import footer from './img/foot.png'
import whiteLogo from './img/logo-white.png'
import email from './img/email.png'
import github from './img/github.png'
import linkedin from './img/linkedin.png'
import twitter from './img/twitter-icon.png'

export default function Footer() {
    return (
        <footer>
        <div className="relative pt-[90px] sm:pt-[300px] font-extralight">
            <img 
                className="w-full absolute inset-y-0 left-0 lg:h-full lg:w-full lg:max-w-none" 
                src={footer} 
                alt="Illustration of a blue lighthouse on the left hand side." 
                />
            <div className="relative max-w-7xl mx-auto px-4 sm:px-6 bg-[#192E5B]">
                <div className="lg:flex text-white">
                <div>
                    <a href="https://rotational.io/" target="_blank" rel="noreferrer">
                        <img 
                            src={whiteLogo} 
                            alt="Rotational Labs logo" 
                            className="sm:h-auto h-16" 
                            />
                    </a>
                </div>
                    <div className="sm:ml-24 mt-5 sm:mt-0 max-w-lg pb-8">
                         <h5 className="font-bold text-lg">Contact Details</h5>
                            <p className="mt-3 text-lg">
                                We love to talk through complex scaling and performance issues. What's slowing you down? Drop us a line for a friendly chat. We'll make it worth your while.
                            </p>
                    </div>
                    <div className="sm:ml-24 sm:mt-0 mt-8">
                        <ul className="sm:grid grid-cols-2 lg:flex flex-col gap-0 mr-5 xl:grid grid-cols-2 gap-x-20 gap-8">
                                <li className="pb-8">
                                    <a href="mailto:info@rotational.io" className="flex items-center" target="_blank" rel="noreferrer">
                                        <img 
                                            src={email} 
                                            alt="Envelope" 
                                            className="p-4 bg-white rounded-lg mr-4" 
                                            />
                                        <span className="text-lg">info@rotational.io</span>
                                    </a>
                                </li>
                                <li className="pb-8">
                                    <a href="https://github.com/rotationalio" className="flex items-center" target="_blank" rel="noreferrer">
                                        <img 
                                            src={github} 
                                            alt="GitHub logo" 
                                            className="p-4 bg-white rounded-lg mr-4" 
                                            />
                                        <span className="text-lg">rotationalio</span>
                                    </a>
                                </li>
                                <li className="pb-8">
                                    <a href="https://twitter.com/rotationalio" className="flex items-center" target="_blank" rel="noreferrer">
                                        <img 
                                            src={twitter} 
                                            alt="Twitter logo" 
                                            className="p-4 bg-white rounded-lg mr-4" 
                                            />
                                        <span className="text-lg">rotationalio</span>
                                    </a>
                                </li>
                                <li>
                                    <a href="https://www.linkedin.com/company/rotational" className="flex items-center" target="_blank" rel="noreferrer">
                                        <img 
                                            src={linkedin} 
                                            alt="LinkedIn logo" 
                                            className="p-4 bg-white rounded-lg mr-4" 
                                            />
                                        <span className="text-lg">Rotational</span>
                                    </a>
                                </li>
                        </ul>
                    </div>
                </div>
                <div className="sm:flex justify-between border-t py-6 text-white mt-12 sm:mt-32">
                    <p className="text-base lg:text-xl">Copyright © {new Date().getFullYear()} Rotational Labs, LLC, All Rights Reserved</p>
                    <ul className="sm:mt-0 mt-4 flex">
                        <li className="border-r pr-4 mr-4 text-base lg:text-xl"><a href="https://rotational.io/">Privacy Policy</a></li>
                        <li className="text-base lg:text-xl"><a href="https://rotational.io/">Terms of Use</a></li>
                    </ul>
                </div>
            </div>
        </div>
    </footer>
    )
}