import { render, screen } from '@testing-library/react';
import React from 'react';
import { vi } from 'vitest';

import { QueryClientWrapper } from '../../../../utils/test-utils';
import ProjectList from '../ProjectList';

vi.mock('react-router-dom', async () => ({
  ...vi.importMock('react-router-dom'),
  Link: ({ children, to, ...rest }: { children: JSX.Element; to: string }) =>
    React.createElement('a', { href: to, ...rest }, children),
}));

describe('<ProjectList />', () => {
  it('should disable add project button', () => {
    render(<ProjectList />, { wrapper: QueryClientWrapper });

    expect(screen.getByTestId('create__project-btn')).toBeDisabled();
  });
});
