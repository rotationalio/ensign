import { memo } from 'react';

import { EXTRENAL_LINKS, ROUTES } from '@/application';
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
                  <a href={ROUTES.HOME}>Ensign</a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.DOCUMENTATION} target="_blank" rel="noreferrer">
                    Documentation
                  </a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.SDKs} target="_blank" rel="noreferrer">
                    SDKs
                  </a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.SERVER} target="_blank" rel="noreferrer">
                    Status
                  </a>
                </li>
              </ul>
            </div>
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">COMPANY</h3>
              <ul>
                <li>
                  <a href={EXTRENAL_LINKS.SERVICES}>Services</a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.BLOG}>Blog</a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.OPEN_SOURCE}>Open Source</a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.ABOUT}>About</a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.CONTACT}>Contact Us</a>
                </li>
              </ul>
            </div>
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">INFORMATION</h3>
              <ul>
                <li>
                  <a href={EXTRENAL_LINKS.PRIVACY}>Privacy</a>
                </li>
                <li>
                  <a href={EXTRENAL_LINKS.TERMS}>Terms</a>
                </li>
              </ul>
            </div>
          </div>
          <div className="mt-12 max-w-7xl justify-between border-t px-6 sm:mt-32 sm:flex">
            <div className="mx-auto mt-8 sm:mt-0 xl:ml-5">
              <div className="mx-auto grid grid-cols-2 gap-x-20 gap-y-6 sm:mt-4 md:gap-x-32 lg:grid-cols-4 xl:ml-20">
                <div>
                  <a
                    href={EXTRENAL_LINKS.TWITTER}
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
                    href={EXTRENAL_LINKS.GITHUB}
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
                    href={EXTRENAL_LINKS.LINKEDIN}
                    className="icon-hover"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <LinkedInIcon />
                    <span className="mt-1 ml-4">LinkedIn</span>
                  </a>
                </div>
                <div>
                  <a href={EXTRENAL_LINKS.EMAIL_US} className="icon-hover">
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
                <a href={EXTRENAL_LINKS.PRIVACY}>Privacy Policy</a>
              </li>
              <li className="">
                <a href={EXTRENAL_LINKS.TERMS}>Terms of Use</a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default memo(LandingFooter);
