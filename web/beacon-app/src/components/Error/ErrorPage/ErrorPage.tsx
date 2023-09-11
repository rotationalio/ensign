import { t, Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { isRouteErrorResponse, useRouteError } from 'react-router-dom';

import NotFoundOtters from '@/assets/images/not-found-otters.svg';
import { Link } from '@/components/ui/Link';

type ErrorPageProps = {
  errorMessage?: string;
  errorCause?: string;
  errorTitle?: string;
};

export const render404 = () => {
  return (
    <section className="mx-auto my-20  flex max-w-4xl place-items-center items-center justify-center rounded-lg border border-solid border-primary-800 text-2xl">
      <div className="my-10 mx-auto max-w-xl">
        <h1 className="mt-4 text-2xl font-bold text-gray-800">
          <Trans>Sorry, we can’t find that page. (404)</Trans>
        </h1>
        <Trans>
          <p className="mt-4">
            Return to
            <Link href="/"> rotational.app </Link>
            or please contact us at support@rotational.io for assistance.
          </p>
        </Trans>
        <img src={NotFoundOtters} alt="not found otters" className="mx-auto mt-20" />
      </div>
    </section>
  );
};

export default function ErrorPage({ errorMessage, errorCause, errorTitle }: ErrorPageProps) {
  const { error } = useRouteError() as { error: Error };
  const Error = useRouteError() as { error: Error };
  if (isRouteErrorResponse(Error) && Error?.status === 404) {
    return render404();
  }

  return (
    <section
      className="mx-auto my-20  flex max-w-4xl place-items-center items-center justify-center rounded-lg border border-solid border-primary-800 text-2xl"
      data-testid="error-page"
    >
      <div className="my-10 mx-auto max-w-xl">
        <h1 className="text-2xl font-bold text-gray-800">
          {' '}
          {errorTitle || t`Sorry, we’re having trouble loading this page.`}
        </h1>

        <p className="text-xl text-gray-600">
          <pre>{(error?.cause as any) || errorCause}</pre>
        </p>

        <p className="text-xl text-gray-600">{error?.message || errorMessage}</p>

        <Button className="mt-4" onClick={() => window.location.reload()}>
          Reload
        </Button>
      </div>
    </section>
  );
}
