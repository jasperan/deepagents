export const CANVAS = {
  width: 1920,
  height: 1080,
  fps: 30,
  defaultDurationInFrames: 660,
  heroDurationInFrames: 720,
} as const;

export const SOCIAL_CANVAS = {
  width: 1080,
  height: 1920,
  fps: 30,
} as const;

export const palette = {
  bg: '#070B14',
  bgSoft: '#101827',
  panel: 'rgba(10, 16, 28, 0.72)',
  panelStrong: 'rgba(12, 18, 32, 0.92)',
  panelBorder: 'rgba(255,255,255,0.10)',
  text: '#F4F7FB',
  textDim: 'rgba(244,247,251,0.70)',
  textMute: 'rgba(244,247,251,0.50)',
  oracle: '#F80000',
  oracleSoft: '#FF6A5A',
  cyan: '#4ED8FF',
  blue: '#4D7CFE',
  violet: '#8B6CFF',
  green: '#59F6B4',
  gold: '#F8C04A',
  warning: '#FF9B61',
  shadow: '0 28px 80px rgba(0, 0, 0, 0.45)',
  shadowSoft: '0 16px 48px rgba(0, 0, 0, 0.28)',
} as const;

export const themeMap = {
  oracle: {
    accent: palette.oracle,
    accent2: palette.oracleSoft,
    glow: 'rgba(248, 0, 0, 0.35)',
  },
  vector: {
    accent: palette.cyan,
    accent2: palette.blue,
    glow: 'rgba(78, 216, 255, 0.35)',
  },
  cloud: {
    accent: palette.blue,
    accent2: palette.violet,
    glow: 'rgba(77, 124, 254, 0.35)',
  },
  gold: {
    accent: palette.gold,
    accent2: palette.oracleSoft,
    glow: 'rgba(248, 192, 74, 0.32)',
  },
  violet: {
    accent: palette.violet,
    accent2: palette.cyan,
    glow: 'rgba(139, 108, 255, 0.32)',
  },
} as const;

export type ThemeName = keyof typeof themeMap;

export const spacing = {
  xs: 10,
  sm: 16,
  md: 24,
  lg: 32,
  xl: 48,
  xxl: 64,
  xxxl: 88,
} as const;

export const radii = {
  sm: 18,
  md: 26,
  lg: 34,
  pill: 999,
} as const;

export const sceneDurations = {
  intro: 132,
  evidence: 168,
  architecture: 168,
  benefits: 132,
  outro: 108,
  transition: 12,
} as const;
