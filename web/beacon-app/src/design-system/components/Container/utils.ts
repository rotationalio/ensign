import { ContainerVariant } from './Container.types';
export const setVariantStyle = (variant: ContainerVariant) => {
  switch (variant) {
    case 'default':
      return 'max-w-[1662px] mx-auto px-4 sm:px-6 lg:px-8';
    case 'dash':
      return 'max-w-[1440] mx-auto px-4 sm:px-6 lg:px-8';
    case 'base':
      return 'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8';
    default:
      return 'max-w-[1440px] mx-auto px-4 sm:px-6 lg:px-8';
  }
};

export enum CONTAINER_VARIANT {
  DEFAULT = 'default',
  DASH = 'dash',
  BASE = 'base',
}
