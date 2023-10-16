// capitalize a string
export const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const getInitials = (name: string) => {
  const nameArray = name.split(' ');
  const initials =
    nameArray.length >= 2
      ? nameArray[0].charAt(0) + nameArray[1].charAt(0)
      : nameArray[0].charAt(0) + nameArray[0].charAt(1);

  // console.log('[] initials', initials);
  return initials.toUpperCase();
};
