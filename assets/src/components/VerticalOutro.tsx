import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';

export const VerticalOutro = ({
  line,
  cta,
  theme,
}: {
  line: string;
  cta: string;
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];
  const opacity = fadeIn(frame, 0, 14);
  const shift = slideUp(frame, 0, 18, 14);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: 22,
        opacity,
        transform: `translateY(${shift}px)`,
      }}
    >
      <div
        style={{
          color: palette.text,
          fontFamily: fonts.display,
          fontSize: 82,
          lineHeight: 0.98,
          letterSpacing: '-0.06em',
          fontWeight: 700,
        }}
      >
        {line}
      </div>
      <div
        style={{
          alignSelf: 'flex-start',
          padding: '16px 22px',
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
        {cta}
      </div>
    </div>
  );
};
