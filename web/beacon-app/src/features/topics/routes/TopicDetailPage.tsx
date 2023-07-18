import { Heading } from '@rotational/beacon-core';
import invariant from 'invariant';
import { useParams } from 'react-router-dom';

import AppLayout from '@/components/layout/AppLayout';

const TopicDetailPage = () => {
  const param = useParams<{ id: string }>();
  const { id: topicID } = param;

  invariant(topicID, 'topic id is required');
  return (
    <AppLayout>
      <Heading as="h1" className="flex items-center text-lg font-semibold">
        Topic Name
      </Heading>
      {/*      <Suspense
        fallback={
          <div>
            <Loader />
          </div>
        }
      ></Suspense> */}
    </AppLayout>
  );
};

export default TopicDetailPage;
