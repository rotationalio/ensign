import { Meta, Story } from '@storybook/react';

import SideBar from './Sidebar';

export default {
  title: 'component/layout/SideBar',
  component: SideBar,
} as Meta<typeof SideBar>;

const Template: Story<typeof SideBar> = (args) => (
  <div>
    <SideBar {...args} />
  </div>
);

export const Default = Template.bind({});
Default.args = {};
