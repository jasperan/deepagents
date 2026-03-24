import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

export const ComparisonPanel = ({
  title,
  leftTitle,
  rightTitle,
  leftItems,
  rightItems,
  theme,
}: {
  title: string;
  leftTitle: string;
  rightTitle: string;
  leftItems: string[];
  rightItems: string[];
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  const renderColumn = (heading: string, items: string[], tone: 'left' | 'right') => {
    const accent = tone === 'left' ? 'rgba(255,255,255,0.16)' : colors.accent;

    return (
      <div
        style={{
          flex: 1,
          display: 'flex',
          flexDirection: 'column',
          gap: spacing.md,
          padding: 24,
          borderRadius: 24,
          background: tone === 'left' ? 'rgba(255,255,255,0.03)' : 'rgba(255,255,255,0.05)',
          boxShadow: tone === 'right' ? `inset 0 0 0 1px ${colors.accent}` : 'none',
        }}
      >
        <div
          style={{
            color: tone === 'right' ? colors.accent : palette.textDim,
            fontFamily: fonts.mono,
            fontSize: 18,
            letterSpacing: '0.12em',
            textTransform: 'uppercase',
          }}
        >
          {heading}
        </div>
        <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
          {items.map((item, index) => {
            const start = 10 + index * 4;
            const opacity = fadeIn(frame, start, 10);
            const shift = slideUp(frame, start, 16, 10);

            return (
              <div
                key={item}
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
                    marginTop: 12,
                    background: accent,
                    flexShrink: 0,
                  }}
                />
                <div
                  style={{
                    color: palette.text,
                    fontFamily: fonts.body,
                    fontSize: 23,
                    lineHeight: 1.35,
                  }}
                >
                  {item}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  return (
    <GlassPanel
      style={{
        height: '100%',
        padding: 28,
        display: 'flex',
        flexDirection: 'column',
        gap: spacing.lg,
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
      <div style={{display: 'flex', gap: spacing.lg, flex: 1}}>
        {renderColumn(leftTitle, leftItems, 'left')}
        {renderColumn(rightTitle, rightItems, 'right')}
      </div>
    </GlassPanel>
  );
};
