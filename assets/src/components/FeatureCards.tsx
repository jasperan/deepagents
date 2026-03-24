import {useCurrentFrame} from 'remotion';

import type {BenefitCard} from '../data/videos';
import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

export const FeatureCards = ({
  cards,
  theme,
}: {
  cards: BenefitCard[];
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <div style={{display: 'grid', gridTemplateColumns: 'repeat(3, minmax(0, 1fr))', gap: spacing.lg}}>
      {cards.map((card, index) => {
        const start = 6 + index * 8;
        const opacity = fadeIn(frame, start, 12);
        const shift = slideUp(frame, start, 24, 12);

        return (
          <GlassPanel
            key={card.title}
            style={{
              minHeight: 280,
              padding: 28,
              display: 'flex',
              flexDirection: 'column',
              gap: spacing.md,
              opacity,
              transform: `translateY(${shift}px)`,
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
              {card.eyebrow}
            </div>
            <div
              style={{
                color: palette.text,
                fontFamily: fonts.display,
                fontSize: 38,
                lineHeight: 1.08,
                fontWeight: 700,
                letterSpacing: '-0.04em',
              }}
            >
              {card.title}
            </div>
            <div
              style={{
                color: palette.textDim,
                fontFamily: fonts.body,
                fontSize: 24,
                lineHeight: 1.45,
              }}
            >
              {card.body}
            </div>
          </GlassPanel>
        );
      })}
    </div>
  );
};
