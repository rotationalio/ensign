import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';

function UpdateAlert() {
  return (
    <>
      <p>
        <Trans>
          A new version of Ensign is available. Click <span className="font-bold">"Update"</span> to
          use the latest version.
        </Trans>
      </p>
      <div>
        <Button
          size="small"
          onClick={() => window.location.reload()}
          data-testid="update-alert-btn"
        >
          <Trans>Update</Trans>
        </Button>
      </div>
    </>
  );
}

export default UpdateAlert;
