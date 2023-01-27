import { Meta, Story } from '@storybook/react';

import LandingFooter from './Footer';

export default {
  title: 'LandingFooter',
  component: LandingFooter,
} as Meta;

const Template: Story = (args) => <LandingFooter {...args} />;

export const Default = Template.bind({});
Default.args = {};
