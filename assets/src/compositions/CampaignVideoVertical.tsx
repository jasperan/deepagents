import {fade} from '@remotion/transitions/fade';
import {slide} from '@remotion/transitions/slide';
import {TransitionSeries, linearTiming} from '@remotion/transitions';
import {useCurrentFrame, useVideoConfig} from 'remotion';

import {MusicBed} from '../components/MusicBed';
import {VerticalHeader} from '../components/VerticalHeader';
import {VerticalOutro} from '../components/VerticalOutro';
import {VerticalPanelRenderer} from '../components/VerticalPanelRenderer';
import {VerticalPhoneFrame} from '../components/VerticalPhoneFrame';
import {VerticalStatStack} from '../components/VerticalStatStack';
import {getVerticalVideoById} from '../data/videos-vertical';
import {fonts} from '../design/typography';
import {palette, themeMap} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';

export type CampaignVideoVerticalProps = {
  videoId: string;
};

const splitDurations = (total: number) => {
  const transition = 10;
  const content = total - transition * 3;
  const hook = Math.round(content * 0.22);
  const proof = Math.round(content * 0.3);
  const detail = Math.round(content * 0.28);
  const outro = content - hook - proof - detail;

  return {hook, proof, detail, outro, transition};
};

const SectionLabel = ({
  text,
  theme,
}: {
  text: string;
  theme: ReturnType<typeof getVerticalVideoById>['theme'];
}) => {
  const frame = useCurrentFrame();
  const opacity = fadeIn(frame, 0, 10);
  const shift = slideUp(frame, 0, 12, 10);
  const colors = themeMap[theme];

  return (
    <div
      style={{
        color: colors.accent,
        fontFamily: fonts.mono,
        fontSize: 18,
        letterSpacing: '0.14em',
        textTransform: 'uppercase',
        opacity,
        transform: `translateY(${shift}px)`,
      }}
    >
      {text}
    </div>
  );
};

const CaptionText = ({
  text,
}: {
  text: string;
}) => {
  const frame = useCurrentFrame();
  const opacity = fadeIn(frame, 8, 10);
  const shift = slideUp(frame, 8, 12, 10);

  return (
    <div
      style={{
        color: palette.textDim,
        fontFamily: fonts.body,
        fontSize: 28,
        lineHeight: 1.3,
        opacity,
        transform: `translateY(${shift}px)`,
      }}
    >
      {text}
    </div>
  );
};

const HookScene = ({videoId}: CampaignVideoVerticalProps) => {
  const video = getVerticalVideoById(videoId);

  return (
    <VerticalPhoneFrame theme={video.theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 28, justifyContent: 'center', flex: 1}}>
        <VerticalHeader
          number={video.number}
          kicker="native 9:16"
          hook={video.hook}
          subhook={video.subhook}
          theme={video.theme}
        />
      </div>
    </VerticalPhoneFrame>
  );
};

const ProofScene = ({videoId}: CampaignVideoVerticalProps) => {
  const video = getVerticalVideoById(videoId);

  return (
    <VerticalPhoneFrame theme={video.theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 24, flex: 1}}>
        <SectionLabel text={video.proofTitle} theme={video.theme} />
        <VerticalPanelRenderer panel={video.proofPanel} theme={video.theme} />
        <CaptionText text={video.proofCaption} />
      </div>
    </VerticalPhoneFrame>
  );
};

const DetailScene = ({videoId}: CampaignVideoVerticalProps) => {
  const video = getVerticalVideoById(videoId);

  return (
    <VerticalPhoneFrame theme={video.theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: 24, flex: 1}}>
        <SectionLabel text={video.detailTitle} theme={video.theme} />
        <div style={{display: 'flex', flexDirection: 'column', gap: 20, flex: 1}}>
          <VerticalPanelRenderer panel={video.detailPanel} theme={video.theme} />
          <CaptionText text={video.detailCaption} />
          <VerticalStatStack chips={video.statChips} theme={video.theme} />
        </div>
      </div>
    </VerticalPhoneFrame>
  );
};

const OutroScene = ({videoId}: CampaignVideoVerticalProps) => {
  const video = getVerticalVideoById(videoId);

  return (
    <VerticalPhoneFrame theme={video.theme}>
      <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'flex-end', gap: 28, flex: 1}}>
        <VerticalOutro line={video.closeLine} cta={video.closeCta} theme={video.theme} />
      </div>
    </VerticalPhoneFrame>
  );
};

export const CampaignVideoVertical = ({videoId}: CampaignVideoVerticalProps) => {
  const video = getVerticalVideoById(videoId);
  const {fps} = useVideoConfig();
  const parts = splitDurations(video.durationInFrames);
  const timing = linearTiming({durationInFrames: parts.transition});

  return (
    <>
      <MusicBed theme={video.theme} />
      <TransitionSeries>
        <TransitionSeries.Sequence durationInFrames={parts.hook}>
          <HookScene videoId={video.id} />
        </TransitionSeries.Sequence>
        <TransitionSeries.Transition presentation={fade()} timing={timing} />
        <TransitionSeries.Sequence durationInFrames={parts.proof}>
          <ProofScene videoId={video.id} />
        </TransitionSeries.Sequence>
        <TransitionSeries.Transition
          presentation={slide({direction: 'from-bottom'})}
          timing={linearTiming({durationInFrames: Math.max(8, timing.getDurationInFrames({fps}))})}
        />
        <TransitionSeries.Sequence durationInFrames={parts.detail}>
          <DetailScene videoId={video.id} />
        </TransitionSeries.Sequence>
        <TransitionSeries.Transition presentation={fade()} timing={timing} />
        <TransitionSeries.Sequence durationInFrames={parts.outro}>
          <OutroScene videoId={video.id} />
        </TransitionSeries.Sequence>
      </TransitionSeries>
    </>
  );
};
