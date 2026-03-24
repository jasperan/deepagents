import {Easing, interpolate, spring} from 'remotion';

export const fadeIn = (frame: number, start: number, duration: number) => {
  return interpolate(frame, [start, start + duration], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.cubic),
  });
};

export const fadeOut = (frame: number, start: number, duration: number) => {
  return interpolate(frame, [start, start + duration], [1, 0], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.in(Easing.cubic),
  });
};

export const slideUp = (frame: number, start: number, distance: number, duration: number) => {
  return interpolate(frame, [start, start + duration], [distance, 0], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.cubic),
  });
};

export const slideRight = (frame: number, start: number, distance: number, duration: number) => {
  return interpolate(frame, [start, start + duration], [-distance, 0], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.cubic),
  });
};

export const softSpring = (frame: number, fps: number, delay = 0, durationInFrames = 32) => {
  return spring({
    frame: Math.max(0, frame - delay),
    fps,
    durationInFrames,
    config: {damping: 200},
  });
};

export const clamp = (value: number, min: number, max: number) => {
  return Math.min(max, Math.max(min, value));
};

export const wave = (frame: number, amplitude: number, speed = 0.03, offset = 0) => {
  return Math.sin(frame * speed + offset) * amplitude;
};
