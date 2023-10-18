import 'node_modules/react-modal-video/scss/modal-video.scss';

import React, { Fragment } from 'react';
import ModalVideo from 'react-modal-video';
interface Props {
  key?: string;
  videoId: string;
  channel?: string; // youtube, vimeo, etc. but let use the default youtube
  isOpen: boolean;
  onClose: () => void;
}

const BCModalVideo = ({ key, videoId, isOpen, onClose }: Props) => {
  const k = Math.floor(Math.random() * 1000);
  return (
    <Fragment key={key || k}>
      <ModalVideo isOpen={isOpen} videoId={videoId} onClose={onClose} allowFullScreen />
    </Fragment>
  );
};

export default BCModalVideo;
