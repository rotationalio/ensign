import { Meta, Story } from '@storybook/react';
import MaintenaceMode from './MaintenaceMode';

export default {
  title: '/component/MaintenanceMode',
} as Meta;

const Template: Story = (args) => <MaintenaceMode {...args} />;

export const Default = Template.bind({});
Default.args = {};