import { Meta, Story } from '@storybook/react';

import Footer from './Footer';

export default {
  title: 'LandingFooter',
  component: Footer,
} as Meta;

const Template: Story = (args) => <Footer {...args} />;

export const Default = Template.bind({});
Default.args = {};
