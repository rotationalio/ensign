import { Meta, Story } from '@storybook/react';

import Avatar from './Avatar';
import { AvatarProps } from './Avatar.type';

export default {
  title: 'ui/Avatar',
  component: Avatar,
  argTypes: {
    src: {
      description: 'The src attribute for the img element.',
      control: {
        type: 'text',
      },
    },
    srcSet: {
      description:
        'The srcSet attribute for the img element. Use this attribute for responsive image display.',
      control: {
        type: 'text',
      },
    },
    alt: {
      description:
        'Used in combination with src or srcSet to provide an alt attribute for the rendered img element.',
      control: {
        type: 'text',
      },
    },
  },
} as Meta<AvatarProps>;

const Template: Story<AvatarProps> = (args) => <Avatar {...args} />;

export const Default = Template.bind({});
Default.args = {
  src: 'https://images.unsplash.com/photo-1492633423870-43d1cd2775eb?&w=128&h=128&dpr=2&q=80',
  alt: 'Colm Tuite',
};
