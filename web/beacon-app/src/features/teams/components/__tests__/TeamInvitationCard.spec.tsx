import { render, screen } from '@testing-library/react';
import React from 'react';

import TeamInvitationCard from '../TeamInvitationCard';

describe('<TeamInvitationCard />', () => {
  it('should display right data', () => {
    const data = {
      email: 'fylipagy@mailinator.com',
      org_name: 'Cameron Mayer Trading',
      inviter_name: 'Elijah Merritt',
      role: 'Observer',
      has_account: false,
    };

    render(<TeamInvitationCard data={data} />);

    expect(screen.getByTestId('inviter_name').textContent).toBe(data.inviter_name);
    expect(screen.getByTestId('org_name').textContent).toBe(data.org_name);
    expect(screen.getByTestId('role').textContent).toBe(data.role);
  });

  it('should display correct message for new invited member', () => {
    const data = {
      email: 'fylipagy@mailinator.com',
      org_name: 'Cameron Mayer Trading',
      inviter_name: 'Elijah Merritt',
      role: 'Observer',
      has_account: false,
    };

    const { container } = render(<TeamInvitationCard data={data} />);

    expect(container.textContent).toBe(
      `You've Been Invited!You've been invited by ${data.inviter_name} to join the ${data.org_name} organization as ${data.role} on Ensign! Create your account today.`
    );
  });

  it('should display correct message for existing invited member', () => {
    const data = {
      email: 'fylipagy@mailinator.com',
      org_name: 'Cameron Mayer Trading',
      inviter_name: 'Elijah Merritt',
      role: 'Observer',
      has_account: true,
    };

    const { container } = render(<TeamInvitationCard data={data} />);

    expect(container.textContent).toBe(
      `You've Been Invited!You've been invited by ${data.inviter_name} to join the ${data.org_name} organization as ${data.role} on Ensign! Log in to accept the invitation.`
    );
  });
});
