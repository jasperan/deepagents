import type {ReactNode} from 'react';

import {fade} from '@remotion/transitions/fade';
import {slide} from '@remotion/transitions/slide';
import {TransitionSeries, linearTiming} from '@remotion/transitions';
import {AbsoluteFill, interpolate, useCurrentFrame, useVideoConfig} from 'remotion';

import {AnimatedBackground} from '../components/AnimatedBackground';
import {FeatureCards} from '../components/FeatureCards';
import {OutroCard} from '../components/OutroCard';
import {PanelRenderer} from '../components/PanelRenderer';
import {SceneHeader} from '../components/SceneHeader';
import {getVideoById} from '../data/videos';
import {fonts} from '../design/typography';
import {palette, sceneDurations, spacing, themeMap} from '../design/tokens';
import {fadeIn, slideUp} from '../utils/animation';

export type CampaignVideoProps = {
  videoId: string;
};

const sectionTitleStyle = {
  fontFamily: fonts.display,
  fontSize: 58,
  lineHeight: 1.06,
  letterSpacing: '-0.05em',
  fontWeight: 700,
  color: palette.text,
} as const;

const SceneScaffold = ({
  theme,
  children,
}: {
  theme: ReturnType<typeof getVideoById>['theme'];
  children: ReactNode;
}) => {
  return (
    <AbsoluteFill>
      <AnimatedBackground theme={theme} />
      <AbsoluteFill
        style={{
          padding: '72px 88px 64px',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        {children}
      </AbsoluteFill>
    </AbsoluteFill>
  );
};

const BrandMark = ({theme}: {theme: ReturnType<typeof getVideoById>['theme']}) => {
  const colors = themeMap[theme];

  return (
    <div
      style={{
        position: 'absolute',
        right: 88,
        bottom: 44,
        color: palette.textMute,
        fontFamily: fonts.mono,
        fontSize: 18,
        letterSpacing: '0.14em',
        textTransform: 'uppercase',
      }}
    >
      deepagents-oracle <span style={{color: colors.accent}}>campaign master</span>
    </div>
  );
};

const SectionHeading = ({
  eyebrow,
  title,
  description,
  theme,
}: {
  eyebrow: string;
  title: string;
  description?: string;
  theme: ReturnType<typeof getVideoById>['theme'];
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];
  const opacity = fadeIn(frame, 0, 16);
  const shift = slideUp(frame, 0, 18, 16);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: 14,
        opacity,
        transform: `translateY(${shift}px)`,
      }}
    >
      <div
        style={{
          color: colors.accent,
          fontFamily: fonts.mono,
          fontSize: 20,
          letterSpacing: '0.14em',
          textTransform: 'uppercase',
        }}
      >
        {eyebrow}
      </div>
      <div style={sectionTitleStyle}>{title}</div>
      {description ? (
        <div
          style={{
            color: palette.textDim,
            fontFamily: fonts.body,
            fontSize: 28,
            lineHeight: 1.35,
            maxWidth: 1180,
          }}
        >
          {description}
        </div>
      ) : null}
    </div>
  );
};

const BulletStack = ({
  bullets,
  theme,
}: {
  bullets: string[];
  theme: ReturnType<typeof getVideoById>['theme'];
}) => {
  const frame = useCurrentFrame();
  const colors = themeMap[theme];

  return (
    <div style={{display: 'flex', flexDirection: 'column', gap: 18}}>
      {bullets.map((bullet, index) => {
        const start = 10 + index * 5;
        const opacity = fadeIn(frame, start, 10);
        const shift = slideUp(frame, start, 16, 10);

        return (
          <div
            key={bullet}
            style={{
              display: 'flex',
              gap: 14,
              opacity,
              transform: `translateY(${shift}px)`,
              padding: '18px 20px',
              borderRadius: 22,
              background: 'rgba(255,255,255,0.04)',
              boxShadow: `inset 0 0 0 1px ${palette.panelBorder}`,
            }}
          >
            <div
              style={{
                marginTop: 12,
                width: 10,
                height: 10,
                borderRadius: 999,
                background: colors.accent,
                flexShrink: 0,
              }}
            />
            <div
              style={{
                color: palette.text,
                fontFamily: fonts.body,
                fontSize: 24,
                lineHeight: 1.38,
              }}
            >
              {bullet}
            </div>
          </div>
        );
      })}
    </div>
  );
};

const IntroScene = ({videoId}: CampaignVideoProps) => {
  const video = getVideoById(videoId);
  const frame = useCurrentFrame();
  const colors = themeMap[video.theme];

  return (
    <SceneScaffold theme={video.theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: spacing.xxl, flex: 1, justifyContent: 'center'}}>
        <SceneHeader
          number={video.number}
          kicker={video.kicker}
          title={video.title}
          claim={video.claim}
          theme={video.theme}
        />
        <div style={{display: 'flex', gap: 18, flexWrap: 'wrap'}}>
          {video.introPoints.map((point, index) => {
            const opacity = fadeIn(frame, 14 + index * 5, 10);
            const y = slideUp(frame, 14 + index * 5, 18, 10);

            return (
              <div
                key={point}
                style={{
                  padding: '16px 22px',
                  borderRadius: 999,
                  background: 'rgba(255,255,255,0.05)',
                  boxShadow: `inset 0 0 0 1px ${colors.accent}`,
                  color: palette.text,
                  fontFamily: fonts.body,
                  fontSize: 22,
                  opacity,
                  transform: `translateY(${y}px)`,
                }}
              >
                {point}
              </div>
            );
          })}
        </div>
      </div>
      <BrandMark theme={video.theme} />
    </SceneScaffold>
  );
};

const EvidenceScene = ({videoId}: CampaignVideoProps) => {
  const video = getVideoById(videoId);

  return (
    <SceneScaffold theme={video.theme}>
      <SectionHeading
        eyebrow="Proof"
        title={video.evidenceTitle}
        description={video.claim}
        theme={video.theme}
      />
      <div style={{display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: spacing.lg, marginTop: spacing.xl, minHeight: 500}}>
        <PanelRenderer panel={video.leftPanel} theme={video.theme} />
        <PanelRenderer panel={video.rightPanel} theme={video.theme} />
      </div>
      <div style={{marginTop: spacing.xl}}>
        <BulletStack bullets={video.evidenceBullets} theme={video.theme} />
      </div>
      <BrandMark theme={video.theme} />
    </SceneScaffold>
  );
};

const ArchitectureScene = ({videoId}: CampaignVideoProps) => {
  const video = getVideoById(videoId);

  return (
    <SceneScaffold theme={video.theme}>
      <SectionHeading
        eyebrow="Architecture"
        title={video.architectureTitle}
        description="The repo structure and code paths make the Oracle story visible instead of abstract."
        theme={video.theme}
      />
      <div style={{display: 'grid', gridTemplateColumns: '1.7fr 1fr', gap: spacing.lg, marginTop: spacing.xl, flex: 1}}>
        <div style={{minHeight: 0}}>
          <PanelRenderer panel={video.architecturePanel} theme={video.theme} />
        </div>
        <div style={{display: 'flex', alignItems: 'stretch'}}>
          <BulletStack bullets={video.architectureBullets} theme={video.theme} />
        </div>
      </div>
      <BrandMark theme={video.theme} />
    </SceneScaffold>
  );
};

const BenefitsScene = ({videoId}: CampaignVideoProps) => {
  const video = getVideoById(videoId);
  const frame = useCurrentFrame();
  const opacity = fadeIn(frame, 0, 12);
  const y = slideUp(frame, 0, 18, 12);

  return (
    <SceneScaffold theme={video.theme}>
      <div style={{display: 'flex', flexDirection: 'column', gap: spacing.xl, flex: 1}}>
        <div style={{opacity, transform: `translateY(${y}px)`}}>
          <SectionHeading
            eyebrow="Why it matters"
            title={video.benefitTitle}
            description="This is where the technical proof widens into product value."
            theme={video.theme}
          />
        </div>
        <FeatureCards cards={video.benefitCards} theme={video.theme} />
      </div>
      <BrandMark theme={video.theme} />
    </SceneScaffold>
  );
};

const OutroScene = ({videoId}: CampaignVideoProps) => {
  const video = getVideoById(videoId);
  const frame = useCurrentFrame();
  const scale = interpolate(frame, [0, 40], [0.98, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
  });

  return (
    <SceneScaffold theme={video.theme}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          flex: 1,
          transform: `scale(${scale})`,
        }}
      >
        <OutroCard
          line={video.closingLine}
          badges={video.closingBadges}
          cta={video.closingCta}
          theme={video.theme}
        />
      </div>
      <BrandMark theme={video.theme} />
    </SceneScaffold>
  );
};

export const CampaignVideo = ({videoId}: CampaignVideoProps) => {
  const video = getVideoById(videoId);
  const {fps} = useVideoConfig();
  const timing = linearTiming({durationInFrames: sceneDurations.transition});

  return (
    <TransitionSeries>
      <TransitionSeries.Sequence durationInFrames={sceneDurations.intro}>
        <IntroScene videoId={video.id} />
      </TransitionSeries.Sequence>
      <TransitionSeries.Transition presentation={fade()} timing={timing} />
      <TransitionSeries.Sequence durationInFrames={sceneDurations.evidence}>
        <EvidenceScene videoId={video.id} />
      </TransitionSeries.Sequence>
      <TransitionSeries.Transition
        presentation={slide({direction: 'from-right'})}
        timing={linearTiming({durationInFrames: Math.max(10, timing.getDurationInFrames({fps}))})}
      />
      <TransitionSeries.Sequence durationInFrames={sceneDurations.architecture}>
        <ArchitectureScene videoId={video.id} />
      </TransitionSeries.Sequence>
      <TransitionSeries.Transition presentation={fade()} timing={timing} />
      <TransitionSeries.Sequence durationInFrames={sceneDurations.benefits}>
        <BenefitsScene videoId={video.id} />
      </TransitionSeries.Sequence>
      <TransitionSeries.Transition presentation={slide({direction: 'from-bottom'})} timing={timing} />
      <TransitionSeries.Sequence durationInFrames={video.durationInFrames ? video.durationInFrames - 552 : sceneDurations.outro}>
        <OutroScene videoId={video.id} />
      </TransitionSeries.Sequence>
    </TransitionSeries>
  );
};
