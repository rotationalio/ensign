import { Meta, Story } from '@storybook/react'
import TenantQuickStart from './TenantQuickStart';

export default {
    title: 'ui/Welcome/TenantQuickStart',
    component: TenantQuickStart,
  } as Meta;
  
  const Template: Story = (args) => <TenantQuickStart {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};