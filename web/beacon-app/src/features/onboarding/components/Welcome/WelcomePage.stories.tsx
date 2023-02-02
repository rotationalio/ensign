import { Meta, Story } from '@storybook/react';
import WelcomePage from "./WelcomePage";

export default {
    title: 'onboarding/Welcome/WelcomePage',
    component: WelcomePage,
  } as Meta;
  
  const Template: Story = (args) => <WelcomePage {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};