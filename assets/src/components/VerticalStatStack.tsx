import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';

export const VerticalStatStack = ({
  chips,
  theme,
}: {
  chips: string[];
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
      {chips.map((chip, index) => {
        const opacity = fadeIn(frame, 10 + index * 4, 10);
        const shift = slideUp(frame, 10 + index * 4, 14, 10);

        return (
          <div
            key={chip}
            style={{
              padding: '14px 18px',
              borderRadius: 20,
              background: 'rgba(255,255,255,0.05)',
              boxShadow: `inset 0 0 0 1px ${colors.accent}`,
              color: palette.text,
              fontFamily: fonts.mono,
              fontSize: 18,
              letterSpacing: '0.12em',
              textTransform: 'uppercase',
              opacity,
              transform: `translateY(${shift}px)`,
            }}
          >
            {chip}
          </div>
        );
      })}
    </div>
  );
};
