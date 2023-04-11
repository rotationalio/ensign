import { useRouteError } from 'react-router-dom';

// error page with reload button

export default function ErrorPage() {
  const { error } = useRouteError() as { error: Error };

  return (
    <div className="flex h-full w-full flex-col items-center justify-center">
      <h1 className="text-2xl font-bold text-gray-800">Something went wrong</h1>
      <p className="text-xl text-gray-600">{error.message}</p>
      <button onClick={() => window.location.reload()}>Reload</button>
    </div>
  );
}
