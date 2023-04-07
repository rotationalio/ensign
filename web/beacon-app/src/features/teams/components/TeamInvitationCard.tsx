import { Card } from '@rotational/beacon-core';
import { useLocation } from 'react-router-dom';

import { ROUTES } from '@/application';

export default function TeamInvitationCard() {
  const location = useLocation();
  const isNewInvitePage = location.pathname == ROUTES.NEW_INVITATION;
  const isExistingInvitePage = location.pathname == ROUTES.EXISTING_INVITATION;

  return (
    <>
      <Card className="bg-[#ECF6FF]">
        <div className="ml-5 p-6">
          <Card.Header>
            <h1 className="mb-4 text-base font-bold">You've Been Invited!</h1>
          </Card.Header>

          {isNewInvitePage && (
            <Card.Body>
              <p>
                You've been invited by <span className="font-bold">(inviter name)</span> to join the{' '}
                <span className="font-bold">(org name)</span> organization as{' '}
                <span className="font-bold">(role)</span> on Ensign! Create your account today.
              </p>
            </Card.Body>
          )}

          {isExistingInvitePage && (
            <Card.Body>
              <p>
                You've been invited by <span className="font-bold">(inviter name)</span> to join the{' '}
                <span className="font-bold">(org name)</span> organization as{' '}
                <span className="font-bold">(role)</span> on Ensign! Log in to accept the
                invitation.
              </p>
            </Card.Body>
          )}
        </div>
      </Card>
    </>
  );
}
