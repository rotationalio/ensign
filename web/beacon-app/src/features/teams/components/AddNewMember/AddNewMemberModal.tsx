import { Modal } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { toast } from 'react-hot-toast';
import { useCreateMember } from '@/features/members/hooks/useCreateMember';

import AddNewMemberForm from './AddNewMemberForm';

type AddNewMemberModalProps = {
  isOpened: boolean;
  onClose: () => void;
};

function AddNewMemberModal({ onClose, isOpened }: AddNewMemberModalProps) {
  const { createMember, isCreatingMember, wasMemberCreated, hasMemberFailed, error, reset } =
    useCreateMember();

  const handleSubmit = async (values: any) => {
    await createMember(values);
  };

  useEffect(() => {
    if (wasMemberCreated) {
      toast.success('Success! You have invited your teammate to join your organization.');
      onClose();
      reset();
    }
  }, [wasMemberCreated, onClose, reset]);

  useEffect(() => {
    if (hasMemberFailed) {
      toast.error(
        (error as any)?.response?.data?.error ||
          `Could not create member. Please try again or contact support, if the problem continues.`
      );
      reset();
    }
  }, [hasMemberFailed, error, reset]);

  return (
    <div className="relative">
      <Modal
        open={isOpened}
        title="Invite New Team Member"
        containerClassName="overflow-scroll max-h-[100vh] max-w-[100vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="memberCreationModal"
        onClose={onClose}
      >
        <>
          <AddNewMemberForm onSubmit={handleSubmit} isSubmitting={isCreatingMember} />
        </>
      </Modal>
    </div>
  );
}

export default AddNewMemberModal;
