import { Trans } from '@lingui/macro';
import { AriaButton as Button, Modal } from '@rotational/beacon-core';
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
          <p className="my-3">
            <Trans>Your eventing chariot awaits</Trans>
          </p>
          <Link to={PATH_DASHBOARD.ROOT}>
            <Button color="secondary" size="large">
              <Trans>Take the reins</Trans>
            </Button>
          </Link>
          <Link to={ROUTES.COMPLETE}>
            <Button color="ghost" className="font-normal text-blue-500">
              <Trans>Close</Trans>
            </Button>
          </Link>
        </div>
      </Modal>
    </>
  );
}
