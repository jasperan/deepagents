import type {CSSProperties, ReactNode} from 'react';

import {palette, radii} from '../design/tokens';

export const GlassPanel = ({
  children,
  style,
}: {
  children: ReactNode;
  style?: CSSProperties;
}) => {
  return (
    <div
      style={{
        background: palette.panel,
        border: `1px solid ${palette.panelBorder}`,
        borderRadius: radii.lg,
        boxShadow: palette.shadow,
        backdropFilter: 'blur(18px)',
        overflow: 'hidden',
        ...style,
      }}
    >
      {children}
    </div>
  );
};
