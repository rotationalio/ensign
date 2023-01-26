import email from '/src/assets/icons/email.png';
import github from '/src/assets/icons/github.png';
import linkedin from '/src/assets/icons/linkedin.png';
import twitter from '/src/assets/icons/twitter.png';
import otter from '/src/assets/images/footer-otter.png';

function Footer() {
  return (
    <footer className="bg-footer bg-cover bg-no-repeat">
      <div className="pt-64 2xl:pt-80">
        <div className="mx-auto max-w-7xl">
          <div className="grid grid-cols-4 pb-20">
            <img src={otter} alt="Sea otter" />
            <div className="font-bold leading-loose">
              <h3 className="font-['Quattrocento']">PRODUCT</h3>
              <ul>
                <li>
                  <a href="https://rotational.io/ensign">Ensign</a>
                </li>
                {/* <li>
                  <a href="#">Documentation</a>
                </li>
                <li>
                  <a href="#">Status</a>
                </li> */}
              </ul>
            </div>
            <div className="font-bold leading-loose">
              <h3 className="font-['Quattrocento']">COMPANY</h3>
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
              <h3 className="font-['Quattrocento']">INFORMATION</h3>
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
                  <span className="ml-4">Twitter</span>
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
                  <span className="ml-4">GitHub</span>
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
                  <span className="mt-2 ml-4">LinkedIn</span>
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
                  <span className="ml-4">Email</span>
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

export default Footer;
