import { Meta, Story } from '@storybook/react';

import AccessDashboard from './AccessDashboard';

export default {
  title: 'ui/AccessDashboard',
  component: AccessDashboard,
} as Meta;

const Template: Story = (args) => <AccessDashboard {...args} />;

export const Default = Template.bind({});
Default.args = {};
