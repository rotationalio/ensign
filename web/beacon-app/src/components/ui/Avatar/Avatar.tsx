import { twMerge } from 'tailwind-merge';

import getInitials from '@/utils/getInitials';

import { StyledAvatar, StyledAvatarFallback, StyledAvatarImage } from './Avatar.styles';
import { AvatarProps } from './Avatar.type';

const Avatar: React.FC<AvatarProps> = ({ className, fallbackProps, ...imageProps }) => {
  return (
    <StyledAvatar>
      <StyledAvatarImage {...imageProps} className={twMerge(className)} />
      <StyledAvatarFallback
        delayMs={600}
        {...fallbackProps}
        className={twMerge('capitalize', fallbackProps?.className)}
      >
        {getInitials(imageProps.alt)}
      </StyledAvatarFallback>
    </StyledAvatar>
  );
};

export default Avatar;
