import { useState } from 'react';

import { BlueBars } from '@/components/icons/blueBars';

import { DeleteOrg } from '../DeleteOrg';

export default function OrganizationDetails() {
  const [showButton, setShowButton] = useState(false);

  const handleToggleBars = () => {
    const open = showButton;
    setShowButton(!open);
  };

  const handleClose = () => setShowButton(false);

  return (
    <>
      <h3 className="mt-10 text-2xl font-bold">Organization Dashboard</h3>
      <h4 className="mt-10 max-w-4xl border-t border-primary-900 pt-4 text-xl font-bold">
        Organization Details
      </h4>
      <section className="mt-8 max-w-4xl rounded-md border-2 border-secondary-500 pl-6">
        <div className="mr-4 mt-2 flex justify-end">
          <BlueBars onClick={handleToggleBars} />
          <div className="relative left-12">{showButton && <DeleteOrg close={handleClose} />}</div>
        </div>
        <div className="flex gap-4 py-8">
          <h6 className="font-bold">Organization Name:</h6>
          <span>Name</span>
        </div>
        <div className="flex gap-36 pb-8">
          <h6 className="font-bold">URL:</h6>
          <span>Domain</span>
        </div>
        <div className="flex gap-32 pb-8">
          <h6 className="font-bold">Org ID:</h6>
          <span>ID</span>
        </div>
        <div className="flex gap-28 pb-8">
          <h6 className="font-bold">Owner:</h6>
          <span className="ml-3">Owner</span>
        </div>
        <div className="flex gap-28 pb-8">
          <h6 className="font-bold">Created:</h6>
          <span className="ml-1">Created</span>
        </div>
      </section>
    </>
  );
}
