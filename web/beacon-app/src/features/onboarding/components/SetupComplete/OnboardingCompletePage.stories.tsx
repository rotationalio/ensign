import { Meta, Story } from '@storybook/react';

import OnboardingCompletePage from './OnboardingCompletePage';

export default {
  title: 'onboarding/SetupTenantComplete/OnboardingComplete',
  component: OnboardingCompletePage,
} as Meta;

const Template: Story = (args) => <OnboardingCompletePage {...args} />;

export const Default = Template.bind({});
Default.args = {};
