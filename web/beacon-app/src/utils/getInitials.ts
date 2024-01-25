export const getInitials = (name: string) => {
  const nameArray = name?.split(' ') || [''];
  const initials =
    nameArray.length >= 2
      ? nameArray[0].charAt(0) + nameArray[1].charAt(0)
      : nameArray[0].charAt(0) + nameArray[0].charAt(1);

  return initials.toUpperCase();
};
