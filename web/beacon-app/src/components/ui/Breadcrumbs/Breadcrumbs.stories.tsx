import { Meta, Story } from '@storybook/react';

import Breadcrumbs from './Breadcrumbs';

export default {
  title: 'ui/Breadcrumbs',
  component: Breadcrumbs,
} as Meta<typeof Breadcrumbs>;

const Template: Story<typeof Breadcrumbs> = (args) => (
  <Breadcrumbs {...args}>
    <Breadcrumbs.Item>Item 1</Breadcrumbs.Item>
    <Breadcrumbs.Item>Item 2</Breadcrumbs.Item>
    <Breadcrumbs.Item>Item 2</Breadcrumbs.Item>
  </Breadcrumbs>
);

export const Default = Template.bind({});
Default.args = {};
