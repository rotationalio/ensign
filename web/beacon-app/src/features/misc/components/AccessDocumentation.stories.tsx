import { Meta, Story } from '@storybook/react';

import AccessDocumentation from './AccessDocumentation';

export default {
  title: 'features/misc/components/AccessDocumentation',
  component: AccessDocumentation,
} as Meta<typeof AccessDocumentation>;

const Template: Story<typeof AccessDocumentation> = (args) => <AccessDocumentation {...args} />;

export const Default = Template.bind({});
Default.args = {};
