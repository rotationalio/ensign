import { t } from '@lingui/macro';
import { Modal } from '@rotational/beacon-core';

import ChangePasswordForm from './ChangePasswordForm';

type ChangeRoleModalProps = {
  open: boolean;
  values?: any;
  handleModalClose: () => void;
};

function ChangePasswordModal({ open, handleModalClose, values }: ChangeRoleModalProps) {
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
        <ChangePasswordForm handleSubmit={handleSubmit} initialValues={values} />
      </>
    </Modal>
  );
}

export default ChangePasswordModal;
