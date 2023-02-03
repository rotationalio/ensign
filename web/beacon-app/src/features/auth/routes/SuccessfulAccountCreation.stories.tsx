import { Meta, Story } from '@storybook/react';

import SuccessfulAccountCreation from './SuccessfulAccountCreation';

export default {
  title: 'features/auth/components/SuccessfulAccountCreation',
  component: SuccessfulAccountCreation,
} as Meta<typeof SuccessfulAccountCreation>;

const Template: Story<typeof SuccessfulAccountCreation> = (args) => (
  <SuccessfulAccountCreation {...args} />
);

export const Default = Template.bind({});
Default.args = {};
