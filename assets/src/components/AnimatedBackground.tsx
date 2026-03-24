import {AbsoluteFill, useCurrentFrame} from 'remotion';

import {palette, themeMap, type ThemeName} from '../design/tokens';
import {wave} from '../utils/animation';

export const AnimatedBackground = ({theme}: {theme: ThemeName}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  const x1 = 180 + wave(frame, 90, 0.018, 0.1);
  const y1 = 120 + wave(frame, 70, 0.022, 1.4);
  const x2 = 1480 + wave(frame, 120, 0.014, 0.8);
  const y2 = 760 + wave(frame, 80, 0.02, 2.2);
  const rotation = wave(frame, 3.5, 0.01, 0.7);

  return (
    <AbsoluteFill
      style={{
        background: `radial-gradient(circle at 18% 18%, rgba(255,255,255,0.04), transparent 28%), linear-gradient(135deg, ${palette.bg} 0%, ${palette.bgSoft} 100%)`,
        overflow: 'hidden',
      }}
    >
      <div
        style={{
          position: 'absolute',
          inset: -220,
          backgroundImage:
            'linear-gradient(rgba(255,255,255,0.06) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.06) 1px, transparent 1px)',
          backgroundSize: '96px 96px',
          opacity: 0.12,
          transform: `rotate(${rotation}deg) scale(1.08)`,
        }}
      />
      <div
        style={{
          position: 'absolute',
          width: 620,
          height: 620,
          borderRadius: 9999,
          left: x1,
          top: y1,
          filter: 'blur(140px)',
          background: colors.glow,
          opacity: 0.95,
        }}
      />
      <div
        style={{
          position: 'absolute',
          width: 560,
          height: 560,
          borderRadius: 9999,
          left: x2,
          top: y2,
          filter: 'blur(150px)',
          background: 'rgba(77, 124, 254, 0.22)',
          opacity: 0.8,
        }}
      />
      <div
        style={{
          position: 'absolute',
          inset: 0,
          background: 'linear-gradient(180deg, rgba(7,11,20,0.1) 0%, rgba(7,11,20,0.48) 100%)',
        }}
      />
    </AbsoluteFill>
  );
};
