import type {Panel} from '../data/videos';
import type {ThemeName} from '../design/tokens';
import {CalloutPanel} from './CalloutPanel';
import {CodeWindow} from './CodeWindow';
import {ComparisonPanel} from './ComparisonPanel';
import {DatabaseTable} from './DatabaseTable';
import {FlowPanel} from './FlowPanel';
import {PathRouterDiagram} from './PathRouterDiagram';
import {TerminalWindow} from './TerminalWindow';

export const PanelRenderer = ({panel, theme}: {panel: Panel; theme: ThemeName}) => {
  switch (panel.kind) {
    case 'code':
      return <CodeWindow {...panel} theme={theme} />;
    case 'terminal':
      return <TerminalWindow {...panel} theme={theme} />;
    case 'table':
      return <DatabaseTable {...panel} theme={theme} />;
    case 'callout':
      return <CalloutPanel {...panel} theme={theme} />;
    case 'comparison':
      return <ComparisonPanel {...panel} theme={theme} />;
    case 'flow':
      return <FlowPanel {...panel} theme={theme} />;
    case 'router':
      return <PathRouterDiagram {...panel} theme={theme} />;
    default:
      return null;
  }
};
