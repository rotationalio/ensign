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
    <Modal
      open={open}
      title={t`Change Password`}
      containerClassName="w-[25vw]"
      onClose={handleModalClose}
    >
      <>
        <div className="my-5">
          <p className="text-base">
            <Trans>
              <span className="font-">Note: </span>You will be logged out when you change your
              password. You are required to log in again with your new password.
            </Trans>
          </p>
        </div>
        <ChangePasswordForm handleSubmit={handleSubmit} initialValues={values} />
      </>
    </Modal>
  );
}

export default ChangePasswordModal;
