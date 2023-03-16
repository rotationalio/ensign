import { memo } from 'react';

import EmailIcon from '@/assets/icons/emailIcon';
import GitHubIcon from '@/assets/icons/githubIcon';
import LinkedInIcon from '@/assets/icons/linkedinIcon';
import TwitterIcon from '@/assets/icons/twitterIcon';
import SeaOtter from '@/assets/images/seaOtter';

function LandingFooter() {
  return (
    <footer className="bg-footer bg-cover bg-no-repeat text-white ">
      <div className="pt-72 2xl:pt-80">
        <div className="mx-auto max-w-7xl">
          <div className="mx-auto grid-cols-4 text-center sm:ml-0 sm:grid sm:text-left">
            <SeaOtter />
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">PRODUCT</h3>
              <ul>
                <li>
                  <a href="https://rotational.io/ensign">Ensign</a>
                </li>
                <li>
                  <a href="https://ensign.rotational.dev/getting-started/">Documentation</a>
                </li>
                {/* <li>
                  <a href="#">Status</a>
                </li> */}
              </ul>
            </div>
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">COMPANY</h3>
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
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">INFORMATION</h3>
              <ul>
                <li>
                  <a href="https://rotational.io/privacy">Privacy</a>
                </li>
                <li>
                  <a href="https://rotational.io/terms">Terms</a>
                </li>
              </ul>
            </div>
            {/* <div className="font-bold leading-loose">
              <h3>SDKS</h3>
              <ul>
                <li>
                  <a href="#">Go</a>
                </li>
                <li>
                  <a href="#">Python</a>
                </li>
              </ul>
            </div> */}
          </div>
          <div className="mt-12 max-w-7xl justify-between border-t px-6 sm:mt-32 sm:flex">
            <div className="mx-auto mt-8 sm:mt-0 xl:ml-5">
              <div className="mx-auto grid grid-cols-2 gap-x-20 gap-y-6 sm:mt-4 md:gap-x-32 lg:grid-cols-4 xl:ml-20">
                <div>
                  <a
                    href="https://twitter.com/rotationalio"
                    className="icon-hover"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <TwitterIcon />
                    <span className="ml-4">Twitter</span>
                  </a>
                </div>
                <div>
                  <a
                    href="https://github.com/rotationalio"
                    className="icon-hover"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <GitHubIcon />
                    <span className="ml-4">GitHub</span>
                  </a>
                </div>
                <div>
                  <a
                    href="https://www.linkedin.com/company/rotational"
                    className="icon-hover"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <LinkedInIcon />
                    <span className="mt-1 ml-4">LinkedIn</span>
                  </a>
                </div>
                <div>
                  <a
                    href="mailto:info@rotational.io"
                    className="icon-hover"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <EmailIcon />
                    <span className="mt-1 ml-4">Email</span>
                  </a>
                </div>
              </div>
            </div>
          </div>
          <div className="mt-4 justify-between px-6 py-10 text-white sm:flex">
            <p className="">
              Copyright © {new Date().getFullYear()} Rotational Labs, Inc · All Rights Reserved
            </p>
            <ul className="mt-4 flex sm:mt-0">
              <li className="mr-4 border-r pr-4">
                <a href="https://rotational.io/privacy">Privacy Policy</a>
              </li>
              <li className="">
                <a href="https://rotational.io/terms">Terms of Use</a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default memo(LandingFooter);
