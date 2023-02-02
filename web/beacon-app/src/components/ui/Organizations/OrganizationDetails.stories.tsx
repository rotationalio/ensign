import { Meta, Story } from '@storybook/react';
import OrganizationDetails from './OrganizationDetails';

export default {
  title: 'ui/OrganizationDetails',
  component: OrganizationDetails,
} as Meta;

const Template: Story = (args) => <OrganizationDetails {...args} />;

export const Default = Template.bind({});
Default.args = {};