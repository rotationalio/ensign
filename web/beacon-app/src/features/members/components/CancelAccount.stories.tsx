import { Meta, Story } from '@storybook/react';

import CancelAccount from './CancelAccount';

export default {
  title: 'members/CancelAccount',
  component: CancelAccount,
} as Meta;

const Template: Story = (args) => <CancelAccount {...args} />;

export const Default = Template.bind({});
Default.args = {};
