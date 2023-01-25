import React from 'react';

import email from '/src/assets/images/email.png';
import otter from '/src/assets/images/footer-otter.png';
import github from '/src/assets/images/github.png';
import linkedin from '/src/assets/images/linkedin.png';
import twitter from '/src/assets/images/twitter-icon.png';

function Footer() {
  return (
    <footer className="bg-footer bg-no-repeat bg-cover">
      <div className="pt-64 2xl:pt-80 font-extralight">
        <div className="max-w-7xl mx-auto">
          <div className="grid grid-cols-5 pb-20">
            <img src={otter} alt="Sea otter" />
            <div className="font-bold leading-loose">
              <h3>PRODUCT</h3>
              <ul>
                <li>
                  <a href="https://rotational.io/ensign">Ensign</a>
                </li>
                <li>
                  <a href="#">Documentation</a>
                </li>
                <li>
                  <a href="#">Status</a>
                </li>
              </ul>
            </div>
            <div className="font-bold leading-loose">
              <h3>COMPANY</h3>
              <ul>
                <li>
                  <a href="https://rotational.io/services">Services</a>
                </li>
                <li>
                  <a href="https://rotational.io/blog">Blog</a>
                </li>
                <li>
                  <a href="https://rotational.io/opensource">Open Source</a>
                </li>
                <li>
                  <a href="https://rotational.io/about">About</a>
                </li>
                <li>
                  <a href="https://rotational.io/contact">Contact Us</a>
                </li>
              </ul>
            </div>
            <div className="font-bold leading-loose">
              <h3>INFORMATION</h3>
              <ul>
                <li>
                  <a href="https://rotational.io/privacy">Privacy</a>
                </li>
                <li>
                  <a href="https://rotational.io/terms">Terms</a>
                </li>
              </ul>
            </div>
            <div className="font-bold leading-loose">
              <h3>SDKS</h3>
              <ul>
                <li>
                  <a href="#">Go</a>
                </li>
                <li>
                  <a href="#">Python</a>
                </li>
              </ul>
            </div>
          </div>
          <div className="border-t py-6">
            <div className="grid grid-cols-4 ml-10">
              <a
                href="https://twitter.com/rotationalio"
                className="flex items-center"
                target="_blank"
                rel="noreferrer"
              >
                <img src={twitter} alt="Twitter logo" className="p-4 bg-white rounded-lg mr-3" />
                <span className="text-lg">rotationalio</span>
              </a>
              <a
                href="https://github.com/rotationalio"
                className="flex items-center"
                target="_blank"
                rel="noreferrer"
              >
                <img src={github} alt="GitHub logo" className="p-4 bg-white rounded-lg mr-3" />
                <span className="text-lg">rotationalio</span>
              </a>
              <a
                href="https://www.linkedin.com/company/rotational"
                className="flex items-center"
                target="_blank"
                rel="noreferrer"
              >
                <img src={linkedin} alt="LinkedIn logo" className="p-4 bg-white rounded-lg mr-3" />
                <span className="text-lg">Rotational</span>
              </a>
              <a
                href="mailto:info@rotational.io"
                className="flex items-center"
                target="_blank"
                rel="noreferrer"
              >
                <img src={email} alt="Envelope" className="p-4 bg-white rounded-lg mr-3" />
                <span className="text-lg">info@rotational.io</span>
              </a>
            </div>
          </div>
          <div className="sm:flex justify-between border-t py-6 text-white mt-4">
            <p className="text-base lg:text-xl">
              Copyright © {new Date().getFullYear()} Rotational Labs, Inc, All Rights Reserved
            </p>
            <ul className="sm:mt-0 mt-4 flex">
              <li className="border-r pr-4 mr-4 text-base lg:text-xl">
                <a href="https://rotational.io/privacy">Privacy Policy</a>
              </li>
              <li className="text-base lg:text-xl">
                <a href="https://rotational.io/terms">Terms of Use</a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
