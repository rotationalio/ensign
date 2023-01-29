import { Link } from 'react-router-dom';
import checkmark from '/src/assets/images/checkmark.png';

import data from '/src/assets/images/hosted-data-icon.png';


export default function AccessDashboard() {
  return (
    <section className="grid max-w-6xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <img src={data} alt="" className="mx-auto mt-6" />
      <div>
        <h2 className="mt-8 font-bold">
          Set up Your Tenant <span className="font-normal">(required)</span>
        </h2>
        <p className="mt-8">
          A tenant is a collection of settings. The tenant is your locus of control when setting up
          projects and topics.
        </p>
      </div>
      <div>
      {/*     Make green cirlce the background image
       */}{' '}
      <img src={checkmark} alt="" />
      <Link to="/">
        Access Dashboard
        </Link>
      </div>
    </section>
  );
}
