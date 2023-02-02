/* Creates a URL with an organization domain */
/* Ex. Rotational Labs -> ensign.rotational.io/rotational-labs */

export default function slugify(domain: string) {
    return "ensign.rotational.io/" + 
    domain
        .normalize('NFKD') /* Splits the base character and its accent */
        .replace( /[\u0300-\u036f]/g, '' ) /* Deletes all the accents */
        .toLowerCase()
        .trim()
        .replace(/\s+/g, '-')
        .replace(/[^\w\-]+/g, '')
        .replace(/\-\-+/g, '-');
}