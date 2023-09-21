import { t, Trans } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';

import ChangePasswordForm from './ChangePasswordForm';

type ChangePasswordFormProps = {
  open: boolean;
  values?: any;
  handleModalClose: () => void;
};

function ChangePasswordModal({ open, handleModalClose, values }: ChangePasswordFormProps) {
  const handleSubmit = (values: any) => {
    console.log('values', values);
  };

  return (
    <Modal open={open} title={t`Change Password`} onClose={handleModalClose}>
      <>
        <div className="my-5">
          <p className="text-base">
            <Trans>
              <span className="font-bold">Note: </span>You will be logged out when you change your
              password. You will be required to log in again with your new password.
            </Trans>
          </p>
        </div>
        <ChangePasswordForm handleSubmit={handleSubmit} initialValues={values} />
      </>
    </Modal>
  );
}

export default ChangePasswordModal;
