import { Meta, Story } from '@storybook/react'
import SetupYourTenant from './SetupYourTenant';

export default {
    title: 'ui/Welcome/SetupYourTenant',
    component: SetupYourTenant,
  } as Meta;
  
  const Template: Story = (args) => <SetupYourTenant {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};