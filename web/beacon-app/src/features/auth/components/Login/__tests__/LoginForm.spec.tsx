import React from 'react';
import { vi } from 'vitest';

import { customRender } from '../../../../../utils/test-utils';
import LoginForm from '../LoginForm';

describe('LoginForm', () => {
  it('should submit the form', async () => {
    const handleSubmit = vi.fn();
    customRender(<LoginForm onSubmit={handleSubmit} />);
  });
});
