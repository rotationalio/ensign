/* eslint-disable testing-library/no-node-access */
import { render, screen } from '@testing-library/react';
import React from 'react';

import QuickView from '../QuickView';

describe('QuickView', () => {
  it('should display right data', () => {
    const data = [
      {
        name: 'Active Projects',
        value: 4,
      },
      {
        name: 'Topics',
        value: 1,
      },
      {
        name: 'API Keys',
        value: 3,
      },
      {
        name: 'Data Storage',
        value: 2,
      },
    ];

    render(<QuickView data={data} />);
    expect(screen.getByRole('heading', { name: /active projects/i })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /active projects/i }).nextSibling?.textContent).toBe(
      '4'
    );

    expect(screen.getByRole('heading', { name: /topics/i })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /topics/i }).nextSibling?.textContent).toBe('1');

    expect(screen.getByRole('heading', { name: /api keys/i })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /api keys/i }).nextSibling?.textContent).toBe('3');

    expect(screen.getByRole('heading', { name: /data storage/i })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /data storage/i }).nextSibling?.textContent).toBe(
      '2'
    );
  });
});
