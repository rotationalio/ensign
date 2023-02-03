import { Meta, Story } from '@storybook/react';

import Registration from './Registration';

export default {
  title: 'features/auth/routes/Registration',
  component: Registration,
} as Meta<typeof Registration>;

const Template: Story<typeof Registration> = (args) => <Registration {...args} />;

export const Default = Template.bind({});
Default.args = {};
