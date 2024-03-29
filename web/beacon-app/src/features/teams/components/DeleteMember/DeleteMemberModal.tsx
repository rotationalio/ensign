import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';

import { useDeleteMember } from '@/features/members/hooks/useDeleteMember';

import DeleteMemberForm from './DeleteMemberForm';

type DeleteMemberModalProps = {
  onOpen: {
    opened: boolean;
    member?: any;
  };
  onClose: () => void;
};

function DeleteMemberModal({ onClose, onOpen }: DeleteMemberModalProps) {
  const { opened, member } = onOpen;
  const { deleteMember, isDeletingMember, wasMemberDeleted, hasMemberFailed, error, reset } =
    useDeleteMember(member?.id);
  const initialValues = {
    id: member?.id,
    name: member?.name || member?.email || '-',
    delete_agreement: false,
  };

  const onDelete = () => {
    deleteMember();
  };

  useEffect(() => {
    if (wasMemberDeleted) {
      toast.success('Success! You have removed your teammate from your organization');
      onClose();
      reset();
    }
  }, [wasMemberDeleted, onClose, reset]);

  useEffect(() => {
    if (hasMemberFailed) {
      toast.error(
        (error as any)?.response?.data?.error ||
          `We could not complete the action. Please try again or contact support, if the problem continues.`
      );
      reset();
    }
  }, [hasMemberFailed, error, reset]);

  return (
    <div className="relative">
      <Modal
        open={opened}
        title="Remove Team Member"
        containerClassName="overflow-scroll w-[50vh] max-h-[100vh] max-w-[100vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="delete-member-modal"
        onClose={onClose}
      >
        <>
          <DeleteMemberForm
            onSubmit={onDelete}
            isSubmitting={isDeletingMember}
            initialValues={initialValues}
          />
        </>
      </Modal>
    </div>
  );
}

export default DeleteMemberModal;
