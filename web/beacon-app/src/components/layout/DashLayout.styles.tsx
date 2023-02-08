import styled from 'styled-components';

import { SIDEBAR_WIDTH, TOPBAR_HEIGHT } from '@/constants/dash-layout';

export const MainStyle = styled.main`
  height: calc(100vh - ${TOPBAR_HEIGHT});
  margin-left: ${SIDEBAR_WIDTH}px;
  width: calc(100vw - ${SIDEBAR_WIDTH}px);
  margin-top: ${TOPBAR_HEIGHT}px;
`;
