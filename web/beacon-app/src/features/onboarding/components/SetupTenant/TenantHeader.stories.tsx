import { Meta, Story } from '@storybook/react';

import TenantHeader from './TenantHeader';

export default {
  title: 'onboarding/SetupTenant/TenantHeader',
  component: TenantHeader,
} as Meta;

const Template: Story = (args) => <TenantHeader {...args} />;

export const Default = Template.bind({});
Default.args = {};
