import { Trans } from '@lingui/macro';
import { memo } from 'react';

import email from '@/assets/icons/email.png';
import github from '@/assets/icons/github.png';
import linkedin from '@/assets/icons/linkedin.png';
import twitter from '@/assets/icons/twitter.png';
import otter from '@/assets/images/footer-otter.png';

function LandingFooter() {
  return (
    <footer className="bg-footer bg-cover bg-no-repeat text-white">
      <div className="pt-64 2xl:pt-80">
        <div className="mx-auto max-w-7xl">
          <div className="grid grid-cols-4 pb-20">
            <img src={otter} alt="Sea otter" />
            <div className="font-bold leading-loose">
              <h3 className="font-['Quattrocento'] font-light">
                <Trans>PRODUCT</Trans>
              </h3>
              <ul>
                <li>
                  <a href="https://rotational.io/ensign">
                    <Trans>Ensign</Trans>
                  </a>
                </li>
                <li>
                  <a href="https://ensign.rotational.dev/getting-started/">
                    <Trans>Documentation</Trans>
                  </a>
                </li>
                {/* <li>
                  <a href="#">Status</a>
                </li> */}
              </ul>
            </div>
            <div className="font-bold leading-loose">
              <h3 className="font-['Quattrocento'] font-light">
                <Trans>COMPANY</Trans>
              </h3>
              <ul>
                <li>
                  <a href="https://rotational.io/services">
                    <Trans>Services</Trans>
                  </a>
                </li>
                <li>
                  <a href="https://rotational.io/blog">
                    <Trans>Blog</Trans>
                  </a>
                </li>
                <li>
                  <a href="https://rotational.io/opensource">
                    <Trans>Open Source</Trans>
                  </a>
                </li>
                <li>
                  <a href="https://rotational.io/about">
                    <Trans>About</Trans>
                  </a>
                </li>
                <li>
                  <a href="https://rotational.io/contact">
                    <Trans>Contact Us</Trans>
                  </a>
                </li>
              </ul>
            </div>
            <div className="font-bold leading-loose">
              <h3 className="font-['Quattrocento'] font-light">
                <Trans>INFORMATION</Trans>
              </h3>
              <ul>
                <li>
                  <a href="https://rotational.io/privacy">
                    <Trans>Privacy</Trans>
                  </a>
                </li>
                <li>
                  <a href="https://rotational.io/terms">
                    <Trans>Terms</Trans>
                  </a>
                </li>
              </ul>
            </div>
            {/* <div className="font-bold leading-loose">
              <h3 className="font-['Quattrocento']">SDKS</h3>
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
              <div className="hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="https://twitter.com/rotationalio"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={twitter} alt="Twitter logo" className="" />
                  <span className="ml-4">
                    <Trans>Twitter</Trans>
                  </span>
                </a>
              </div>
              <div className="hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="https://github.com/rotationalio"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={github} alt="GitHub logo" className="" />
                  <span className="ml-4">
                    <Trans>GitHub</Trans>
                  </span>
                </a>
              </div>
              <div className="hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="https://www.linkedin.com/company/rotational"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={linkedin} alt="LinkedIn logo" className="" />
                  <span className="mt-2 ml-4">
                    <Trans>LinkedIn</Trans>
                  </span>
                </a>
              </div>
              <div className="mt-2 hover:max-w-[50%] hover:rounded-full hover:bg-icon-hover">
                <a
                  href="mailto:info@rotational.io"
                  className="flex items-center"
                  target="_blank"
                  rel="noreferrer"
                >
                  <img src={email} alt="Envelope" className="" />
                  <span className="ml-4">
                    <Trans>Email</Trans>
                  </span>
                </a>
              </div>
            </div>
          </div>
          <div className="mt-4 justify-between py-10 text-white sm:flex">
            <p className="">
              <Trans>
                Copyright © {new Date().getFullYear()} Rotational Labs, Inc · All Rights Reserved
              </Trans>
            </p>
            <ul className="mt-4 flex sm:mt-0">
              <li className="mr-4 border-r pr-4">
                <a href="https://rotational.io/privacy">
                  <Trans>Privacy Policy</Trans>
                </a>
              </li>
              <li className="">
                <a href="https://rotational.io/terms">
                  <Trans>Terms of Use</Trans>
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default memo(LandingFooter);
