import { Meta, Story } from '@storybook/react';
import MaintenanceMode from './MaintenanceMode';

export default {
  title: '/component/MaintenanceMode',
} as Meta;

const Template: Story = (args) => <MaintenanceMode {...args} />;

export const Default = Template.bind({});
Default.args = {};