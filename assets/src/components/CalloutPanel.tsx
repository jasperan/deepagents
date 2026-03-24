import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

export const CalloutPanel = ({
  title,
  body,
  bullets,
  stat,
  theme,
}: {
  title: string;
  body: string;
  bullets?: string[];
  stat?: string;
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <GlassPanel
      style={{
        height: '100%',
        padding: 28,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
      }}
    >
      <div style={{display: 'flex', flexDirection: 'column', gap: spacing.lg}}>
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
        <div
          style={{
            color: palette.text,
            fontFamily: fonts.display,
            fontSize: 42,
            lineHeight: 1.12,
            fontWeight: 700,
            letterSpacing: '-0.04em',
          }}
        >
          {body}
        </div>
        {bullets ? (
          <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
            {bullets.map((bullet, index) => {
              const start = 10 + index * 5;
              const opacity = fadeIn(frame, start, 10);
              const shift = slideUp(frame, start, 18, 10);

              return (
                <div
                  key={bullet}
                  style={{
                    display: 'flex',
                    gap: 14,
                    opacity,
                    transform: `translateY(${shift}px)`,
                  }}
                >
                  <div
                    style={{
                      marginTop: 13,
                      width: 8,
                      height: 8,
                      borderRadius: 999,
                      background: colors.accent,
                      flexShrink: 0,
                    }}
                  />
                  <div
                    style={{
                      color: palette.textDim,
                      fontFamily: fonts.body,
                      fontSize: 23,
                      lineHeight: 1.4,
                    }}
                  >
                    {bullet}
                  </div>
                </div>
              );
            })}
          </div>
        ) : null}
      </div>
      {stat ? (
        <div
          style={{
            alignSelf: 'flex-start',
            padding: '12px 18px',
            borderRadius: 999,
            border: `1px solid ${colors.accent}`,
            color: colors.accent,
            fontFamily: fonts.mono,
            fontSize: 18,
            letterSpacing: '0.1em',
            textTransform: 'uppercase',
            background: 'rgba(255,255,255,0.04)',
          }}
        >
          {stat}
        </div>
      ) : null}
    </GlassPanel>
  );
};
