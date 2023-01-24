import React from 'react';
import email from 'src/assets/images/email.png';
import github from 'src/assets/images/github.png';
import linkedin from 'src/assets/images/linkedin.png';
import twitter from 'src/assets/images/twitter-icon.png';

const Footer = () => {
  return (
    <footer>
      <div className="bg-footer relative pt-[90px] pt-[350px] font-extralight">
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6">
          <div className="lg:grid grid-cols-3 text-white"></div>
          <div className="ml-10 xl:ml-5 sm:mt-0 mt-8">
            <ul className="sm:grid grid-cols-2 lg:flex flex-col gap-0 mr-5 xl:grid grid-cols-2 gap-x-28 gap-y-4">
              <li className="pb-8">
                <a
                  href="mailto:info@rotational.io"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={email} alt="Envelope" className="p-4 bg-white rounded-lg mr-3" />
                  <span className="text-lg">info@rotational.io</span>
                </a>
              </li>
              <li className="pb-8">
                <a
                  href="https://github.com/rotationalio"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={github} alt="GitHub logo" className="p-4 bg-white rounded-lg mr-3" />
                  <span className="text-lg">rotationalio</span>
                </a>
              </li>
              <li className="pb-8">
                <a
                  href="https://twitter.com/rotationalio"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={twitter} alt="Twitter logo" className="p-4 bg-white rounded-lg mr-3" />
                  <span className="text-lg">rotationalio</span>
                </a>
              </li>
              <li>
                <a
                  href="https://www.linkedin.com/company/rotational"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img
                    src={linkedin}
                    alt="LinkedIn logo"
                    className="p-4 bg-white rounded-lg mr-3"
                  />
                  <span className="text-lg">Rotational</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
        <div className="sm:flex justify-between border-t py-6 text-white mt-12 sm:mt-32">
          <p className="text-base lg:text-xl">
            Copyright © {new Date().getFullYear()} Rotational Labs, Inc, All Rights Reserved
          </p>
          <ul className="sm:mt-0 mt-4 flex">
            <li className="border-r pr-4 mr-4 text-base lg:text-xl">
              <a href="https://rotational.io/">Privacy Policy</a>
            </li>
            <li className="text-base lg:text-xl">
              <a href="https://rotational.io/">Terms of Use</a>
            </li>
          </ul>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
