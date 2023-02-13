import { Modal } from '@rotational/beacon-core';
import { Fragment } from 'react';

import { Close } from '@/components/icons/close';

export default function CancelAcctModal(props: any) {
  return (
    <div>
      <Modal title="Cancel Account" open={true} containerClassName="max-w-md">
        <Fragment key=".0">
          <Close onClick={props.close} className="absolute top-4 right-8"></Close>
          <p className="pb-4">
            Please contact us at <span className="font-bold">support@rotational.io</span> to cancel
            your account. Please include your name, email, and Org ID in your request to cancel your
            account. We are working on an automated process to cancel accounts and appreciate your
            patience.
          </p>
          <p className="pb-4">
            You are the <span className="font-bold">Owner</span> of this account. If you cancel your
            account, your Organization, Tenant, Project, and Topic and all associated data will be{' '}
            <span className="font-bold">permanently</span> deleted.
          </p>
        </Fragment>
      </Modal>
    </div>
  );
}
