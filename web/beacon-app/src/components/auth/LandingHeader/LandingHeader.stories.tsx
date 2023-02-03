import { Meta, Story } from '@storybook/react';

import { LandingHeader } from '.';

export default {
  title: 'component/landing-page/Header',
} as Meta;

const Template: Story = (args) => <LandingHeader {...args} />;

export const Default = Template.bind({});
Default.args = {};
