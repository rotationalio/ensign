import * as RadixAvatar from '@radix-ui/react-avatar';
import styled from 'styled-components';

export const StyledAvatar = styled(RadixAvatar.Root)`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  vertical-align: middle;
  overflow: hidden;
  user-select: none;
  width: 45px;
  height: 45px;
  border-radius: 100%;
  background-color: gray;
`;

export const StyledAvatarImage = styled(RadixAvatar.Image)`
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: inherit;
`;

export const StyledAvatarFallback = styled(RadixAvatar.Fallback)`
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--colors-secondary-900);
  color: #fff;
  font-size: 15px;
  line-height: 1;
  font-weight: 500;
`;
