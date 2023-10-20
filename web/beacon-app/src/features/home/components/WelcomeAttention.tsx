import { Trans } from '@lingui/macro';
import React, { useState } from 'react';

import BCModalVideo from '@/components/common/Modal/BCModalVideos';
// todo: ensure to use a better image , this one seems creepy
import { Image } from '@/components/ui/Image';

const welcomevideo = {
  preview_image: 'https://i.ytimg.com/vi/wurObU34Kes/maxres1.jpg',
  videoID: 'wurObU34Kes',
};

const WelcomeAttention = () => {
  const [isOpen, setIsOpen] = useState(false);

  const videoModalHandler = () => {
    setIsOpen(true);
  };

  return (
    <>
      <div
        className="px-auto mb-8 mt-4 flex flex-col justify-between gap-4 space-x-10 rounded-md border border-black/30 bg-[#F7F9FB]  p-4 text-justify xl:flex-row"
        data-cy="ensign-welcome"
      >
        <div className="flex flex-col space-y-10 sm:w-4/5">
          <p className="text-md">
            <Trans>
              <span className="font-bold">Welcome to Ensign</span>, your all-in-one platform for
              real-time data management. Ensign is a flexible database meets streaming engine for
              data teams to build and deploy real-time models, data products, and services.
            </Trans>
          </p>
          <p>
            <Trans>Ready to dive in? Learn how to use Ensign or start your first project.</Trans>
          </p>
        </div>

        <div className="flex" data-cy="welcome-video">
          <button onClick={videoModalHandler} data-cy="welcome-vid-bttn">
            <Image
              src={welcomevideo.preview_image}
              alt="welcome video"
              className="float-right h-[77px] w-[139px]"
            />
          </button>
        </div>
      </div>
      <BCModalVideo
        videoId={welcomevideo.videoID}
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
      />
    </>
  );
};

export default WelcomeAttention;
