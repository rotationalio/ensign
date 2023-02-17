import { render } from '@testing-library/react';
import { vi } from 'vitest';

import { I18nWrapper } from '@/utils/test-utils';

import LoginForm from '../LoginForm';

describe('first', () => {
  it('should submit the form', async () => {
    // await act(async () => {
    //   dynamicActivate('en');
    // });
    const handleSubmit = vi.fn();
    render(<LoginForm onSubmit={handleSubmit} />, {
      wrapper: I18nWrapper,
    });

    // // eslint-disable-next-line testing-library/no-unnecessary-act
    // await act(async () => {
    //   await userEvent.type(screen.getByTestId('email'), 'john.doe@example.com');
    //   await userEvent.type(screen.getByTestId('password'), 'Password123@@_');
    // });
    //expect(handleSubmit).toHaveBeenCalledTimes(1);
  });
});
