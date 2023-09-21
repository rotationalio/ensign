import { t, Trans } from '@lingui/macro';
import { Heading, Loader, Toast } from '@rotational/beacon-core';
import { Suspense, useState } from 'react';

import { SentryErrorBoundary } from '@/components/Error';
import SettingsButton from '@/components/ui/Settings/Settings';

import { useFetchProfile } from '../hooks/useFetchProfile';
import { CancelAcctModal } from './CancelModal';
// import ChangePasswordModal from './ChangePassword/ChangePasswordModal';
import MemberDetailInfo from './MemberInfo';
export default function MemberDetails() {
  const [isCancelModalOpen, setIsCancelModalOpen] = useState(false);
  // const [isChangePasswordModalOpen, setIsChangePasswordModalOpen] = useState(false);
  const { profile, hasProfileFailed, isFetchingProfile, error } = useFetchProfile();

  if (isFetchingProfile) {
    return <Loader size="lg" />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasProfileFailed}
        variant="danger"
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

  // const onCloseChangePasswordModal = () => {
  //   setIsChangePasswordModalOpen(false);
  // };

  // const openChangePasswordModal = () => {
  //   setIsChangePasswordModalOpen(true);
  // };

  return (
    <>
      <Suspense fallback={<Loader />}>
        <SentryErrorBoundary
          fallback={<div>We are unable to fetch your member, please try again.</div>}
        >
          <div className="my-10">
            <div className="flex items-center justify-between rounded-md bg-[#F7F9FB] px-6 py-3">
              <Heading as="h1" className="flex w-40 items-center gap-5 text-2xl font-semibold">
                <span className="mr-2">
                  <Trans>User Profile</Trans>
                </span>
              </Heading>
              <div className="flex w-full justify-end">
                <SettingsButton
                  key="org-action"
                  data={[
                    {
                      name: t`Cancel Account`,
                      onClick: openCancelModal,
                    },
                    // {
                    //   name: t`Change Password`,
                    //   onClick: () => openChangePasswordModal(),
                    // },
                  ]}
                />
              </div>
            </div>
            <div className="mx-6">
              <MemberDetailInfo data={profile} />
            </div>
          </div>

          <CancelAcctModal close={onCloseCancelModal} isOpen={isCancelModalOpen} />
          {/* <ChangePasswordModal
            open={isChangePasswordModalOpen}
            handleModalClose={onCloseChangePasswordModal}
          /> */}
        </SentryErrorBoundary>
      </Suspense>
    </>
  );
}
