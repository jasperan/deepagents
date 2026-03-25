# Deep Agents Oracle Video Campaign

Remotion workspace for the `deepagents-oracle` marketing series.

## Commands

```bash
cd assets
npm install
npm run dev
npm run compositions
npm run typecheck
npm run render:all
npm run render:posters
```

## What this workspace contains

- 10 campaign videos in 16:9
- a shared visual system and reusable scene components
- content grounded in real repo artifacts from the Oracle fork

## Notes

- Renders go to `assets/renders/`
- Poster frames render to `assets/renders/posters/`
- The compositions now include theme-based music beds from `assets/public/audio/`
- Voiceover timing lives in `assets/voiceover/campaign-voiceover.md`
- Generated TTS is not included yet. The current pass prepares scripts and timings for either human recording or later ElevenLabs generation.
- Future screenshots, logos, and final music/voiceover assets can be dropped into `assets/public/`
