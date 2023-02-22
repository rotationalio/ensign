import { t, Trans } from '@lingui/macro';
import { Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { BlueBars } from '@/components/icons/blueBars';

import { useFetchMembers } from '../hooks/useFetchMembers';
import CancelAccount from './CancelAccount';

export default function MemberDetails() {
  const [isOpen, setIsOpen] = useState(false);

  const handleToggleBars = () => {
    const open = isOpen;
    setIsOpen(!open);
  };

  const handleClose = () => setIsOpen(false);

  const { member, hasMemberFailed, isFetchingMember, error } = useFetchMembers();

  const { id, name, created } = member;

  if (isFetchingMember) {
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasMemberFailed}
        onClose={handleClose}
        variant="danger"
        title={t`We are unable to fetch your member, please try again.`}
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  return (
    <>
      <h3 className="mt-10 text-2xl font-bold">
        <Trans>User Profile</Trans>
      </h3>
      <h4 className="mt-10 border-t border-primary-900 pt-4 text-xl font-bold">
        <Trans>User Details</Trans>
      </h4>
      <section className="mt-8 rounded-md border-2 border-secondary-500 pl-6">
        <div className="mr-4 mt-3 flex justify-end">
          <BlueBars onClick={handleToggleBars} />
          <div>{isOpen && <CancelAccount close={handleClose} />} </div>
        </div>
        <div className="flex gap-16 pt-4 pb-8">
          <h6 className="font-bold">
            <Trans>User Name:</Trans>
          </h6>
          <span className="ml-3">{name}</span>
        </div>
        {/*         <div className="flex gap-12 pb-8">
          <h6 className="font-bold">Email Address:</h6>
          <span>test@example.com</span>
        </div> */}
        <div className="flex gap-24 pb-8">
          <h6 className="font-bold">
            <Trans>User ID:</Trans>
          </h6>
          <span className="ml-3">{id}</span>
        </div>
        <div className="flex gap-24 pb-8">
          <h6 className="font-bold">
            <Trans>Created:</Trans>
          </h6>
          <span className="ml-2">{created}</span>
        </div>
      </section>
    </>
  );
}
