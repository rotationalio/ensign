/* eslint-disable jsx-a11y/click-events-have-key-events */
import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React from 'react';

import BCModalVideo from '@/components/common/Modal/BCModalVideos';
import { Image } from '@/components/ui/Image';
import { STARTER_VIDEOS } from '@/features/home/util/utils';

interface StarterVideoProps {
  preview_image: string;
  title: string;
  ytVideoId: string;
  key: string;
}

const StarterVideo = ({ preview_image, title, key }: StarterVideoProps) => {
  return (
    <div key={key} className="flex flex-col gap-2 space-x-2 p-4 hover:font-bold">
      <h2 className="ml-2 flex">{title}</h2>
      <Image src={preview_image} alt={title} className="" />
    </div>
  );
};

const StarterVideos = () => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [videoID, setVideoID] = React.useState('');
  const openVideoHandler = (data: any) => {
    console.log('data', data);
    setIsOpen(true);
    setVideoID(data.ytVideoId);
  };

  const onClose = () => {
    setIsOpen(false);
  };

  return (
    <div className="starter-videos">
      <Heading as="h1" className="pt-10 text-lg font-semibold">
        <Trans>Starter Videos</Trans>
      </Heading>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {STARTER_VIDEOS.map((video, idx) => (
          <button onClick={() => openVideoHandler(video)} key={idx}>
            <StarterVideo
              key={video.title}
              preview_image={video.preview_image}
              title={video.title}
              ytVideoId={video.ytVideoId}
            />
          </button>
        ))}

        <BCModalVideo videoId={videoID} isOpen={isOpen} onClose={onClose} />
      </div>
    </div>
  );
};

export default StarterVideos;
