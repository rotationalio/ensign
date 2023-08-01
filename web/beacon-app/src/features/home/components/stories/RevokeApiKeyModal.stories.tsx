import { Meta, Story } from '@storybook/react';

import RevokeApiKeyModal, { RevokeApiKeyModalProps } from '../RevokeApiKeyModal';

export default {
  title: 'beacon/RevokeApiKeyModal',
} as Meta<RevokeApiKeyModalProps>;

const Template: Story<RevokeApiKeyModalProps> = (args) => <RevokeApiKeyModal {...args} />;

export const Default = Template.bind({});
Default.args = {
  open: true,
};
