import { Button, Toast, useMenu } from '@rotational/beacon-core';
import { useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import { BlueBars } from '@/components/icons/blueBars';
import { Dropdown as Menu } from '@/components/ui/Dropdown';
import { useOrgStore } from '@/store';

import { useFetchOrg } from '../hooks/useFetchOrgDetail';
import { getOrgData } from '../utils';
import { DeleteOrgModal } from './DeleteOrgModal';

export default function OrganizationDetails() {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'org-action' });
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const getOrg = useOrgStore.getState() as any;

  const { org: organization, hasOrgFailed, isFetchingOrg, error } = useFetchOrg(getOrg.org);

  if (isFetchingOrg) {
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasOrgFailed}
        variant="danger"
        title="We are unable to fetch your organization, please try again."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const onCloseDeleteModal = () => {
    setIsDeleteModalOpen(false);
  };

  const openDeleteModal = () => {
    setIsDeleteModalOpen(true);
  };

  return (
    <>
      <CardListItem data={getOrgData(organization)} className="my-5">
        <div className="flex w-full justify-end">
          <Button
            variant="ghost"
            className="bg-transparent flex justify-end border-none outline-none focus:outline-none "
            onClick={open}
            size="xsmall"
          >
            <BlueBars />
          </Button>

          <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
            <Menu.Item onClick={openDeleteModal}>Delete</Menu.Item>
          </Menu>
        </div>
        <DeleteOrgModal close={onCloseDeleteModal} isOpen={isDeleteModalOpen} />
      </CardListItem>
    </>
  );
}
