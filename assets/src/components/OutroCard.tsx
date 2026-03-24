import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';

export const OutroCard = ({
  line,
  badges,
  cta,
  theme,
}: {
  line: string;
  badges: string[];
  cta: string;
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];
  const opacity = fadeIn(frame, 0, 16);
  const shift = slideUp(frame, 0, 20, 16);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: spacing.xl,
        opacity,
        transform: `translateY(${shift}px)`,
      }}
    >
      <div
        style={{
          color: palette.text,
          fontFamily: fonts.display,
          fontSize: 92,
          lineHeight: 1.02,
          fontWeight: 700,
          letterSpacing: '-0.06em',
          maxWidth: 1400,
        }}
      >
        {line}
      </div>
      <div style={{display: 'flex', gap: 16, flexWrap: 'wrap'}}>
        {badges.map((badge, index) => {
          const badgeOpacity = fadeIn(frame, 8 + index * 5, 10);

          return (
            <div
              key={badge}
              style={{
                padding: '14px 20px',
                borderRadius: 999,
                border: `1px solid ${colors.accent}`,
                background: 'rgba(255,255,255,0.04)',
                color: colors.accent,
                fontFamily: fonts.mono,
                fontSize: 18,
                letterSpacing: '0.12em',
                textTransform: 'uppercase',
                opacity: badgeOpacity,
              }}
            >
              {badge}
            </div>
          );
        })}
      </div>
      <div
        style={{
          color: palette.textDim,
          fontFamily: fonts.body,
          fontSize: 30,
          lineHeight: 1.35,
          maxWidth: 1100,
        }}
      >
        {cta}
      </div>
    </div>
  );
};
