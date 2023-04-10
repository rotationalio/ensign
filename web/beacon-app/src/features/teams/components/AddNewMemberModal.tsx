import { Modal } from '@rotational/beacon-core';
import { toast } from 'react-hot-toast';

import { Close } from '@/components/icons/close';
import Button from '@/components/ui/Button/Button';
import { useCreateMember } from '@/features/members/hooks/useCreateMember';

import AddNewMemberForm from './AddNewMemberForm';

type AddNewMemberModalProps = {
  isOpened: boolean;
  onClose: () => void;
};

function AddNewMemberModal({ onClose, isOpened }: AddNewMemberModalProps) {
  const { createMember, isCreatingMember, wasMemberCreated, hasMemberFailed, error } =
    useCreateMember();

  const handleSubmit = async (values: any) => {
    await createMember(values);
  };

  if (wasMemberCreated) {
    toast.success('Member created successfully');
    onClose();
  }
  if (hasMemberFailed) {
    toast.error(
      (error as any)?.response?.data?.error ||
        `Member creation failed, please try again or contact support if the problem persists.`
    );
  }

  return (
    <div className="relative">
      <Modal
        open={isOpened}
        title="Invite New Team Member"
        containerClassName="overflow-scroll h-[40vh] max-h-[100vh] max-w-[100vw] lg:max-w-[50vw] no-scrollbar"
        data-testid="memberCreationModal"
      >
        <>
          <AddNewMemberForm onSubmit={handleSubmit} isSubmitting={isCreatingMember} />
          <Button
            onClick={onClose}
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

export default AddNewMemberModal;
