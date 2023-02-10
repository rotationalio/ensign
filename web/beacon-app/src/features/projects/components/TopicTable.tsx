import { Table, Toast } from '@rotational/beacon-core';
import { useState } from 'react';

import { TableHeading } from '@/components/common/TableHeader';
import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';
import { Topic } from '@/features/topics/types/topicService';

interface TopicTableProps {
  projectID: string;
}

export const TopicTable = ({ projectID }: TopicTableProps) => {
  const [items, setItems] = useState<Topic[]>([]);
  const [, setIsOpen] = useState(false);
  const handleClose = () => setIsOpen(false);

  const { topics, isFetchingTopics, wasTopicsFetched, hasTopicsFailed, error } =
    useFetchTopics(projectID);

  if (isFetchingTopics) {
    // TODO: add loading state
    return <div>Loading...</div>;
  }

  if (error) {
    return (
      <Toast
        isOpen={hasTopicsFailed}
        onClose={handleClose}
        variant="danger"
        title="Sorry we are having trouble fetching your topics, please try again later."
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  if (wasTopicsFetched && topics) {
    const newItems = topics.topics || [];
    setItems(newItems);
  }

  return (
    <div>
      <TableHeading>Topics</TableHeading>
      <Table
        columns={[
          { Header: 'Topics ID', accessor: 'id' },
          { Header: 'Name', accessor: 'name' },
        ]}
        data={items}
      />
    </div>
  );
};

export default TopicTable;
