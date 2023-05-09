import { Button, Modal } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { PATH_DASHBOARD, ROUTES } from '@/application';
import HeavyCheckMark from '@/components/icons/heavy-check-mark';

export type SuccessfulTenantCreationModalProps = {
  open: boolean;
};

export default function SuccessfulTenantCreationModal({
  open,
}: SuccessfulTenantCreationModalProps) {
  return (
    <>
      <Modal open={open} title="Success!" size="medium">
        <div className="grid place-items-center gap-3">
          <HeavyCheckMark className="h-20 w-20" />
          <p className="my-3">Your eventing chariot awaits</p>
          <Link to={PATH_DASHBOARD.ROOT}>
            <Button variant="secondary" size="large">
              Take the reins
            </Button>
          </Link>
          <Link to={ROUTES.COMPLETE}>
            <Button variant="ghost" className="font-normal text-blue-500">
              Close
            </Button>
          </Link>
        </div>
      </Modal>
    </>
  );
}
