export const getImgWrapperStyle = (imgPosition: string) => {
  switch (imgPosition) {
    case 'top':
      return 'flex items-center';
    case 'left' || 'right':
      return 'grid grid-cols-2';
    case 'top-middle':
      return 'flex items-center';
    default:
      return 'column';
  }
};
