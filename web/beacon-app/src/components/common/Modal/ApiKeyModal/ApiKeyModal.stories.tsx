import { Meta, Story } from '@storybook/react';

import ApiKeyModal, { ApiKeyModalProps } from './ApiKeyModal';

export default {
  title: 'beacon/ApiKeyModal',
} as Meta<ApiKeyModalProps>;

const Template: Story<ApiKeyModalProps> = (args) => <ApiKeyModal {...args} />;

export const Default = Template.bind({});
Default.args = {
  open: true,
};
