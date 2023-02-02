import AccessDashboard from "./AccessDashboard";
import { Meta, Story } from "@storybook/react";


export default {
    title: 'ui/AccessDashboard',
    component: AccessDashboard,
  } as Meta;
  
  const Template: Story = (args) => <AccessDashboard {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};