import AppLayout from '@/components/layout/AppLayout';
import { useOrgStore } from '@/store';
import { setCookie } from '@/utils/cookies';

import QuickStart from '../components/QuickStart';
import StarterVideos from '../components/StarterVideos/StarterVideo';
import WelcomeAttention from '../components/WelcomeAttention';

export default function Home() {
  const store = useOrgStore((state) => state) as any;
  const isAuthenticated = store.isAuthenticated;
  setCookie('authenticatedUser', isAuthenticated as string);

  return (
    <AppLayout>
      <WelcomeAttention />
      <QuickStart />
      <StarterVideos />
    </AppLayout>
  );
}
