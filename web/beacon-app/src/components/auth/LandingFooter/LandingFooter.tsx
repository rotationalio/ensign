import { memo } from 'react';

import { EXTERNAL_LINKS, ROUTES } from '@/application';
import { appConfig } from '@/application/config';
import EmailIcon from '@/assets/icons/emailIcon';
import GitHubIcon from '@/assets/icons/githubIcon';
import LinkedInIcon from '@/assets/icons/linkedinIcon';
import TwitterIcon from '@/assets/icons/twitterIcon';
function LandingFooter() {
  const { version: appVersion, revision: gitRevision } = appConfig;
  return (
    <footer className="bg-footer bg-cover bg-no-repeat text-white ">
      <div className="pt-72 2xl:pt-80">
        <div className="mx-auto max-w-7xl">
          <div className="mx-auto grid-cols-3 text-center sm:ml-0 sm:grid sm:text-left">
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">PRODUCT</h3>
              <ul>
                <li>
                  <a href={ROUTES.HOME}>Ensign</a>
                </li>
                <li>
                  <a href={EXTERNAL_LINKS.DOCUMENTATION} target="_blank" rel="noreferrer">
                    Documentation
                  </a>
                </li>
                <li>
                  <a href={EXTERNAL_LINKS.SDKs} target="_blank" rel="noreferrer">
                    SDKs
                  </a>
                </li>
                <li>
                  <a href={EXTERNAL_LINKS.SERVER} target="_blank" rel="noreferrer">
                    Status
                  </a>
                </li>
              </ul>
            </div>
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">COMPANY</h3>
              <ul>
                <li>
                  <a href={EXTERNAL_LINKS.SERVICES}>Services</a>
                </li>
                <li>
                  <a href={EXTERNAL_LINKS.BLOG}>Blog</a>
                </li>
                <li>
                  <a href={EXTERNAL_LINKS.ABOUT}>About</a>
                </li>
              </ul>
            </div>
            <div className="pt-4 font-bold leading-loose">
              <h3 className="font-light">COMMUNITY</h3>
              <ul>
                <li>
                  <a href={EXTERNAL_LINKS.DATA_PLAYGROUND}>Data Playground</a>
                </li>
                <li>
                  <a href={EXTERNAL_LINKS.OPEN_SOURCE}>Open Source</a>
                </li>
              </ul>
            </div>
          </div>
          <div className="mt-12 max-w-7xl justify-between border-t px-6 sm:mt-32 sm:flex">
            <div className="mx-auto mt-8 sm:mt-0 xl:ml-5">
              <div className="mx-auto grid grid-cols-2 gap-x-20 gap-y-6 sm:mt-4 md:gap-x-32 lg:grid-cols-4">
                <div>
                  <a
                    href={EXTERNAL_LINKS.TWITTER}
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
                    href={EXTERNAL_LINKS.GITHUB}
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
                    href={EXTERNAL_LINKS.LINKEDIN}
                    className="icon-hover"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <LinkedInIcon />
                    <span className="mt-1 ml-4">LinkedIn</span>
                  </a>
                </div>
                <div>
                  <a href={EXTERNAL_LINKS.EMAIL_US} className="icon-hover">
                    <EmailIcon />
                    <span className="mt-1 ml-4">Email</span>
                  </a>
                </div>
              </div>
            </div>
          </div>
          <div className="mt-4 justify-between px-6 pt-10  text-white sm:flex">
            <p className="">
              Copyright © Rotational Labs, Inc. {new Date().getFullYear()} · All Rights Reserved
            </p>

            <ul className="mt-4 flex sm:mt-0">
              <li className="mr-4 border-r pr-4">
                <a href={EXTERNAL_LINKS.PRIVACY}>Privacy Policy</a>
              </li>
              <li className="">
                <a href={EXTERNAL_LINKS.TERMS}>Terms of Use</a>
              </li>
            </ul>
          </div>
          <div className="justify-between py-3 px-6 text-center">
            <p>
              {appVersion && <span className="text-xs text-white">App Version {appVersion} </span>}
              {gitRevision && (
                <span className="text-xs text-white">& Git Revision {gitRevision} </span>
              )}
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default memo(LandingFooter);
