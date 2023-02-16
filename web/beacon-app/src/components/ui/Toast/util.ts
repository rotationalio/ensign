export const getToastBgVariantStyle = (variant: string) => {
  switch (variant) {
    case 'success':
      return 'bg-green-success';
    case 'danger':
      return 'bg-danger';
    case 'warning':
      return 'bg-warning';
    case 'info':
      return 'bg-primary';
    case 'primary':
      return 'bg-primary';
    case 'secondary':
      return 'bg-secondary-900';
    default:
      return 'bg-white';
  }
};

export const getToastColorVariantStyle = (variant: string) => {
  switch (variant) {
    case 'success':
    case 'danger':
    case 'warning':
    case 'info':
    case 'primary':
    case 'secondary':
      return 'text-white';
    default:
      return 'text-gray-900';
  }
};
