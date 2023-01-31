import { Meta, Story } from '@storybook/react';

import LeftSideBar from './LeftSidebar';

export default {
  title: 'component/layout/LeftSideBar',
  component: LeftSideBar,
} as Meta<typeof LeftSideBar>;

const Template: Story<typeof LeftSideBar> = (args) => (
  <div>
    <LeftSideBar {...args} />
  </div>
);

export const Default = Template.bind({});
Default.args = {};
