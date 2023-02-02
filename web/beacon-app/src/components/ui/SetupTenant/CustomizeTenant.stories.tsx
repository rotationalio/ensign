import { Meta, Story } from '@storybook/react'
import CustomizeTenant from './CustomizeTenant';

export default {
    title: 'ui/Welcome/CustomizeTenant',
    component: CustomizeTenant,
  } as Meta;
  
  const Template: Story = (args) => <CustomizeTenant {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};