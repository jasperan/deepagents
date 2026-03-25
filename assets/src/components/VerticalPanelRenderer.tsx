import type {ReactNode} from 'react';

import {useCurrentFrame} from 'remotion';

import type {Panel} from '../data/videos';
import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

const accentForTone = (theme: ThemeName, tone?: 'oracle' | 'vector' | 'cloud') => {
  if (tone === 'vector') {
    return themeMap.vector.accent;
  }

  if (tone === 'cloud') {
    return themeMap.cloud.accent;
  }

  return themeMap[theme].accent;
};

const PanelShell = ({
  title,
  theme,
  children,
}: {
  title: string;
  theme: ThemeName;
  children: ReactNode;
}) => {
  const colors = themeMap[theme];

  return (
    <GlassPanel
      style={{
        padding: 26,
        display: 'flex',
        flexDirection: 'column',
        gap: spacing.lg,
        minHeight: 760,
      }}
    >
      <div
        style={{
          color: colors.accent,
          fontFamily: fonts.mono,
          fontSize: 18,
          letterSpacing: '0.14em',
          textTransform: 'uppercase',
        }}
      >
        {title}
      </div>
      {children}
    </GlassPanel>
  );
};

const CodePanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'code'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();
  const highlights = new Set(panel.highlight ?? []);

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 8}}>
        {panel.lines.slice(0, 6).map((line, index) => {
          const opacity = fadeIn(frame, 6 + index * 4, 10);
          const shift = slideUp(frame, 6 + index * 4, 14, 10);
          const isHot = highlights.has(index);

          return (
            <div
              key={`${index}-${line}`}
              style={{
                display: 'grid',
                gridTemplateColumns: '48px 1fr',
                gap: 16,
                padding: '10px 12px',
                borderRadius: 18,
                background: isHot ? 'rgba(78,216,255,0.08)' : 'rgba(255,255,255,0.03)',
                boxShadow: isHot ? `inset 0 0 0 1px ${themeMap[theme].accent}` : 'none',
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              <div
                style={{
                  color: palette.textMute,
                  fontFamily: fonts.mono,
                  fontSize: 18,
                  textAlign: 'right',
                }}
              >
                {index + 1}
              </div>
              <div
                style={{
                  color: palette.text,
                  fontFamily: fonts.mono,
                  fontSize: 24,
                  lineHeight: 1.45,
                  whiteSpace: 'pre-wrap',
                }}
              >
                {line || ' '}
              </div>
            </div>
          );
        })}
      </div>
    </PanelShell>
  );
};

const TerminalPanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'terminal'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 16}}>
        {panel.lines.map((line, index) => {
          const opacity = fadeIn(frame, 6 + index * 6, 10);
          const shift = slideUp(frame, 6 + index * 6, 14, 10);
          const color =
            line.kind === 'command'
              ? colors.accent
              : line.kind === 'success'
                ? palette.green
                : line.kind === 'warning'
                  ? palette.warning
                  : line.kind === 'note'
                    ? palette.textMute
                    : palette.textDim;

          return (
            <div
              key={`${index}-${line.text}`}
              style={{
                display: 'flex',
                gap: 12,
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              <div
                style={{
                  width: 8,
                  height: 8,
                  borderRadius: 999,
                  background: color,
                  marginTop: 12,
                  flexShrink: 0,
                }}
              />
              <div
                style={{
                  color,
                  fontFamily: fonts.mono,
                  fontSize: 24,
                  lineHeight: 1.45,
                  whiteSpace: 'pre-wrap',
                }}
              >
                {line.text}
              </div>
            </div>
          );
        })}
      </div>
    </PanelShell>
  );
};

const CalloutPanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'callout'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div
        style={{
          color: palette.text,
          fontFamily: fonts.display,
          fontSize: 52,
          lineHeight: 1.08,
          letterSpacing: '-0.05em',
          fontWeight: 700,
        }}
      >
        {panel.body}
      </div>
      {panel.bullets ? (
        <div style={{display: 'flex', flexDirection: 'column', gap: 14}}>
          {panel.bullets.map((bullet, index) => {
            const opacity = fadeIn(frame, 8 + index * 4, 10);
            const shift = slideUp(frame, 8 + index * 4, 12, 10);

            return (
              <div
                key={bullet}
                style={{
                  display: 'flex',
                  gap: 12,
                  opacity,
                  transform: `translateY(${shift}px)`,
                }}
              >
                <div
                  style={{
                    width: 8,
                    height: 8,
                    borderRadius: 999,
                    background: themeMap[theme].accent,
                    marginTop: 12,
                    flexShrink: 0,
                  }}
                />
                <div
                  style={{
                    color: palette.textDim,
                    fontFamily: fonts.body,
                    fontSize: 26,
                    lineHeight: 1.35,
                  }}
                >
                  {bullet}
                </div>
              </div>
            );
          })}
        </div>
      ) : null}
      {panel.stat ? (
        <div
          style={{
            alignSelf: 'flex-start',
            padding: '14px 18px',
            borderRadius: 999,
            background: 'rgba(255,255,255,0.05)',
            boxShadow: `inset 0 0 0 1px ${themeMap[theme].accent}`,
            color: themeMap[theme].accent,
            fontFamily: fonts.mono,
            fontSize: 18,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {panel.stat}
        </div>
      ) : null}
    </PanelShell>
  );
};

const TablePanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'table'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 14}}>
        {panel.rows.map((row, index) => {
          const opacity = fadeIn(frame, 8 + index * 4, 10);
          const shift = slideUp(frame, 8 + index * 4, 12, 10);
          const hot = panel.highlightRow === index;

          return (
            <div
              key={`${index}-${row.join('-')}`}
              style={{
                display: 'flex',
                flexDirection: 'column',
                gap: 10,
                padding: '18px 18px',
                borderRadius: 20,
                background: hot ? 'rgba(248,0,0,0.08)' : 'rgba(255,255,255,0.04)',
                boxShadow: hot ? `inset 0 0 0 1px ${themeMap[theme].accent}` : `inset 0 0 0 1px ${palette.panelBorder}`,
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              {row.map((cell, cellIndex) => (
                <div key={`${cellIndex}-${cell}`} style={{display: 'flex', flexDirection: 'column', gap: 4}}>
                  <div
                    style={{
                      color: palette.textMute,
                      fontFamily: fonts.mono,
                      fontSize: 16,
                      letterSpacing: '0.12em',
                      textTransform: 'uppercase',
                    }}
                  >
                    {panel.columns[cellIndex]}
                  </div>
                  <div
                    style={{
                      color: palette.text,
                      fontFamily: cellIndex === 0 ? fonts.mono : fonts.body,
                      fontSize: 24,
                      lineHeight: 1.3,
                    }}
                  >
                    {cell}
                  </div>
                </div>
              ))}
            </div>
          );
        })}
      </div>
    </PanelShell>
  );
};

const ComparisonPanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'comparison'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();

  const renderBucket = (title: string, items: string[], tone: 'left' | 'right') => {
    const accent = tone === 'right' ? themeMap[theme].accent : palette.textMute;

    return (
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 12,
          padding: '18px 18px',
          borderRadius: 20,
          background: 'rgba(255,255,255,0.04)',
          boxShadow: tone === 'right' ? `inset 0 0 0 1px ${themeMap[theme].accent}` : `inset 0 0 0 1px ${palette.panelBorder}`,
        }}
      >
        <div
          style={{
            color: accent,
            fontFamily: fonts.mono,
            fontSize: 17,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {title}
        </div>
        {items.map((item) => (
          <div
            key={item}
            style={{
              color: palette.text,
              fontFamily: fonts.body,
              fontSize: 24,
              lineHeight: 1.32,
            }}
          >
            • {item}
          </div>
        ))}
      </div>
    );
  };

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 16, opacity: fadeIn(frame, 6, 10)}}>
        {renderBucket(panel.leftTitle, panel.leftItems, 'left')}
        {renderBucket(panel.rightTitle, panel.rightItems, 'right')}
      </div>
    </PanelShell>
  );
};

const FlowPanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'flow'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 16}}>
        {panel.steps.map((step, index) => {
          const opacity = fadeIn(frame, 8 + index * 4, 10);
          const shift = slideUp(frame, 8 + index * 4, 12, 10);
          const accent = accentForTone(theme, step.tone);

          return (
            <div
              key={`${step.title}-${index}`}
              style={{
                display: 'flex',
                flexDirection: 'column',
                gap: 10,
                padding: '18px 18px',
                borderRadius: 20,
                background: 'rgba(255,255,255,0.05)',
                boxShadow: `inset 0 0 0 1px ${accent}`,
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              <div
                style={{
                  color: accent,
                  fontFamily: fonts.mono,
                  fontSize: 16,
                  letterSpacing: '0.12em',
                  textTransform: 'uppercase',
                }}
              >
                {step.eyebrow}
              </div>
              <div
                style={{
                  color: palette.text,
                  fontFamily: fonts.display,
                  fontSize: 34,
                  lineHeight: 1.06,
                  letterSpacing: '-0.04em',
                  fontWeight: 700,
                }}
              >
                {step.title}
              </div>
              <div
                style={{
                  color: palette.textDim,
                  fontFamily: fonts.body,
                  fontSize: 24,
                  lineHeight: 1.32,
                }}
              >
                {step.body}
              </div>
            </div>
          );
        })}
      </div>
    </PanelShell>
  );
};

const RouterPanelView = ({panel, theme}: {panel: Extract<Panel, {kind: 'router'}>; theme: ThemeName}) => {
  const frame = useCurrentFrame();

  return (
    <PanelShell title={panel.title} theme={theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 14}}>
        {panel.routes.map((route, index) => {
          const opacity = fadeIn(frame, 8 + index * 4, 10);
          const shift = slideUp(frame, 8 + index * 4, 12, 10);
          const accent = accentForTone(theme, route.tone);

          return (
            <div
              key={`${route.path}-${route.target}`}
              style={{
                display: 'flex',
                flexDirection: 'column',
                gap: 10,
                padding: '16px 18px',
                borderRadius: 20,
                background: 'rgba(255,255,255,0.05)',
                boxShadow: `inset 0 0 0 1px ${accent}`,
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              <div
                style={{
                  color: accent,
                  fontFamily: fonts.mono,
                  fontSize: 17,
                  letterSpacing: '0.12em',
                }}
              >
                {route.path}
              </div>
              <div
                style={{
                  color: palette.text,
                  fontFamily: fonts.display,
                  fontSize: 30,
                  lineHeight: 1.08,
                  fontWeight: 700,
                }}
              >
                ↓ {route.target}
              </div>
            </div>
          );
        })}
        {panel.footer ? (
          <div
            style={{
              color: palette.textDim,
              fontFamily: fonts.body,
              fontSize: 24,
              lineHeight: 1.32,
            }}
          >
            {panel.footer}
          </div>
        ) : null}
      </div>
    </PanelShell>
  );
};

export const VerticalPanelRenderer = ({panel, theme}: {panel: Panel; theme: ThemeName}) => {
  switch (panel.kind) {
    case 'code':
      return <CodePanelView panel={panel} theme={theme} />;
    case 'terminal':
      return <TerminalPanelView panel={panel} theme={theme} />;
    case 'callout':
      return <CalloutPanelView panel={panel} theme={theme} />;
    case 'table':
      return <TablePanelView panel={panel} theme={theme} />;
    case 'comparison':
      return <ComparisonPanelView panel={panel} theme={theme} />;
    case 'flow':
      return <FlowPanelView panel={panel} theme={theme} />;
    case 'router':
      return <RouterPanelView panel={panel} theme={theme} />;
    default:
      return null;
  }
};
