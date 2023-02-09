import { Modal } from '@rotational/beacon-core';
import { Fragment } from 'react';

import { Close } from '@/components/icons/close';

export default function DeleteOrgModal(props: any) {
  return (
    <div>
      <Modal title="Delete Organization" open={true} containerClassName="max-w-md">
        <Fragment key=".0">
          <Close onClick={props.close} className="absolute top-4 right-8"></Close>
          <p className="pb-4">
            Please contact us at <span className="font-bold">support@rotational.io</span> to delete
            your organization. Please include your name, email, and Org ID in your request to delete
            the organization. We are working on an automated process to delete organizations and
            appreciate your patience.
          </p>
          <p className="pb-4">
            Please note that deleting the organization will{' '}
            <span className="font-bold">permanently</span> delete the Tenant, Project, Users, Topics
            and all data associated with the organization.
          </p>
        </Fragment>
      </Modal>
    </div>
  );
}
