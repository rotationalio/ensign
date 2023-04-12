import { Card } from '@rotational/beacon-core';

export default function TeamInvitationCard({ data }: { data: any }) {
  const hasAccount = data?.has_account || false;

  return (
    <>
      <Card className="overflow-hidden !rounded-lg bg-[#ECF6FF] px-8 py-5">
        <Card.Header>
          <h1 className="mb-4 text-base font-bold">You've Been Invited!</h1>
        </Card.Header>

        {hasAccount ? (
          <Card.Body>
            <p>
              You've been invited by{' '}
              <span className="font-bold" data-testid="inviter_name">
                {data?.inviter_name}
              </span>{' '}
              to join the{' '}
              <span className="font-bold" data-testid="org_name">
                {data?.org_name}
              </span>{' '}
              organization as{' '}
              <span className="font-bold" data-testid="role">
                {data?.role}
              </span>{' '}
              on Ensign! Log in to accept the invitation.
            </p>
          </Card.Body>
        ) : (
          <Card.Body>
            <p>
              You've been invited by{' '}
              <span className="font-bold" data-testid="inviter_name">
                {data?.inviter_name}
              </span>{' '}
              to join the{' '}
              <span className="font-bold" data-testid="org_name">
                {data?.org_name}
              </span>{' '}
              organization as{' '}
              <span className="font-bold" data-testid="role">
                {data?.role}
              </span>{' '}
              on Ensign! Create your account today.
            </p>
          </Card.Body>
        )}
      </Card>
    </>
  );
}
