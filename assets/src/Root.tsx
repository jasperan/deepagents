import type {ComponentProps} from 'react';

import {Composition, Folder} from 'remotion';

import {CampaignVideo} from './compositions/CampaignVideo';
import {CampaignVideoVertical} from './compositions/CampaignVideoVertical';
import {campaignVideos} from './data/videos';
import {verticalCampaignVideos} from './data/videos-vertical';
import {CANVAS, SOCIAL_CANVAS} from './design/tokens';

export const RemotionRoot = () => {
  return (
    <>
      <Folder name="oracle-campaign">
        {campaignVideos.map((video) => (
          <Composition
            key={video.renderId}
            id={video.renderId}
            component={CampaignVideo}
            width={CANVAS.width}
            height={CANVAS.height}
            fps={CANVAS.fps}
            durationInFrames={video.durationInFrames ?? CANVAS.defaultDurationInFrames}
            defaultProps={{videoId: video.id} satisfies ComponentProps<typeof CampaignVideo>}
          />
        ))}
      </Folder>
      <Folder name="oracle-campaign-vertical">
        {verticalCampaignVideos.map((video) => (
          <Composition
            key={video.renderId}
            id={video.renderId}
            component={CampaignVideoVertical}
            width={SOCIAL_CANVAS.width}
            height={SOCIAL_CANVAS.height}
            fps={SOCIAL_CANVAS.fps}
            durationInFrames={video.durationInFrames}
            defaultProps={{videoId: video.id} satisfies ComponentProps<typeof CampaignVideoVertical>}
          />
        ))}
      </Folder>
    </>
  );
};
