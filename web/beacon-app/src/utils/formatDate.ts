export const isDate = (value: any) => value instanceof Date;

export const formatDate = (date: Date, format?: string, type?: string) => {
  if (!format) {
    return Intl.DateTimeFormat(undefined, {
      localeMatcher: 'best fit',
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    }).format(date);
  }
  if (type === 'time') {
    return Intl.DateTimeFormat(undefined, {
      localeMatcher: 'best fit',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  }
  if (type === 'date') {
    return Intl.DateTimeFormat(undefined, {
      localeMatcher: 'best fit',
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    }).format(date);
  }
};
