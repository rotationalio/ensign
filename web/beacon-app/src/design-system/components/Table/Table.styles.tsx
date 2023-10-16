import { rgba } from 'polished';
import styled from 'styled-components';

const StyledBase = styled.div`
  display: flex;
  flex-direction: column;
`;

const StyledTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  tbody:before {
    content: '-';
    display: block;
    line-height: 0.6em;
    color: transparent;
  }
`;

const StyledTh = styled.th`
  text-align: ${(props) => (props.align ? props.align : 'left')};
  font-size: 18px;
  font-weight: 700;
  opacity: 0.65;
`;

const StyledTd = styled.td`
  text-align: ${(props) => (props.align ? props.align : 'left')};
`;

const StyledTr = styled.tr`
  border-bottom: 2px solid ${rgba(150, 150, 150, 0.5)};
`;

export { StyledBase, StyledTable, StyledTd, StyledTh, StyledTr };
