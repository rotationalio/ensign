import { Meta, Story } from '@storybook/react';

import Tag, { TagProps } from './Tag';

export default {
  title: 'ui/Tag',
  component: Tag,
  argTypes: {
    variant: {
      description: 'The variant of the tag.',
      control: {
        type: 'select',
        options: ['primary', 'secondary', 'success', 'warning', 'danger', 'info', 'light', 'dark'],
      },
    },
    size: {
      description: 'The size of the tag.',
      control: {
        type: 'select',
        options: ['small', 'medium', 'large'],
      },
    },
    leftIcon: {
      description: 'The icon to display on the left side of the tag.',
      control: {
        type: 'text',
      },
    },
    rightIcon: {
      description: 'The icon to display on the right side of the tag.',
      control: {
        type: 'text',
      },
    },
  },
} as Meta<TagProps>;

const Template: Story<TagProps> = (args) => <Tag {...args} />;
export const Default = Template.bind({});
Default.args = {
  children: 'Tag',
};
export const Primary = Template.bind({});
Primary.args = {
  children: 'Tag',
  variant: 'primary',
};
