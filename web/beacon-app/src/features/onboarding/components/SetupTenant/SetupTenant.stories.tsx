import { Meta, Story } from '@storybook/react';

import SetupTenant from './SetupTenant';

export default {
  title: 'onboarding/SetupTenant/SetupTenant',
  component: SetupTenant,
} as Meta;

const Template: Story = (args) => <SetupTenant {...args} />;
export const Default = Template.bind({});
Default.args = {};
