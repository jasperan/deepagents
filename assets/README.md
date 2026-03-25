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
npm run render:all-vertical
npm run render:posters-vertical
```

## What this workspace contains

- 10 campaign videos in 16:9
- 10 native 9:16 social variants with mixed runtimes
- a shared visual system and reusable scene components
- content grounded in real repo artifacts from the Oracle fork

## Notes

- Renders go to `assets/renders/`
- Poster frames render to `assets/renders/posters/`
- Vertical poster frames render to `assets/renders/posters-vertical/`
- Vertical masters render to `assets/renders/vertical/`
- The compositions now include theme-based music beds from `assets/public/audio/`
- Voiceover timing lives in `assets/voiceover/campaign-voiceover.md`
- Generated TTS is not included yet. The current pass prepares scripts and timings for either human recording or later ElevenLabs generation.
- The 9:16 variants are rebuilt as native vertical edits, not cropped 16:9 masters.
- Future screenshots, logos, and final music/voiceover assets can be dropped into `assets/public/`
