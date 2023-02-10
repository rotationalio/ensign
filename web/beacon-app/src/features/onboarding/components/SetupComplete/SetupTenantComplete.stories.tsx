import { Meta, Story } from '@storybook/react';

import SetupTenantComplete from './SetupTenantComplete';

export default {
  title: 'onboarding/SetupTenantComplete/SetupTenantComplete',
  component: SetupTenantComplete,
} as Meta;

const Template: Story = (args) => <SetupTenantComplete {...args} />;

export const Default = Template.bind({});
Default.args = {};
