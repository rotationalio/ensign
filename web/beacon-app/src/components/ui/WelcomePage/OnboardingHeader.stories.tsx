import { Meta, Story } from '@storybook/react'
import OnboardingHeader from './OnboardingHeader';

export default {
    title: 'ui/Welcome/OnboardingHeader',
    component: OnboardingHeader,
  } as Meta;
  
  const Template: Story = (args) => <OnboardingHeader {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};