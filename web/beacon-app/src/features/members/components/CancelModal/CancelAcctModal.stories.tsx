import { Meta, Story } from '@storybook/react';

import CancelAcctModal from './CancelAcctModal';
interface CancelAcctModalProps {
  close: () => void;
  isOpen: boolean;
}
export default {
  title: 'beacon/CancelAcctModal',
} as Meta;

const Template: Story<CancelAcctModalProps> = (args) => <CancelAcctModal {...args} />;

export const Default = Template.bind({});
Default.args = {};
