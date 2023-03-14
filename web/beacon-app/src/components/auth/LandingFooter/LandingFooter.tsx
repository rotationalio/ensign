import { memo } from 'react';

import EmailIcon from '@/assets/icons/emailIcon';
import GitHubIcon from '@/assets/icons/githubIcon';
import LinkedInIcon from '@/assets/icons/linkedinIcon';
import TwitterIcon from '@/assets/icons/twitterIcon';
import SeaOtter from '@/assets/images/seaOtter';

function LandingFooter() {
  return (
    <footer className="bg-footer bg-cover bg-no-repeat text-white ">
      <div className="pt-64 2xl:pt-80">
        <div className="mx-auto max-w-7xl">
          <div className="grid grid-cols-4 pb-20">
            <SeaOtter />
            <div className="font-bold leading-loose">
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
            <div className="font-bold leading-loose">
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
            <div className="font-bold leading-loose">
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
          <div className="border-t py-6">
            <div className="ml-20 grid grid-cols-4">
              <div className="pl-2 pt-2 hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="https://twitter.com/rotationalio"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <TwitterIcon />
                  <span className="ml-4">Twitter</span>
                </a>
              </div>
              <div className="pt-2 pl-2 hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="https://github.com/rotationalio"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <GitHubIcon />
                  <span className="ml-4">GitHub</span>
                </a>
              </div>
              <div className="pt-1 pl-2 hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="https://www.linkedin.com/company/rotational"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <LinkedInIcon />
                  <span className="mt-1 ml-4">LinkedIn</span>
                </a>
              </div>
              <div className="py-2 pl-2 hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="mailto:info@rotational.io"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <EmailIcon />
                  <span className="-mt-2 ml-4">Email</span>
                </a>
              </div>
            </div>
          </div>
          <div className="mt-4 justify-between py-10 text-white sm:flex">
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
