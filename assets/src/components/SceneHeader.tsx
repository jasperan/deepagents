import {useCurrentFrame, useVideoConfig} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp, softSpring} from '../utils/animation';
import {padVideoNumber} from '../utils/format';

export const SceneHeader = ({
  number,
  kicker,
  title,
  claim,
  theme,
}: {
  number: number;
  kicker: string;
  title: string;
  claim: string;
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const {fps} = useVideoConfig();
  const colors = themeMap[theme];
  const intro = softSpring(frame, fps, 0, 30);
  const opacity = fadeIn(frame, 0, 18);
  const titleShift = slideUp(frame, 0, 28, 18);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: spacing.md,
        opacity,
        transform: `translateY(${titleShift}px) scale(${0.98 + intro * 0.02})`,
      }}
    >
      <div style={{display: 'flex', alignItems: 'center', gap: spacing.sm}}>
        <div
          style={{
            padding: '10px 18px',
            borderRadius: 999,
            border: `1px solid ${colors.accent}`,
            background: 'rgba(255,255,255,0.04)',
            color: colors.accent,
            fontFamily: fonts.mono,
            fontSize: 20,
            letterSpacing: '0.14em',
            textTransform: 'uppercase',
          }}
        >
          {padVideoNumber(number)} / 10
        </div>
        <div
          style={{
            color: palette.textDim,
            fontFamily: fonts.mono,
            fontSize: 20,
            letterSpacing: '0.16em',
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
          fontSize: 88,
          lineHeight: 1.04,
          fontWeight: 700,
          maxWidth: 1180,
          letterSpacing: '-0.05em',
        }}
      >
        {title}
      </div>
      <div
        style={{
          color: palette.textDim,
          fontFamily: fonts.body,
          fontSize: 31,
          lineHeight: 1.35,
          maxWidth: 1060,
        }}
      >
        {claim}
      </div>
    </div>
  );
};
