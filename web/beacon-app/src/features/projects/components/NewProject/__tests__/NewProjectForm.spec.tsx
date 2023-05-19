import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it, vi } from 'vitest';

// import selectEvent from 'react-select-event';
import NewProjectForm from '../NewProjectForm';
import type { NewProjectFormProps } from './NewProjectForm';

const renderComponent = (props) => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  return render(<NewProjectForm {...props} />, { wrapper });
};

vi.mock('@lingui/macro', () => ({
  t: (str) => str,
  Trans: ({ children }) => children,
}));

vi.mock('@/features/projects/hooks/useNewProjectForm', () => ({
  useNewProjectForm: vi.fn(),
}));

describe('NewProjectForm', () => {
  beforeAll(() => {
    vi.spyOn(React, 'useEffect').mockImplementation((f) => f());
  });
  it(' should render the form', () => {
    const props: NewProjectFormProps = {
      onSubmit: vi.fn(),
    };
    renderComponent(props);
    expect(screen.getByTestId('project-name')).toBeInTheDocument();
  });

  it('should be able submit the form', async () => {
    const props = {
      onSubmit: vi.fn(),
    };

    renderComponent(props);

    const input = screen.getByTestId('project-name');
    const submit = screen.getByTestId('prj-submit-btn');

    await (act as any)(async () => {
      fireEvent.change(input, { target: { value: 'test' } });
      fireEvent.click(submit);
    });

    expect(props.onSubmit).toHaveBeenCalled();
  });

  it('should disable submit button when isDisabled is true', async () => {
    const props = {
      onSubmit: vi.fn(),
      isDisabled: true,
    };

    renderComponent(props);

    const submit = screen.getByTestId('prj-submit-btn');

    expect(submit).toBeDisabled();
  });

  it('should display error message if name is empty', async () => {
    const props = {
      onSubmit: vi.fn(),
    };

    renderComponent(props);

    const input = screen.getByTestId('project-name');
    const submit = screen.getByTestId('prj-submit-btn');

    await (act as any)(async () => {
      fireEvent.change(input, { target: { value: '' } });
      fireEvent.click(submit);
    });

    expect(screen.getByText('Project name is required.')).toBeInTheDocument();
  });
});
