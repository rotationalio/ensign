import { Meta, Story } from '@storybook/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();
import TenantTable from './TenantTable';

export default {
  title: 'tenants/TenantTable',
  component: TenantTable,
} as Meta;

const Template: Story = (args) => {
  return (
    <QueryClientProvider client={queryClient}>
      <TenantTable {...args} />
    </QueryClientProvider>
  );
};

export const Default = Template.bind({});
Default.args = {
  tenants: [
    {
      id: '1',
      name: 'Tenant 1',
      env: 'dev',
      cloud: 'aws',
      region: 'us-east-1',
      created: '2021-09-01',
    },
    {
      id: '2',
      name: 'Tenant 2',
      env: 'dev',
      cloud: 'aws',
      region: 'us-east-1',
      created: '2021-09-01',
    },
  ],
};
