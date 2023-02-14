import { Meta, Story } from '@storybook/react';

import MemberDetails from './MemberDetails';

export default {
  title: 'members/MemberDetails',
  component: MemberDetails,
} as Meta;

const Template: Story = (args) => <MemberDetails {...args} />;

export const Default = Template.bind({});
Default.args = {};
