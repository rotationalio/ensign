import { AvatarFallbackProps, AvatarImageProps } from '@radix-ui/react-avatar';
import { SetRequired } from 'type-fest';

type FallbackProps = Omit<AvatarFallbackProps, 'children'>;

export type AvatarProps = {
  fallbackProps?: FallbackProps;
} & SetRequired<AvatarImageProps, 'alt'>;
