import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

export type FlowStep = {
  eyebrow: string;
  title: string;
  body: string;
  tone?: 'oracle' | 'vector' | 'cloud';
};

const stepColor = (theme: ThemeName, tone?: FlowStep['tone']) => {
  if (tone === 'vector') {
    return themeMap.vector.accent;
  }

  if (tone === 'cloud') {
    return themeMap.cloud.accent;
  }

  return themeMap[theme].accent;
};

export const FlowPanel = ({
  title,
  steps,
  theme,
}: {
  title: string;
  steps: FlowStep[];
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();

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
          color: themeMap[theme].accent,
          fontFamily: fonts.mono,
          fontSize: 18,
          letterSpacing: '0.14em',
          textTransform: 'uppercase',
        }}
      >
        {title}
      </div>
      <div style={{display: 'flex', gap: 18, flex: 1, alignItems: 'stretch'}}>
        {steps.map((step, index) => {
          const opacity = fadeIn(frame, 8 + index * 6, 10);
          const shift = slideUp(frame, 8 + index * 6, 18, 10);
          const accent = stepColor(theme, step.tone);
          const isLast = index === steps.length - 1;

          return (
            <div key={`${step.title}-${index}`} style={{display: 'flex', flex: 1, gap: 18, alignItems: 'center'}}>
              <div
                style={{
                  flex: 1,
                  minHeight: 260,
                  borderRadius: 26,
                  padding: 22,
                  background: 'rgba(255,255,255,0.05)',
                  boxShadow: `inset 0 0 0 1px ${accent}`,
                  display: 'flex',
                  flexDirection: 'column',
                  gap: 14,
                  opacity,
                  transform: `translateY(${shift}px)`,
                }}
              >
                <div
                  style={{
                    color: accent,
                    fontFamily: fonts.mono,
                    fontSize: 17,
                    letterSpacing: '0.12em',
                    textTransform: 'uppercase',
                  }}
                >
                  {step.eyebrow}
                </div>
                <div
                  style={{
                    color: palette.text,
                    fontFamily: fonts.display,
                    fontSize: 34,
                    lineHeight: 1.08,
                    letterSpacing: '-0.04em',
                    fontWeight: 700,
                  }}
                >
                  {step.title}
                </div>
                <div
                  style={{
                    color: palette.textDim,
                    fontFamily: fonts.body,
                    fontSize: 22,
                    lineHeight: 1.35,
                  }}
                >
                  {step.body}
                </div>
              </div>
              {isLast ? null : (
                <div
                  style={{
                    width: 44,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    color: accent,
                    fontFamily: fonts.display,
                    fontSize: 44,
                    opacity,
                  }}
                >
                  →
                </div>
              )}
            </div>
          );
        })}
      </div>
    </GlassPanel>
  );
};
