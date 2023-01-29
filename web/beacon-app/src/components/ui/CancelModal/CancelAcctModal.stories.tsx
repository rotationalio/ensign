import { Meta, Story } from '@storybook/react';
import CancelAcctModal from './CancelAcctModal';

export default {
  title: '/component/CancelAcctModal',
} as Meta;

const Template: Story = (args) => <CancelAcctModal {...args} />;

export const Default = Template.bind({});
Default.args = {};