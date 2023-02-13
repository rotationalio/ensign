import { Meta, Story } from '@storybook/react';

import { CardListItem, CardListItemProps } from '.';

export default {
  title: 'component/common/CardListItem',
  component: CardListItem,
} as Meta;

const Template: Story<CardListItemProps> = (args) => <CardListItem {...args} />;

export const Default = Template.bind({});
Default.args = {
  data: [
    {
      label: 'Project Name:',
      value: 'Acme Systems Global Event Stream',
    },
    {
      label: 'Project ID:',
      value: '1234567890',
    },
    {
      label: 'Project Description:',
      value: 'This is a description of the project',
    },
  ],
  tableClassName: 'border-spacing-y-10',
};
