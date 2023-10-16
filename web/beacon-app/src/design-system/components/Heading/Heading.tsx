export type HeadingLevel = 1 | 2 | 3 | 4 | 5 | 6;

type TextProps<C extends React.ElementType> = {
  level?: HeadingLevel;
  as?: C;
};

type HeadingProps<C extends React.ElementType> = React.PropsWithChildren<TextProps<C>> &
  React.ComponentPropsWithoutRef<C>;

const Heading = <C extends React.ElementType = 'h3'>(props: HeadingProps<C>) => {
  const { as, children, level = 1, ...rest } = props;
  const As = as || `h${level}`;

  return <As {...rest}>{children}</As>;
};

export default Heading;
