/**
 * Get initials of a given name (e.g: Calumn Scott -> CS)
 * @param name
 * @returns string
 */
export default function getInitials(name = '') {
  return name
    .match(/(^\S\S?|\b\S)?/g)
    ?.join('')
    .match(/(^\S|\S$)?/g)
    ?.join('')
    .toUpperCase();
}
