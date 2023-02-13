import { Meta, Story } from '@storybook/react';

import { TableHeading, TableHeadingProps } from '.';

export default {
  title: 'component/common/TableHeading',
  component: TableHeading,
} as Meta;

const Template: Story<TableHeadingProps> = (args) => <TableHeading {...args} />;

export const Default = Template.bind({});

Default.args = {
  children: 'Table Heading',
};
