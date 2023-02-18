import React, { useState } from 'react';
import { Meta, Story } from '@storybook/react';

import DeleteOrgModal from './DeleteOrgModal';

export default {
  title: 'organizations/DeleteOrgModal',
} as Meta;

interface DeleteOrgModalProps {
  close: () => void;
  isOpen: boolean;
}

const [isOpen, setIsOpen] = useState(false);

const handleClose = () => setIsOpen(false);



const Template: Story<DeleteOrgModalProps> = (args) => <DeleteOrgModal {...args} />;

export const Default = Template.bind({});
Default.args = {
  close: handleClose,
  isOpen: isOpen,

};
