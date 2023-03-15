import { Loader, Table, Toast } from '@rotational/beacon-core';

import { TableHeading } from '@/components/common/TableHeader';
import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';
import { Topic } from '@/features/topics/types/topicService';

export const TopicTable = () => {
  // const [, setIsOpen] = useState(false);
  // const handleClose = () => setIsOpen(false);

  const { getTopics, topics, isFetchingTopics, hasTopicsFailed, error } = useFetchTopics();

  if (!topics) {
    getTopics();
  }

  if (isFetchingTopics) {
    // TODO: add loading state
    return <Loader />;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasTopicsFailed}
        duration={3000}
        variant="danger"
        title="Sorry we are having trouble fetching your topics, please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  return (
    <div className="my-5">
      <TableHeading>Topics</TableHeading>
      <div className="overflow-hidden text-sm">
        <Table
          columns={[
            { Header: 'Topics ID', accessor: 'id' },
            { Header: 'Name', accessor: 'name' },
          ]}
          data={(topics.topics as Topic[]) || []}
        />
      </div>
    </div>
  );
};

export default TopicTable;
