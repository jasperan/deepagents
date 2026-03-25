import type {Panel} from './videos';
import {getVideoById} from './videos';
import type {ThemeName} from '../design/tokens';

export type VerticalVideoSpec = {
  id: string;
  renderId: string;
  number: number;
  theme: ThemeName;
  durationInFrames: number;
  posterFrame: number;
  hook: string;
  subhook: string;
  proofTitle: string;
  proofPanel: Panel;
  proofCaption: string;
  detailTitle: string;
  detailPanel: Panel;
  detailCaption: string;
  statChips: string[];
  closeLine: string;
  closeCta: string;
};

export const verticalCampaignVideos: VerticalVideoSpec[] = [
  {
    id: 'memory-survives',
    renderId: 'OracleVertical01MemorySurvives',
    number: 1,
    theme: 'oracle',
    durationInFrames: 330,
    posterFrame: 150,
    hook: 'The process dies. The memory stays.',
    subhook: 'Oracle turns Deep Agents memory into durable state that survives the restart.',
    proofTitle: 'Restart proof',
    proofPanel: getVideoById('memory-survives').rightPanel,
    proofCaption: 'The second session resumes from stored memory instead of starting cold.',
    detailTitle: 'Where it lands',
    detailPanel: getVideoById('memory-survives').architecturePanel,
    detailCaption: 'Memory, history, and skills land in `da_files`, not a fragile side cache.',
    statChips: ['restart-safe', 'cross-session', 'oracle-backed'],
    closeLine: 'Kill the process. Keep the context.',
    closeCta: 'Deep Agents + Oracle AI Database',
  },
  {
    id: 'routing',
    renderId: 'OracleVertical02Routing',
    number: 2,
    theme: 'cloud',
    durationInFrames: 360,
    posterFrame: 180,
    hook: 'Persistence is a routing choice.',
    subhook: 'This fork stays clean because Oracle only carries the paths that matter.',
    proofTitle: 'One line mental model',
    proofPanel: getVideoById('routing').rightPanel,
    proofCaption: '`/memory`, `/history`, and `/skills` go durable. Scratch work stays fast and local.',
    detailTitle: 'Backend flow',
    detailPanel: getVideoById('routing').architecturePanel,
    detailCaption: 'The split is explicit, visible, and easy to explain on screen or in code.',
    statChips: ['compositebackend', 'oraclestorebackend', 'upstream-friendly'],
    closeLine: 'Oracle changes the memory tier, not the whole shape.',
    closeCta: 'Native persistence, additive architecture.',
  },
  {
    id: 'system-of-record',
    renderId: 'OracleVertical03SystemOfRecord',
    number: 3,
    theme: 'gold',
    durationInFrames: 390,
    posterFrame: 210,
    hook: 'Agent memory should have a real storage model.',
    subhook: 'Here, writes resolve into Oracle SQL, rows, timestamps, and commits.',
    proofTitle: 'The write path',
    proofPanel: getVideoById('system-of-record').leftPanel,
    proofCaption: '`MERGE` gives the memory layer durable upsert semantics instead of vague state handling.',
    detailTitle: 'The table',
    detailPanel: getVideoById('system-of-record').architecturePanel,
    detailCaption: 'The `da_files` schema gives the fork a concrete source of truth.',
    statChips: ['merge', 'da_files', 'timestamps'],
    closeLine: 'That is how memory stops feeling like a demo cache.',
    closeCta: 'Oracle as the system of record.',
  },
  {
    id: 'vector-search',
    renderId: 'OracleVertical04VectorSearch',
    number: 4,
    theme: 'vector',
    durationInFrames: 450,
    posterFrame: 240,
    hook: 'Keyword search is a weak memory model.',
    subhook: 'Oracle AI Vector Search lets the agent retrieve by meaning, not just repeated strings.',
    proofTitle: 'Embed on write',
    proofPanel: getVideoById('vector-search').leftPanel,
    proofCaption: 'The vector path is already in the fork through `OracleVectorBackend`.',
    detailTitle: 'Why it hits harder',
    detailPanel: getVideoById('vector-search').architecturePanel,
    detailCaption: 'Semantic retrieval stays inside Oracle instead of adding another vector service to babysit.',
    statChips: ['vector_distance', 'semantic recall', 'in-database'],
    closeLine: 'This is the Oracle moment builders remember.',
    closeCta: 'Smarter memory, same core stack.',
  },
  {
    id: 'one-line-setup',
    renderId: 'OracleVertical05OneLineSetup',
    number: 5,
    theme: 'violet',
    durationInFrames: 330,
    posterFrame: 150,
    hook: 'One line gets you an Oracle-backed agent.',
    subhook: 'The storage story is serious, but the first-run ergonomics stay simple.',
    proofTitle: 'The entrypoint',
    proofPanel: getVideoById('one-line-setup').leftPanel,
    proofCaption: '`create_oracle_deep_agent()` keeps the on-ramp short enough to convert curiosity into usage.',
    detailTitle: 'What it hides for you',
    detailPanel: getVideoById('one-line-setup').architecturePanel,
    detailCaption: 'Config, connection management, backend wiring, memory, and skills get pre-wired together.',
    statChips: ['python', 'fast start', 'oracle defaults'],
    closeLine: 'A stronger stack lands better when the first line stays clean.',
    closeCta: 'Fast start, durable memory.',
  },
  {
    id: 'local-26ai',
    renderId: 'OracleVertical06LocalWith26ai',
    number: 6,
    theme: 'oracle',
    durationInFrames: 360,
    posterFrame: 180,
    hook: 'This Oracle story starts on a laptop.',
    subhook: 'Oracle Database 26ai Free gives the fork a local dev loop that is actually demoable.',
    proofTitle: 'Bring up the stack',
    proofPanel: getVideoById('local-26ai').leftPanel,
    proofCaption: 'A Docker profile gets Oracle up locally without turning setup into a saga.',
    detailTitle: 'Finish the setup',
    detailPanel: getVideoById('local-26ai').rightPanel,
    detailCaption: 'One setup script creates the user and schema so the memory story works immediately.',
    statChips: ['26ai free', 'docker', 'local dev'],
    closeLine: 'A serious database feels different when it runs right next to the app.',
    closeCta: 'Local first, Oracle-backed.',
  },
  {
    id: 'adb-bridge',
    renderId: 'OracleVertical07BridgeToADB',
    number: 7,
    theme: 'cloud',
    durationInFrames: 390,
    posterFrame: 210,
    hook: 'Local to cloud should feel like a deployment step.',
    subhook: 'This fork keeps the same agent shape while the connection target changes.',
    proofTitle: 'The config shift',
    proofPanel: getVideoById('adb-bridge').rightPanel,
    proofCaption: 'Switch modes, point at the DSN, add wallet support, keep moving.',
    detailTitle: 'The three-step bridge',
    detailPanel: getVideoById('adb-bridge').architecturePanel,
    detailCaption: 'Build locally, swap the connection mode, keep the same agent model.',
    statChips: ['adb', 'wallet-ready', 'same architecture'],
    closeLine: 'Same model. Bigger runway.',
    closeCta: 'Oracle 26ai Free to Autonomous Database.',
  },
  {
    id: 'cleaner-architecture',
    renderId: 'OracleVertical08CleanerArchitecture',
    number: 8,
    theme: 'violet',
    durationInFrames: 390,
    posterFrame: 210,
    hook: 'Agent stacks sprawl fast.',
    subhook: 'Oracle pulls memory, semantic retrieval, and durable state toward one stronger center.',
    proofTitle: 'What consolidates',
    proofPanel: getVideoById('cleaner-architecture').rightPanel,
    proofCaption: 'The storage story gets tighter when Oracle carries more of the durable load.',
    detailTitle: 'From sprawl to center',
    detailPanel: getVideoById('cleaner-architecture').architecturePanel,
    detailCaption: 'Fewer seams, less glue, easier to explain internally.',
    statChips: ['less glue', 'one center', 'cleaner ops'],
    closeLine: 'Cleaner memory usually starts with cleaner data architecture.',
    closeCta: 'Oracle gives the stack a center of gravity.',
  },
  {
    id: 'why-oracle-changes-it',
    renderId: 'OracleVertical09ProductionShaped',
    number: 9,
    theme: 'gold',
    durationInFrames: 420,
    posterFrame: 225,
    hook: 'Deep Agents already gives you a strong harness.',
    subhook: 'Oracle hardens the part agent demos usually undersell: memory, continuity, and retrieval.',
    proofTitle: 'What was already strong',
    proofPanel: getVideoById('why-oracle-changes-it').leftPanel,
    proofCaption: 'Planning, tools, sub-agents, and context management are already in place.',
    detailTitle: 'What grows up',
    detailPanel: getVideoById('why-oracle-changes-it').architecturePanel,
    detailCaption: 'Oracle adds continuity across time and a better deployment story around that harness.',
    statChips: ['planning', 'durable memory', 'cloud runway'],
    closeLine: 'Oracle gives Deep Agents a more durable second act.',
    closeCta: 'Stronger harness, stronger memory layer.',
  },
  {
    id: 'hero-film',
    renderId: 'OracleVertical10HeroFilm',
    number: 10,
    theme: 'vector',
    durationInFrames: 480,
    posterFrame: 270,
    hook: 'Deep Agents gives you the harness.',
    subhook: 'Oracle AI Database gives it staying power: durable memory, semantic recall, and a real path from local to cloud.',
    proofTitle: 'What Oracle carries',
    proofPanel: getVideoById('hero-film').leftPanel,
    proofCaption: 'Memory, history, skills, and semantic retrieval converge on the same durable center.',
    detailTitle: 'What the campaign proved',
    detailPanel: getVideoById('hero-film').rightPanel,
    detailCaption: 'Restart persistence, path routing, vector search, local 26ai Free, and the bridge to ADB.',
    statChips: ['durability', 'semantic recall', 'deployment runway'],
    closeLine: 'Deep Agents gets the harness. Oracle gives it staying power.',
    closeCta: 'Built for phones. Built around Oracle.',
  },
];

export const getVerticalVideoById = (videoId: string) => {
  const match = verticalCampaignVideos.find((video) => video.id === videoId);

  if (!match) {
    throw new Error(`Unknown vertical video id: ${videoId}`);
  }

  return match;
};
