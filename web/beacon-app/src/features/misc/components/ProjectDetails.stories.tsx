import { Meta, Story } from '@storybook/react';

import ProjectDetails from './ProjectDetails';

export default {
  title: 'features/misc/components/ProjectDetails',
  component: ProjectDetails,
} as Meta<typeof ProjectDetails>;

const Template: Story<typeof ProjectDetails> = (args) => <ProjectDetails {...args} />;

export const Default = Template.bind({});
Default.args = {};
