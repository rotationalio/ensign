export const isDate = (value: any) => {
  return value instanceof Date && !isNaN(value.getTime());
};

export const formatDate = (date: Date, format?: string, type?: string) => {
  if (!isDate(date)) {
    return 'N/A';
  } else {
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
  }
};
