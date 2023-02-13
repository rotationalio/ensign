import { Button } from '@rotational/beacon-core';
import { QueryErrorResetBoundary } from '@tanstack/react-query';
import { ErrorBoundary } from 'react-error-boundary';

// This is a wrapper around the react-query QueryErrorResetBoundary and the Sentry ErrorBoundary
// that allows us to send errors to Sentry and reset the query cache when an error occurs.

const RTKErrorBoundary: React.FC<any> = ({ children }) => (
  <QueryErrorResetBoundary>
    {({ reset }) => (
      <ErrorBoundary
        onReset={reset}
        onError={() => {
          // send to sentry
        }}
        fallbackRender={({ resetErrorBoundary }) => (
          <div className="item-center text-md justify-center">
            There was an error!
            <Button variant="primary" size="medium" onClick={() => resetErrorBoundary()}>
              Try again
            </Button>
          </div>
        )}
      >
        {children}
      </ErrorBoundary>
    )}
  </QueryErrorResetBoundary>
);

export default RTKErrorBoundary;
