import { Meta, Story } from '@storybook/react';
import DeleteOrg from './DeleteOrg';

export default {
  title: 'organizations/DeleteOrg',
  component: DeleteOrg,
} as Meta;

const Template: Story = (args) => <DeleteOrg {...args} />;

export const Default = Template.bind({});
Default.args = {};