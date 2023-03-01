export const isDate = (value: any) => value instanceof Date;

export const formatDate = (date: Date, format?: string) => {
  if (!format) {
    return Intl.DateTimeFormat(undefined, {
      localeMatcher: 'best fit',
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  }
};
