import { Meta, Story } from '@storybook/react';
import OrganizationsTable from './OrganizationsTable';

export default {
  title: 'organizations/OrganizationsTable',
  component: OrganizationsTable,
} as Meta;

const Template: Story = (args) => <OrganizationsTable {...args} />;

export const Default = Template.bind({});
Default.args = {};