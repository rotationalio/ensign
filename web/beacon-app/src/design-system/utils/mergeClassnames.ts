import { extendTailwindMerge } from 'tailwind-merge';

const mergeClassnames = extendTailwindMerge({
  cacheSize: 0,
});

export default mergeClassnames;
