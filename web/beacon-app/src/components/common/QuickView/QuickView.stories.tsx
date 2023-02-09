import { Meta, Story } from '@storybook/react';

import { QuickView, QuickViewProps } from '.';

export default {
  title: 'component/common/QuickView',
  component: QuickView,
} as Meta;

const Template: Story<QuickViewProps> = (args) => <QuickView {...args} />;

export const Default = Template.bind({});
Default.args = {
  data: [
    {
      name: 'Active Projects',
      value: 10,
    },
    {
      name: 'Topics',
      value: 10,
    },
    {
      name: 'API Keys',
      value: 10,
    },
    {
      name: 'Data Storage (GB)',
      value: 10,
    },
  ],
};
