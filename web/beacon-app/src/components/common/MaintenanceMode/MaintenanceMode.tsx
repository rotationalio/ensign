import { memo } from 'react';
import busyOtters from '/src/assets/images/busy-sea-otters.png';

function MaintenanceMode() {
  return (
    <section className="mx-auto max-w-4xl rounded-lg border border-solid border-primary-800 text-2xl">
      <p className="mx-auto mt-8 max-w-xl">
        Ensign is temporarily undergoing scheduled maintenance. We&#39;ll be back online shortly. See
        our{' '}
        <span className="font-bold text-primary">
          <a href="#">operating status</a>
        </span>{' '}
        for additional information. Contact us any time with questions.
      </p>
      <img src={busyOtters} alt="" className="mx-auto mt-10" />
      <p className="mt-8 pb-20 text-center">
        Enjoy a cup of coffee and catch up on our latest{' '}
        <span className="font-bold text-primary">
          <a href="https://rotational.io/blog">blog posts!</a>
        </span>
      </p>
    </section>
  );
}

export default memo(MaintenanceMode)