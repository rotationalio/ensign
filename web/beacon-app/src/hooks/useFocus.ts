import { useState } from 'react';

type UseFocus = [
  boolean,
  {
    onFocus: () => void;
    onBlur: () => void;
  }
];

export default function useFocus(): UseFocus {
  const [focused, setFocused] = useState(false);

  const onFocus = () => setFocused(true);
  const onBlur = () => setFocused(false);

  return [focused, { onFocus, onBlur }];
}
