import styled from 'styled-components';

import { SIDEBAR_WIDTH, TOPBAR_HEIGHT } from '@/constants/dash-layout';

export const Header = styled.header`
  min-height: ${TOPBAR_HEIGHT}px;
  margin-left: ${SIDEBAR_WIDTH}px;
  border-bottom-width: 1px;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  display: flex;
  align-items: center;
`;
