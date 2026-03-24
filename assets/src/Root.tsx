import type {ComponentProps} from 'react';

import {Composition, Folder} from 'remotion';

import {CampaignVideo} from './compositions/CampaignVideo';
import {campaignVideos} from './data/videos';
import {CANVAS} from './design/tokens';

export const RemotionRoot = () => {
  return (
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
  );
};
