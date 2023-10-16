import * as React from 'react';
/**
 * React Hook to manage a menu
 *
 * It provides the logic and will be used with react context
 * to propagate its return value to all children
 *  for now we will just handle the open/close/isOpen state
 * and the anchor element
 * @param {MenuProps} props
 * @returns {MenuState}
 * @example
 * const menu = useMenu({ id: 'my-menu' })
 * return (
 * <MenuContext.Provider value={menu}>
 *  <button onClick={menu.open}>Open Menu</button>
 * <Menu id={menu.id} anchorEl={menu.anchorEl} open={menu.open}>
 * <MenuItem>Menu Item</MenuItem>
 * </Menu>
 * </MenuContext.Provider>
 * )
 */
export interface MenuProps {
  id: string;
}
export interface MenuState {
  id: string;
  anchorEl: null | HTMLElement;
  open: (event: React.MouseEvent<HTMLElement>) => void;
  close: () => void;
  isOpen: boolean;
}
const useMenu = (props: MenuProps): MenuState => {
  const { id } = props;
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const close = () => {
    setAnchorEl(null);
  };

  const isOpen = Boolean(anchorEl);

  return { id, anchorEl, open, close, isOpen };
};

export default useMenu;
