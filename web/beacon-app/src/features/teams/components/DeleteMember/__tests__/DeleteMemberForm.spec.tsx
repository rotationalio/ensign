import { fireEvent, render, screen } from '@testing-library/react';
import React from 'react';
import { vi } from 'vitest';

import DeleteMemberForm from '../DeleteMemberForm';

const renderComponent = () => {
  const props = {
    onSubmit: vi.fn(),
    initialValues: {
      id: '1',
      name: 'John Doe',
      delete_agreement: false,
    },
    isSubmitting: false,
  };

  render(<DeleteMemberForm {...props} />);
};

describe('DeleteMemberForm', () => {
  it('agree checkbox should be checked', () => {
    renderComponent();
    const has_agreed = screen.getByTestId('delete_agreement');
    fireEvent.click(has_agreed);
    expect(has_agreed).toBeChecked();
  });

  it('agree checkbox should be unchecked by default', () => {
    renderComponent();
    const has_agreed = screen.getByTestId('delete_agreement');

    expect(has_agreed).not.toBeChecked();
  });

  it(' should be disabled when agree checkbox is unchecked', () => {
    renderComponent();

    const removeBtn = screen.getByTestId('remove-btn');
    expect(removeBtn).toBeDisabled();
  });
});
