import { Meta, Story } from '@storybook/react';

import ProjectDetailsStep from './ProjectDetailsStep';
// interface ProjectDetailsStepProps {
//   tenantID: string;
// }

export default {
  title: 'features/misc/components/ProjectDetails',
  component: ProjectDetailsStep,
} as Meta<typeof ProjectDetailsStep>;

const Template: Story = (args) => <ProjectDetailsStep {...args} />;

export const Default = Template.bind({});
Default.args = {};
