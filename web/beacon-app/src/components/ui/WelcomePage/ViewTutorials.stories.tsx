import { Meta, Story } from '@storybook/react'

import ViewTutorials from "./ViewTutorials";

export default {
    title: 'ui/Welcome/ViewTutorials',
    component: ViewTutorials,
  } as Meta;
  
  const Template: Story = (args) => <ViewTutorials {...args} />;
  
  export const Default = Template.bind({});
  Default.args = {};