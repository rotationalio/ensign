import React from 'react';
import { vi } from 'vitest';

import { dynamicActivate } from '../../../../I18n';

vi.mock('react-router-dom', async () => ({
  ...vi.importMock('react-router-dom'),
  Link: ({ children, to, ...rest }: { children: JSX.Element; to: string }) =>
    React.createElement('a', { href: to, ...rest }, children),
}));

describe('<ProjectList />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it.todo(
    'should disable add project button'
    // , () => {
    //   const { debug } = customRender(<ProjectList />);
    //   debug();
    //   expect(screen.getByTestId('create__project-btn')).toBeDisabled();
    // }
  );
});
