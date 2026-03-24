import {useCurrentFrame} from 'remotion';

import type {TerminalLine} from '../data/videos';
import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

const lineColor = (kind: TerminalLine['kind'], accent: string) => {
  switch (kind) {
    case 'command':
      return accent;
    case 'success':
      return palette.green;
    case 'warning':
      return palette.warning;
    case 'note':
      return palette.textMute;
    default:
      return palette.textDim;
  }
};

export const TerminalWindow = ({
  title,
  lines,
  theme,
}: {
  title: string;
  lines: TerminalLine[];
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <GlassPanel
      style={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        background: 'rgba(5, 9, 16, 0.84)',
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          padding: '18px 24px',
          borderBottom: `1px solid ${palette.panelBorder}`,
        }}
      >
        <div style={{display: 'flex', alignItems: 'center', gap: 8}}>
          <div style={{width: 10, height: 10, borderRadius: 999, background: colors.accent}} />
          <div style={{width: 10, height: 10, borderRadius: 999, background: palette.blue}} />
          <div style={{width: 10, height: 10, borderRadius: 999, background: palette.violet}} />
        </div>
        <div
          style={{
            color: palette.textDim,
            fontFamily: fonts.mono,
            fontSize: 17,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {title}
        </div>
        <div
          style={{
            color: palette.textMute,
            fontFamily: fonts.mono,
            fontSize: 16,
          }}
        >
          live
        </div>
      </div>
      <div
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          gap: spacing.sm,
          padding: '26px 28px',
        }}
      >
        {lines.map((line, index) => {
          const start = 6 + index * 8;
          const opacity = fadeIn(frame, start, 10);
          const shift = slideUp(frame, start, 18, 10);

          return (
            <div
              key={`${index}-${line.text}`}
              style={{
                display: 'flex',
                alignItems: 'flex-start',
                gap: 14,
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              <div
                style={{
                  marginTop: 7,
                  width: 8,
                  height: 8,
                  borderRadius: 999,
                  background: line.kind === 'command' ? colors.accent : 'rgba(255,255,255,0.22)',
                  flexShrink: 0,
                }}
              />
              <div
                style={{
                  color: lineColor(line.kind, colors.accent),
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
    </GlassPanel>
  );
};
