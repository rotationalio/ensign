import { Meta, Story } from '@storybook/react';
import LandingFooter from './LandingFooter';

export default {
  title: 'component/landing-page/Footer',
  component: LandingFooter,
} as Meta;

const Template: Story = (args) => <LandingFooter {...args} />;

export const Default = Template.bind({});
Default.args = {};
