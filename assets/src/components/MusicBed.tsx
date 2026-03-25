import {Audio} from '@remotion/media';
import {interpolate, staticFile, useVideoConfig} from 'remotion';

import type {ThemeName} from '../design/tokens';

const trackForTheme: Record<ThemeName, string> = {
  oracle: 'audio/oracle-bed.mp3',
  gold: 'audio/oracle-bed.mp3',
  cloud: 'audio/cloud-bed.mp3',
  violet: 'audio/cloud-bed.mp3',
  vector: 'audio/vector-bed.mp3',
};

const peakVolumeForTheme: Record<ThemeName, number> = {
  oracle: 0.08,
  gold: 0.075,
  cloud: 0.072,
  violet: 0.07,
  vector: 0.078,
};

export const MusicBed = ({theme}: {theme: ThemeName}) => {
  const {durationInFrames} = useVideoConfig();
  const peak = peakVolumeForTheme[theme];

  return (
    <Audio
      src={staticFile(trackForTheme[theme])}
      volume={(f) =>
        interpolate(
          f,
          [0, 24, Math.max(48, durationInFrames - 54), durationInFrames],
          [0, peak, peak, 0],
          {
            extrapolateLeft: 'clamp',
            extrapolateRight: 'clamp',
          },
        )
      }
    />
  );
};
