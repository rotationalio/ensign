import { Button, Loader, Toast, useMenu } from '@rotational/beacon-core';
import { Suspense, useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { SentryErrorBoundary } from '@/components/Error';
import { BlueBars } from '@/components/icons/blueBars';
import { Dropdown as Menu } from '@/components/ui/Dropdown';
import { formatMemberData } from '@/features/members/utils';
import { useOrgStore } from '@/store';

import { useFetchMember } from '../hooks/useFetchMember';
import { CancelAcctModal } from './CancelModal';
export default function MemberDetails() {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'org-action' });
  const [isCancelModalOpen, setIsCancelModalOpen] = useState(false);

  const orgDataState = useOrgStore.getState() as any;

  const { member, hasMemberFailed, isFetchingMember, error } = useFetchMember(orgDataState?.user);

  if (isFetchingMember) {
    return <Loader size="lg" />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasMemberFailed}
        variant="danger"
        title="We are unable to fetch your member, please try again."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const onCloseCancelModal = () => {
    setIsCancelModalOpen(false);
  };

  const openCancelModal = () => {
    setIsCancelModalOpen(true);
  };

  return (
    <>
      <Suspense fallback={<Loader />}>
        <SentryErrorBoundary
          fallback={<div>We are unable to fetch your member, please try again.</div>}
        >
          <CardListItem data={formatMemberData(member)} className="my-5">
            <div className="flex w-full justify-end">
              <Button
                variant="ghost"
                className="bg-transparent flex justify-end border-none"
                onClick={open}
              >
                <BlueBars />
              </Button>

              <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
                <Menu.Item onClick={openCancelModal}>Cancel Account</Menu.Item>
              </Menu>
            </div>
            <CancelAcctModal close={onCloseCancelModal} isOpen={isCancelModalOpen} />
          </CardListItem>
        </SentryErrorBoundary>
      </Suspense>
    </>
  );
}
