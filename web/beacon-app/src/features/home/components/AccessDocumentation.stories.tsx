import { Meta, Story } from '@storybook/react';

import AccessDocumentationStep from './AccessDocumentationStep';

export default {
  title: 'features/misc/components/AccessDocumentation',
  component: AccessDocumentationStep,
} as Meta<typeof AccessDocumentationStep>;

const Template: Story<typeof AccessDocumentationStep> = (args) => (
  <AccessDocumentationStep {...args} />
);

export const Default = Template.bind({});
Default.args = {};
