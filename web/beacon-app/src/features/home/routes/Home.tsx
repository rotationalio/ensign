import AppLayout from '@/components/layout/AppLayout';

import QuickStart from '../components/QuickStart';
import StarterVideos from '../components/StarterVideos/StarterVideo';
import WelcomeAttention from '../components/WelcomeAttention';

export default function Home() {
  return (
    <AppLayout>
      <WelcomeAttention />
      <QuickStart />

      <StarterVideos />
    </AppLayout>
  );
}
