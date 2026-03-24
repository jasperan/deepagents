import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideRight} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

export const CodeWindow = ({
  title,
  language,
  lines,
  highlight,
  theme,
}: {
  title: string;
  language: string;
  lines: string[];
  highlight?: number[];
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];
  const highlighted = new Set(highlight ?? []);

  return (
    <GlassPanel
      style={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          padding: '18px 24px',
          borderBottom: `1px solid ${palette.panelBorder}`,
          background: 'rgba(255,255,255,0.03)',
        }}
      >
        <div style={{display: 'flex', alignItems: 'center', gap: 10}}>
          <div style={{width: 11, height: 11, borderRadius: 999, background: '#FF5F57'}} />
          <div style={{width: 11, height: 11, borderRadius: 999, background: '#FEBC2E'}} />
          <div style={{width: 11, height: 11, borderRadius: 999, background: '#28C840'}} />
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
            color: colors.accent,
            fontFamily: fonts.mono,
            fontSize: 16,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {language}
        </div>
      </div>
      <div
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          gap: 8,
          padding: 24,
          background: 'linear-gradient(180deg, rgba(255,255,255,0.02) 0%, rgba(255,255,255,0.01) 100%)',
        }}
      >
        {lines.map((line, index) => {
          const localStart = 6 + index * 4;
          const opacity = fadeIn(frame, localStart, 10);
          const shift = slideRight(frame, localStart, 18, 10);
          const isHot = highlighted.has(index);

          return (
            <div
              key={`${index}-${line}`}
              style={{
                display: 'grid',
                gridTemplateColumns: '52px 1fr',
                alignItems: 'start',
                gap: spacing.md,
                opacity,
                transform: `translateX(${shift}px)`,
                padding: '6px 10px',
                borderRadius: 16,
                background: isHot ? 'rgba(78, 216, 255, 0.08)' : 'transparent',
                boxShadow: isHot ? `inset 0 0 0 1px ${colors.accent}` : 'none',
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
                  color: isHot ? palette.text : palette.textDim,
                  fontFamily: fonts.mono,
                  fontSize: 24,
                  lineHeight: 1.5,
                  whiteSpace: 'pre-wrap',
                }}
              >
                {line || ' '}
              </div>
            </div>
          );
        })}
      </div>
    </GlassPanel>
  );
};
