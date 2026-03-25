import type {ReactNode} from 'react';

import {AbsoluteFill} from 'remotion';

import {AnimatedBackground} from './AnimatedBackground';
import {fonts} from '../design/typography';
import {palette, type ThemeName} from '../design/tokens';

export const VerticalPhoneFrame = ({
  theme,
  children,
}: {
  theme: ThemeName;
  children: ReactNode;
}) => {
  return (
    <AbsoluteFill>
      <AnimatedBackground theme={theme} />
      <AbsoluteFill
        style={{
          padding: '88px 54px 72px',
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
        }}
      >
        {children}
        <div
          style={{
            position: 'absolute',
            left: 54,
            right: 54,
            bottom: 28,
            display: 'flex',
            justifyContent: 'space-between',
            color: palette.textMute,
            fontFamily: fonts.mono,
            fontSize: 16,
            letterSpacing: '0.14em',
            textTransform: 'uppercase',
          }}
        >
          <span>deepagents-oracle</span>
          <span>9:16 social</span>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
