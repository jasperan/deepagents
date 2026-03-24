import {loadFont as loadInter} from '@remotion/google-fonts/Inter';
import {loadFont as loadSora} from '@remotion/google-fonts/Sora';
import {loadFont as loadMono} from '@remotion/google-fonts/IBMPlexMono';

const inter = loadInter('normal', {
  weights: ['400', '500', '600', '700'],
  subsets: ['latin'],
});

const sora = loadSora('normal', {
  weights: ['400', '600', '700', '800'],
  subsets: ['latin'],
});

const mono = loadMono('normal', {
  weights: ['400', '500', '600'],
  subsets: ['latin'],
});

export const fonts = {
  body: inter.fontFamily,
  display: sora.fontFamily,
  mono: mono.fontFamily,
} as const;
