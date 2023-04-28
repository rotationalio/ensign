/* eslint-disable testing-library/no-node-access */
import { act } from '@testing-library/react';

import { dynamicActivate } from '../../../../I18n';

describe('QuickView', () => {
  beforeEach(() => {
    act(() => {
      dynamicActivate();
    });
  });
  it.todo('should display right data');
  // it.skip('should display right data', () => {
  //   const data = [
  //     {
  //       name: 'Active Projects',
  //       value: 4,
  //     },
  //     {
  //       name: 'Topics',
  //       value: 1,
  //     },
  //     {
  //       name: 'API Keys',
  //       value: 3,
  //     },
  //     {
  //       name: 'Data Storage',
  //       value: 2,
  //     },
  //   ];

  //   customRender(<QuickView data={data} />);
  //   expect(screen.getByRole('heading', { name: /active projects/i })).toBeInTheDocument();
  //   expect(screen.getByRole('heading', { name: /active projects/i }).nextSibling?.textContent).toBe(
  //     '4'
  //   );

  //   expect(screen.getByRole('heading', { name: /topics/i })).toBeInTheDocument();
  //   expect(screen.getByRole('heading', { name: /topics/i }).nextSibling?.textContent).toBe('1');

  //   expect(screen.getByRole('heading', { name: /api keys/i })).toBeInTheDocument();
  //   expect(screen.getByRole('heading', { name: /api keys/i }).nextSibling?.textContent).toBe('3');

  //   expect(screen.getByRole('heading', { name: /data storage/i })).toBeInTheDocument();
  //   expect(screen.getByRole('heading', { name: /data storage/i }).nextSibling?.textContent).toBe(
  //     '2'
  //   );
  // });
});
