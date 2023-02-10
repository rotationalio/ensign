import { Meta, Story } from '@storybook/react';
import TenantTable from './TenantTable';

export default {
  title: 'tenants/TenantTable',
  component: TenantTable,
} as Meta;

const Template: Story = (args) => <TenantTable {...args} />;

export const Default = Template.bind({});
Default.args = {};