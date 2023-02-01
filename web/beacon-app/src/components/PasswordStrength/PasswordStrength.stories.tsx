import { Meta, Story } from '@storybook/react';

import PasswordStrength from './PasswordStrength';

interface Props {
  string: string;
}

export default {
  title: 'component/Common/PasswordStrength',
  component: PasswordStrength,
} as Meta;

const Template: Story<Props> = (args) => <PasswordStrength {...args} />;

export const Default = Template.bind({});
Default.args = {
  string: '1Password@',
};
