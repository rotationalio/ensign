import { Modal } from '@rotational/beacon-core';
import { FormikHelpers } from 'formik';
import { useMemo } from 'react';
import { toast } from 'react-hot-toast';

import { Close } from '@/components/icons/close';
import Button from '@/components/ui/Button/Button';

import { useUpdateMemberRole } from '../hooks/useUpdateMemberRole';
import { ChangeRoleFormDto } from '../types/changeRoleFormDto';
import { Member } from '../types/member';
import ChangeRoleForm from './ChangeRoleForm';

type ChangeRoleModalProps = {
  openChangeRoleModal: {
    opened: boolean;
    member?: Member;
  };
  setOpenChangeRoleModal: React.Dispatch<
    React.SetStateAction<{
      opened: boolean;
      member?: Member | undefined;
    }>
  >;
};

function ChangeRoleModal({ openChangeRoleModal, setOpenChangeRoleModal }: ChangeRoleModalProps) {
  const { member } = openChangeRoleModal;
  const { updateMemberRole } = useUpdateMemberRole();

  const handleSubmit = (values: ChangeRoleFormDto, helpers: FormikHelpers<ChangeRoleFormDto>) => {
    updateMemberRole(
      {
        memberID: openChangeRoleModal?.member?.id,
        role: values.role,
      },
      {
        onError(error, _variables, _context) {
          toast.error((error as any)?.response?.data?.error || 'Something went wrong');
          if ((error as any)?.response.status === 400) {
            setOpenChangeRoleModal({ ...openChangeRoleModal, opened: false });
          }
        },
        onSuccess() {
          toast.success('Success! You have updated your teammate’s role in your organization.');
          setOpenChangeRoleModal({ ...openChangeRoleModal, opened: false });
        },
        onSettled() {
          helpers.setSubmitting(false);
        },
      }
    );
  };

  const initialValues = useMemo(
    () => ({
      name: member?.name || '',
      current_role: member?.role || '',
      role: member?.role || '',
    }),
    [member?.name, member?.role]
  );

  return (
    <div className="relative">
      <Modal
        open={openChangeRoleModal.opened}
        title="Change Role"
        containerClassName="overflow-scroll  max-w-[80vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="keyCreated"
      >
        <>
          <ChangeRoleForm handleSubmit={handleSubmit} initialValues={initialValues} />
          <Button
            onClick={() => setOpenChangeRoleModal({ ...openChangeRoleModal, opened: false })}
            variant="ghost"
            className="absolute top-2 right-2 min-h-fit min-w-fit py-2"
          >
            <Close className="text-primary-900" />
          </Button>
        </>
      </Modal>
    </div>
  );
}

export default ChangeRoleModal;
