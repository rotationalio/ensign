import AppLayout from '@/components/layout/AppLayout';

import QuickStart from '../components/QuickStart';
import WelcomeAttention from '../components/WelcomeAttention';

export default function Home() {
  return (
    <AppLayout>
      <WelcomeAttention />
      <QuickStart />
    </AppLayout>
  );
}
