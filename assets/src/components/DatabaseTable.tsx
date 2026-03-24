import {useCurrentFrame} from 'remotion';

import {fonts} from '../design/typography';
import {palette, spacing, themeMap, type ThemeName} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';
import {GlassPanel} from './GlassPanel';

export const DatabaseTable = ({
  title,
  columns,
  rows,
  highlightRow,
  theme,
}: {
  title: string;
  columns: string[];
  rows: string[][];
  highlightRow?: number;
  theme: ThemeName;
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <GlassPanel
      style={{
        height: '100%',
        padding: 28,
        display: 'flex',
        flexDirection: 'column',
        gap: spacing.md,
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
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: `repeat(${columns.length}, minmax(0, 1fr))`,
          gap: 12,
          padding: '12px 16px',
          borderRadius: 18,
          background: 'rgba(255,255,255,0.05)',
        }}
      >
        {columns.map((column) => (
          <div
            key={column}
            style={{
              color: palette.textMute,
              fontFamily: fonts.mono,
              fontSize: 17,
              letterSpacing: '0.12em',
              textTransform: 'uppercase',
            }}
          >
            {column}
          </div>
        ))}
      </div>
      <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
        {rows.map((row, rowIndex) => {
          const start = 8 + rowIndex * 6;
          const opacity = fadeIn(frame, start, 10);
          const shift = slideUp(frame, start, 14, 10);
          const isHot = highlightRow === rowIndex;

          return (
            <div
              key={`${rowIndex}-${row.join('-')}`}
              style={{
                display: 'grid',
                gridTemplateColumns: `repeat(${columns.length}, minmax(0, 1fr))`,
                gap: 12,
                padding: '18px 16px',
                borderRadius: 20,
                background: isHot ? 'rgba(248, 0, 0, 0.08)' : 'rgba(255,255,255,0.03)',
                boxShadow: isHot ? `inset 0 0 0 1px ${colors.accent}` : 'none',
                opacity,
                transform: `translateY(${shift}px)`,
              }}
            >
              {row.map((cell, cellIndex) => (
                <div
                  key={`${rowIndex}-${cellIndex}-${cell}`}
                  style={{
                    color: palette.text,
                    fontFamily: cellIndex === 0 ? fonts.mono : fonts.body,
                    fontSize: 22,
                    lineHeight: 1.3,
                  }}
                >
                  {cell}
                </div>
              ))}
            </div>
          );
        })}
      </div>
    </GlassPanel>
  );
};
