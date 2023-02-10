import { Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { BlueBars } from '@/components/icons/blueBars';

import { useFetchOrg } from '../../hooks/useFetchOrgDetail';
import { DeleteOrg } from '../DeleteOrg';

export default function OrganizationDetails() {
  const [isOpen, setIsOpen] = useState(false);

  const handleToggleBars = () => {
    const open = isOpen;
    setIsOpen(!open);
  };

  const handleClose = () => setIsOpen(false);

  const { org, hasOrgFailed, isFetchingOrg, error } = useFetchOrg('orgID');

  const { id, name, domain, created } = org;

  if (isFetchingOrg) {
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasOrgFailed}
        onClose={handleClose}
        variant="danger"
        title="We are unable to fetch your organization, please try again."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  return (
    <>
      <h3 className="mt-10 text-2xl font-bold">Organization Dashboard</h3>
      <h4 className="mt-10 max-w-4xl border-t border-primary-900 pt-4 text-xl font-bold">
        Organization Details
      </h4>
      <section className="mt-8 max-w-4xl rounded-md border-2 border-secondary-500 pl-6">
        <div className="mr-4 mt-2 flex justify-end">
          <BlueBars onClick={handleToggleBars} />
          <div className="relative left-12">{isOpen && <DeleteOrg close={handleClose} />}</div>
        </div>
        <div className="flex gap-4 py-8">
          <h6 className="font-bold">Organization Name:</h6>
          <span>{name}</span>
        </div>
        <div className="flex gap-36 pb-8">
          <h6 className="font-bold">URL:</h6>
          <span>{domain}</span>
        </div>
        <div className="flex gap-32 pb-8">
          <h6 className="font-bold">Org ID:</h6>
          <span>{id}</span>
        </div>
        <div className="flex gap-28 pb-8">
          <h6 className="font-bold">Owner:</h6>
          <span className="ml-3">Owner</span>
        </div>
        <div className="flex gap-28 pb-8">
          <h6 className="font-bold">Created:</h6>
          <span className="ml-1">{created}</span>
        </div>
      </section>
    </>
  );
}
