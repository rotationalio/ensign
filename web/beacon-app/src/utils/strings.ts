// capitalize a string
export const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const getInitials = (name: string) => {
  const nameArray = name.split(' ');
  const initials = nameArray[0].charAt(0) + nameArray[1].charAt(0);
  return initials.toUpperCase();
};
