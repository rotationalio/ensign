import { useState } from 'react';

import { BlueBars } from '@/components/icons/blueBars';

import CancelAccount from './CancelAccount';

export default function MemberDetails() {
  const [showButton, setShowButton] = useState(false);

  const handleOpen = () => setShowButton(true);
  const handleClose = () => setShowButton(false);

  return (
    <>
      <h3 className="mt-10 text-2xl font-bold">User Profile</h3>
      <h4 className="mt-10 max-w-4xl border-t border-primary-900 pt-4 text-xl font-bold">
        User Details
      </h4>
      <section className="mt-8 max-w-4xl rounded-md border-2 border-secondary-500 pl-6">
        <div className="mr-4 mt-3 flex justify-end">
          <BlueBars onClick={handleOpen} />
          <div>{showButton && <CancelAccount close={handleClose} />} </div>
        </div>
        <div className="flex gap-16 pt-4 pb-8">
          <h6 className="font-bold">User Name:</h6>
          <span className="ml-3">Ryan Wilder</span>
        </div>
        <div className="flex gap-12 pb-8">
          <h6 className="font-bold">Email Address:</h6>
          <span>test@example.com</span>
        </div>
        <div className="flex gap-24 pb-8">
          <h6 className="font-bold">User ID:</h6>
          <span className="ml-3">ID</span>
        </div>
        <div className="flex gap-24 pb-8">
          <h6 className="font-bold">Created:</h6>
          <span className="ml-2">02.10.2023</span>
        </div>
      </section>
    </>
  );
}
