import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
  // define default options for all queries
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      //cacheTime: 5 * 60 * 1000,
      retry: 0,
      //retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
      refetchInterval: false,
      refetchIntervalInBackground: false,
      refetchOnMount: true,
      refetchOnReconnect: true,
      refetchOnWindowFocus: false,
    },
  },
});

const queryCache = new QueryCache({
  onError: (error) => {
    console.log(error);
  },
  onSuccess: (data) => {
    console.log(data);
  },
});
const mutationCache = new MutationCache({
  onError: (error) => {
    console.log(error);
  },
  onSuccess: (data) => {
    console.log(data);
  },
});

export { mutationCache, queryCache, queryClient, QueryClientProvider };
