import { Meta, Story } from '@storybook/react';

import Loader from './Loader';

export default {
  title: 'ui/Loader',
  component: Loader,
} as Meta<typeof Loader>;

const Template: Story<typeof Loader> = (args) => <Loader {...args} />;

export const Default = Template.bind({});
Default.args = {
  label: 'Initializing...Allocating Resources...Finalizing Tenant',
};
