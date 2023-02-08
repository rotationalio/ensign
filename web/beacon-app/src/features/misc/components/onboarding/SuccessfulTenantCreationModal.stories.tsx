import { Meta, Story } from '@storybook/react';

import SuccessfulTenantCreationModal, {
  SuccessfulTenantCreationModalProps,
} from './SuccessfulTenantCreationModal';

export default {
  title: 'beacon/SuccessFullTenantModal',
} as Meta<SuccessfulTenantCreationModalProps>;

const Template: Story<SuccessfulTenantCreationModalProps> = (args) => (
  <SuccessfulTenantCreationModal {...args} />
);

export const Default = Template.bind({});
Default.args = {
  open: true,
};
