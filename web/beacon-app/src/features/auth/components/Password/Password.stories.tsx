import { Meta, Story } from '@storybook/react';

import Password from './Password';

export default {
  title: 'auth/Password',
  component: Password,
} as Meta;

const Template: Story = (args) => <Password {...args} />;

export const Default = Template.bind({});
Default.args = {};
