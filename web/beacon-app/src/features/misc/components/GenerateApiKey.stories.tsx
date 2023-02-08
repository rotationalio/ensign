import { Meta, Story } from '@storybook/react';

import GenerateApiKey from './GenerateApiKey';

export default {
  title: 'features/misc/components/GenerateApiKey',
  component: GenerateApiKey,
} as Meta<typeof GenerateApiKey>;

const Template: Story<typeof GenerateApiKey> = (args) => <GenerateApiKey {...args} />;

export const Default = Template.bind({});
Default.args = {};
