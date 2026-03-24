import {useCurrentFrame} from 'remotion';

import type {RouterPanel as RouterPanelType} from '../data/videos';
import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, softSpring} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

const toneColor = (theme: ThemeName, tone?: 'oracle' | 'vector' | 'cloud') => {
  if (tone === 'vector') {
    return themeMap.vector.accent;
  }

  if (tone === 'cloud') {
    return themeMap.cloud.accent;
  }

  return themeMap[theme].accent;
};

export const PathRouterDiagram = ({
  title,
  targets,
  routes,
  footer,
  theme,
}: RouterPanelType & {theme: ThemeName}) => {
  const frame = useCurrentFrame();
  const progress = softSpring(frame, 30, 4, 36);

  const sourceX = 180;
  const targetX = 1180;
  const top = 150;
  const sourceGap = 96;
  const targetGap = 220;

  const targetPositions = new Map(
    targets.map((target, index) => [target, top + index * targetGap]),
  );

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
      <div style={{position: 'relative', flex: 1}}>
        <svg
          width="100%"
          height="100%"
          viewBox="0 0 1520 560"
          style={{position: 'absolute', inset: 0}}
        >
          {routes.map((route, index) => {
            const sourceY = top + index * sourceGap;
            const targetY = targetPositions.get(route.target) ?? top;
            const color = toneColor(theme, route.tone);
            const path = `M ${sourceX + 280} ${sourceY + 24} C 700 ${sourceY + 24}, 780 ${targetY + 24}, ${targetX} ${targetY + 24}`;
            const length = 1200;

            return (
              <path
                key={`${route.path}-${route.target}`}
                d={path}
                fill="none"
                stroke={color}
                strokeWidth={4}
                strokeDasharray={length}
                strokeDashoffset={length - length * progress}
                strokeLinecap="round"
                opacity={0.85}
              />
            );
          })}
        </svg>
        <div style={{position: 'absolute', inset: 0}}>
          {routes.map((route, index) => {
            const opacity = fadeIn(frame, 6 + index * 4, 10);
            const y = top + index * sourceGap;
            const color = toneColor(theme, route.tone);

            return (
              <div
                key={`${route.path}-${route.target}-source`}
                style={{
                  position: 'absolute',
                  left: sourceX,
                  top: y,
                  width: 280,
                  padding: '18px 20px',
                  borderRadius: 22,
                  border: `1px solid ${color}`,
                  background: 'rgba(255,255,255,0.05)',
                  color: palette.text,
                  fontFamily: fonts.mono,
                  fontSize: 24,
                  opacity,
                }}
              >
                {route.path}
              </div>
            );
          })}
          {targets.map((target, index) => {
            const opacity = fadeIn(frame, 12 + index * 6, 10);
            const y = top + index * targetGap;

            return (
              <div
                key={target}
                style={{
                  position: 'absolute',
                  left: targetX,
                  top: y,
                  width: 260,
                  padding: '22px 22px',
                  borderRadius: 24,
                  background: 'rgba(255,255,255,0.08)',
                  boxShadow: `inset 0 0 0 1px ${palette.panelBorder}`,
                  color: palette.text,
                  fontFamily: fonts.display,
                  fontSize: 30,
                  fontWeight: 700,
                  lineHeight: 1.15,
                  opacity,
                }}
              >
                {target}
              </div>
            );
          })}
        </div>
      </div>
      {footer ? (
        <div
          style={{
            color: palette.textDim,
            fontFamily: fonts.body,
            fontSize: 23,
            lineHeight: 1.4,
          }}
        >
          {footer}
        </div>
      ) : null}
    </GlassPanel>
  );
};
