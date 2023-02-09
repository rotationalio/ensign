import { Meta, Story } from '@storybook/react';

import DeleteOrgModal from './DeleteOrgModal';

export default {
  title: 'organizations/DeleteOrgModal',
} as Meta;

const Template: Story = (args) => <DeleteOrgModal {...args} />;

export const Default = Template.bind({});
Default.args = {};
