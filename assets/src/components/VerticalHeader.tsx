import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';

export const VerticalHeader = ({
  number,
  kicker,
  hook,
  subhook,
  theme,
}: {
  number: number;
  kicker: string;
  hook: string;
  subhook: string;
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];
  const opacity = fadeIn(frame, 0, 14);
  const shift = slideUp(frame, 0, 20, 14);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: spacing.md,
        opacity,
        transform: `translateY(${shift}px)`,
      }}
    >
      <div style={{display: 'flex', gap: 12, flexWrap: 'wrap'}}>
        <div
          style={{
            padding: '10px 16px',
            borderRadius: 999,
            background: 'rgba(255,255,255,0.06)',
            boxShadow: `inset 0 0 0 1px ${colors.accent}`,
            color: colors.accent,
            fontFamily: fonts.mono,
            fontSize: 18,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {String(number).padStart(2, '0')} / 10
        </div>
        <div
          style={{
            padding: '10px 16px',
            borderRadius: 999,
            background: 'rgba(255,255,255,0.04)',
            color: palette.textDim,
            fontFamily: fonts.mono,
            fontSize: 18,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {kicker}
        </div>
      </div>
      <div
        style={{
          color: palette.text,
          fontFamily: fonts.display,
          fontSize: 92,
          lineHeight: 0.98,
          letterSpacing: '-0.06em',
          fontWeight: 700,
        }}
      >
        {hook}
      </div>
      <div
        style={{
          color: palette.textDim,
          fontFamily: fonts.body,
          fontSize: 28,
          lineHeight: 1.3,
          maxWidth: 880,
        }}
      >
        {subhook}
      </div>
    </div>
  );
};
