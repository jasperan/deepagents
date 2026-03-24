import type {ThemeName} from '../design/tokens';

export type TerminalLine = {
  text: string;
  kind?: 'command' | 'output' | 'success' | 'note' | 'warning';
};

export type CodePanel = {
  kind: 'code';
  title: string;
  language: string;
  lines: string[];
  highlight?: number[];
};

export type TerminalPanel = {
  kind: 'terminal';
  title: string;
  lines: TerminalLine[];
};

export type TablePanel = {
  kind: 'table';
  title: string;
  columns: string[];
  rows: string[][];
  highlightRow?: number;
};

export type CalloutPanel = {
  kind: 'callout';
  title: string;
  body: string;
  bullets?: string[];
  stat?: string;
};

export type ComparisonPanel = {
  kind: 'comparison';
  title: string;
  leftTitle: string;
  rightTitle: string;
  leftItems: string[];
  rightItems: string[];
};

export type FlowPanel = {
  kind: 'flow';
  title: string;
  steps: Array<{
    eyebrow: string;
    title: string;
    body: string;
    tone?: 'oracle' | 'vector' | 'cloud';
  }>;
};

export type RouterPanel = {
  kind: 'router';
  title: string;
  targets: string[];
  routes: Array<{
    path: string;
    target: string;
    tone?: 'oracle' | 'vector' | 'cloud';
  }>;
  footer?: string;
};

export type Panel = CodePanel | TerminalPanel | TablePanel | CalloutPanel | ComparisonPanel | FlowPanel | RouterPanel;

export type BenefitCard = {
  eyebrow: string;
  title: string;
  body: string;
};

export type VideoSpec = {
  id: string;
  renderId: string;
  number: number;
  theme: ThemeName;
  kicker: string;
  title: string;
  claim: string;
  introPoints: string[];
  evidenceTitle: string;
  leftPanel: Panel;
  rightPanel: Panel;
  evidenceBullets: string[];
  architectureTitle: string;
  architecturePanel: Panel;
  architectureBullets: string[];
  benefitTitle: string;
  benefitCards: BenefitCard[];
  closingLine: string;
  closingBadges: string[];
  closingCta: string;
  durationInFrames?: number;
};

export const campaignVideos: VideoSpec[] = [
  {
    id: 'memory-survives',
    renderId: 'Oracle01MemorySurvives',
    number: 1,
    theme: 'oracle',
    kicker: 'durable memory',
    title: 'Memory that survives the restart',
    claim: 'Deep Agents can pick up where it left off because Oracle keeps the memory alive after the process dies.',
    introPoints: ['Oracle-backed persistence', 'Cross-session continuity', 'Real restart proof'],
    evidenceTitle: 'From first run to resumed context',
    leftPanel: {
      kind: 'terminal',
      title: 'Session 1',
      lines: [
        {text: '$ python agent.py', kind: 'command'},
        {text: 'Research LangGraph and save findings', kind: 'output'},
        {text: 'Saved /memory/langgraph.md', kind: 'success'},
        {text: 'Session complete', kind: 'note'},
      ],
    },
    rightPanel: {
      kind: 'terminal',
      title: 'Session 2',
      lines: [
        {text: '$ python agent.py', kind: 'command'},
        {text: 'Reading /memory/langgraph.md', kind: 'output'},
        {text: 'Continuing from prior findings', kind: 'success'},
        {text: 'No manual restore step', kind: 'note'},
      ],
    },
    evidenceBullets: [
      'The repo example already writes findings into `/memory/...` files.',
      'Oracle turns that memory into durable data instead of process-local state.',
      'The restart becomes a continuation, not a reset.',
    ],
    architectureTitle: 'What stays alive in Oracle',
    architecturePanel: {
      kind: 'table',
      title: 'da_files',
      columns: ['namespace', 'file_path', 'role'],
      rows: [
        ['default', '/memory/langgraph.md', 'Agent memory'],
        ['default', '/history/thread-1.md', 'Conversation history'],
        ['default', '/skills/oracle.md', 'Loaded skill'],
      ],
      highlightRow: 0,
    },
    architectureBullets: [
      'Memory, history, and skills all land in the same durable table.',
      'The namespace key keeps agents partitioned cleanly.',
      'This is database-backed memory, not a fragile side cache.',
    ],
    benefitTitle: 'Why that matters',
    benefitCards: [
      {eyebrow: 'Resilience', title: 'Crash tolerant memory', body: 'State survives restarts, deploys, and broken sessions.'},
      {eyebrow: 'Continuity', title: 'Cross-session context', body: 'The next run starts informed instead of empty.'},
      {eyebrow: 'Credibility', title: 'Durability you can explain', body: 'It is easier to trust memory when you know where it lives.'},
    ],
    closingLine: 'Kill the process. Keep the memory.',
    closingBadges: ['ORACLE AI DATABASE', 'CROSS-SESSION', 'DEEP AGENTS'],
    closingCta: 'This fork makes memory durable on purpose.',
  },
  {
    id: 'routing',
    renderId: 'Oracle02Routing',
    number: 2,
    theme: 'cloud',
    kicker: 'composite backend',
    title: 'Where the memory actually lives',
    claim: 'The Oracle story is clean because persistence is routed by path, not smeared across the whole agent runtime.',
    introPoints: ['CompositeBackend', 'Path-based routing', 'Additive to upstream'],
    evidenceTitle: 'The one-line mental model',
    leftPanel: {
      kind: 'code',
      title: 'create_oracle_deep_agent()',
      language: 'python',
      lines: [
        'def _backend_factory(runtime: object) -> CompositeBackend:',
        '    return CompositeBackend(',
        '        default=StateBackend(runtime),',
        '        routes=routes,',
        '    )',
      ],
      highlight: [1, 2, 3, 4],
    },
    rightPanel: {
      kind: 'callout',
      title: 'What gets routed',
      body: 'Persistent paths go into Oracle. Scratch work stays in memory.',
      bullets: ['/memory/*', '/history/*', '/skills/*', 'everything else -> StateBackend'],
      stat: '3 durable path families',
    },
    evidenceBullets: [
      'You did not fork Deep Agents into a different architecture.',
      'You inserted Oracle exactly where persistence matters most.',
      'That keeps the story easy to teach and easy to trust.',
    ],
    architectureTitle: 'Routing as an animation, not a hand wave',
    architecturePanel: {
      kind: 'router',
      title: 'CompositeBackend flow',
      targets: ['OracleStoreBackend', 'StateBackend'],
      routes: [
        {path: '/memory/', target: 'OracleStoreBackend', tone: 'oracle'},
        {path: '/history/', target: 'OracleStoreBackend', tone: 'cloud'},
        {path: '/skills/', target: 'OracleStoreBackend', tone: 'oracle'},
        {path: '/tmp/', target: 'StateBackend', tone: 'vector'},
        {path: '/scratch/', target: 'StateBackend', tone: 'vector'},
      ],
      footer: 'Persistent state moves to Oracle. Ephemeral work stays fast and local.',
    },
    architectureBullets: [
      'The split is visible, explicit, and easy to debug.',
      'Oracle becomes the durable tier without swallowing temporary work.',
      'That keeps the fork close to upstream Deep Agents behavior.',
    ],
    benefitTitle: 'Why builders care',
    benefitCards: [
      {eyebrow: 'Control', title: 'Persistence with boundaries', body: 'Not every file needs a database round trip.'},
      {eyebrow: 'Clarity', title: 'Easy to reason about', body: 'A path prefix tells you exactly where state will land.'},
      {eyebrow: 'Adoption', title: 'Low-friction upgrade', body: 'The Oracle layer feels additive instead of invasive.'},
    ],
    closingLine: 'Persistence becomes a routing decision, not a rewrite.',
    closingBadges: ['COMPOSITEBACKEND', 'ORACLESTOREBACKEND', 'UPSTREAM-FRIENDLY'],
    closingCta: 'That is why the integration reads cleanly on screen and in code.',
  },
  {
    id: 'system-of-record',
    renderId: 'Oracle03SystemOfRecord',
    number: 3,
    theme: 'gold',
    kicker: 'system of record',
    title: 'Oracle as the system of record',
    claim: 'This memory story is believable because file operations resolve into rows, columns, timestamps, and commits inside Oracle.',
    introPoints: ['MERGE-based writes', 'Namespace isolation', 'Durable table model'],
    evidenceTitle: 'The write path is concrete',
    leftPanel: {
      kind: 'code',
      title: 'OracleStoreBackend.write()',
      language: 'sql',
      lines: [
        'MERGE INTO da_files tgt',
        'USING (SELECT :namespace, :file_path FROM DUAL) src',
        'WHEN MATCHED THEN UPDATE SET content = :content',
        'WHEN NOT MATCHED THEN INSERT (...)',
      ],
      highlight: [0, 2, 3],
    },
    rightPanel: {
      kind: 'callout',
      title: 'What that buys you',
      body: 'Memory writes are durable database operations, not opaque blobs hidden in a sidecar.',
      bullets: ['upsert semantics', 'timestamps on every write', 'single durable table for agent state'],
      stat: 'ACID-backed persistence',
    },
    evidenceBullets: [
      'The backend implements real file operations on top of Oracle SQL.',
      'Read, write, edit, glob, grep, upload, and download all land on a durable model.',
      'That is a better demo and a better production story.',
    ],
    architectureTitle: 'The table that makes the story real',
    architecturePanel: {
      kind: 'table',
      title: 'da_files schema',
      columns: ['column', 'type', 'why it matters'],
      rows: [
        ['namespace', 'VARCHAR2(512)', 'Agent isolation'],
        ['file_path', 'VARCHAR2(1024)', 'Stable address'],
        ['content', 'CLOB', 'Durable body'],
        ['modified_at', 'TIMESTAMP', 'Traceable writes'],
      ],
      highlightRow: 2,
    },
    architectureBullets: [
      'The schema is simple enough to explain and strong enough to grow from.',
      'A serious memory story starts with a serious persistence model.',
      'Oracle gives the fork a durable spine.',
    ],
    benefitTitle: 'What the database adds',
    benefitCards: [
      {eyebrow: 'Storage', title: 'Rows, not magic', body: 'Agent state lives in a table you can reason about.'},
      {eyebrow: 'Ops', title: 'Commits and timestamps', body: 'Persistence has visibility and a shape you can inspect.'},
      {eyebrow: 'Trust', title: 'A real source of truth', body: 'The memory layer stops feeling demo-only.'},
    ],
    closingLine: 'The memory is persuasive because the storage model is concrete.',
    closingBadges: ['DA_FILES', 'MERGE', 'TIMESTAMPS'],
    closingCta: 'Oracle turns agent memory into durable application state.',
  },
  {
    id: 'vector-search',
    renderId: 'Oracle04VectorSearch',
    number: 4,
    theme: 'vector',
    kicker: 'semantic memory',
    title: 'Semantic memory with Oracle AI Vector Search',
    claim: 'The stretch goal is the strongest differentiator: memory retrieval can become semantic, not just keyword-based.',
    introPoints: ['Embeddings on write', 'VECTOR column', 'Cosine similarity in-db'],
    evidenceTitle: 'The vector path already exists',
    leftPanel: {
      kind: 'code',
      title: 'OracleVectorBackend',
      language: 'python',
      lines: [
        'embedding_vector = self._embeddings.embed_query(content)',
        "cur.execute(_MERGE_WITH_EMBEDDING_SQL, {'embedding': embedding_str})",
        'result = backend.semantic_grep("resume previous findings")',
      ],
      highlight: [0, 1, 2],
    },
    rightPanel: {
      kind: 'callout',
      title: 'What changes for the agent',
      body: 'A natural-language query can pull the right memory even when the exact keywords are missing.',
      bullets: ['embed on write', 'rank by similarity', 'retrieve intent, not just string overlap'],
      stat: 'VECTOR_DISTANCE(COSINE)',
    },
    evidenceBullets: [
      'The code already stores embeddings in Oracle alongside the file content.',
      'Semantic retrieval makes memory feel smarter without adding a separate vector store.',
      'That is the Oracle AI Database moment in this fork.',
    ],
    architectureTitle: 'How the semantic query resolves',
    architecturePanel: {
      kind: 'comparison',
      title: 'Keyword lookup vs semantic retrieval',
      leftTitle: 'String match only',
      rightTitle: 'Oracle AI Vector Search',
      leftItems: ['Needs exact terms', 'Misses paraphrases', 'Shallow recall'],
      rightItems: ['Embeds the query', 'Ranks by cosine distance', 'Finds semantically related memory'],
    },
    architectureBullets: [
      'Semantic search moves the memory layer from static storage to active retrieval.',
      'Because it stays in Oracle, the architecture stays compact.',
      'This is the most future-facing part of the campaign.',
    ],
    benefitTitle: 'Why this lands',
    benefitCards: [
      {eyebrow: 'Recall', title: 'Better memory retrieval', body: 'Agents recover meaning, not just repeated strings.'},
      {eyebrow: 'Architecture', title: 'No extra vector service', body: 'Search stays in the same database story.'},
      {eyebrow: 'Differentiation', title: 'A real Oracle moment', body: 'This is where Oracle AI Database clearly changes the product.'},
    ],
    closingLine: 'The memory layer stops acting like a folder and starts acting like recall.',
    closingBadges: ['VECTOR', 'SEMANTIC GREP', 'IN-DATABASE'],
    closingCta: 'That is the part builders will remember.',
  },
  {
    id: 'one-line-setup',
    renderId: 'Oracle05OneLineSetup',
    number: 5,
    theme: 'violet',
    kicker: 'developer experience',
    title: 'One line to get an Oracle-backed agent',
    claim: 'The Oracle story is powerful, but the adoption story stays simple because the fork ships a convenience factory.',
    introPoints: ['Simple entrypoint', 'Oracle defaults', 'Low boilerplate'],
    evidenceTitle: 'What the developer writes',
    leftPanel: {
      kind: 'code',
      title: 'examples/oracle_agent/agent.py',
      language: 'python',
      lines: [
        'from deepagents_oracle import create_oracle_deep_agent',
        '',
        'agent = create_oracle_deep_agent(',
        '    system_prompt="Store important findings in /memory/"',
        ')',
      ],
      highlight: [0, 2, 3],
    },
    rightPanel: {
      kind: 'terminal',
      title: 'The launch path',
      lines: [
        {text: '$ pip install deepagents deepagents-oracle', kind: 'command'},
        {text: '$ python examples/oracle_agent/agent.py', kind: 'command'},
        {text: 'Oracle-backed agent online', kind: 'success'},
        {text: 'Memory writes enabled', kind: 'note'},
      ],
    },
    evidenceBullets: [
      'The factory keeps the demo path short enough to be persuasive.',
      'That matters for Python developers who do not want a setup thesis.',
      'You kept the on-ramp smooth while still making Oracle central.',
    ],
    architectureTitle: 'What the helper hides for you',
    architecturePanel: {
      kind: 'callout',
      title: 'create_oracle_deep_agent()',
      body: 'The helper wires config, connection manager, Oracle backend, route mapping, and Deep Agents defaults in one place.',
      bullets: ['OracleConfig()', 'OracleConnectionManager()', 'OracleStoreBackend()', 'CompositeBackend()'],
      stat: 'One opinionated path into the stack',
    },
    architectureBullets: [
      'The helper makes the first demo easy without hiding the architecture from advanced users.',
      'That is exactly what a strong partner package should do.',
      'You gave the fork a friendly surface area.',
    ],
    benefitTitle: 'Why this converts',
    benefitCards: [
      {eyebrow: 'Speed', title: 'Fast first success', body: 'A short setup path keeps curiosity alive.'},
      {eyebrow: 'Clarity', title: 'Good defaults', body: 'The Oracle story feels pre-wired instead of bolted on.'},
      {eyebrow: 'Reach', title: 'Broader audience', body: 'Not everyone wants to hand-assemble a backend graph.'},
    ],
    closingLine: 'A powerful storage story works better when the first line stays short.',
    closingBadges: ['CREATE_ORACLE_DEEP_AGENT', 'PYTHON', 'FAST START'],
    closingCta: 'That is how technical credibility turns into adoption.',
  },
  {
    id: 'local-26ai',
    renderId: 'Oracle06LocalWith26ai',
    number: 6,
    theme: 'oracle',
    kicker: 'local first',
    title: 'Oracle Database 26ai Free on your laptop',
    claim: 'The best part of the Oracle story is that it does not start in a cloud slide deck. It starts with a local container and a quick setup loop.',
    introPoints: ['Docker compose', 'FREEPDB1', 'Laptop-ready demo'],
    evidenceTitle: 'The local loop is short',
    leftPanel: {
      kind: 'terminal',
      title: 'oracle/docker-compose.yml',
      lines: [
        {text: '$ docker compose -f oracle/docker-compose.yml --profile freepdb up -d', kind: 'command'},
        {text: 'oracle-free started on 1521', kind: 'success'},
        {text: 'FREEPDB1 healthy', kind: 'note'},
      ],
    },
    rightPanel: {
      kind: 'terminal',
      title: 'oracle/scripts/setup-oracle.sh',
      lines: [
        {text: '$ bash oracle/scripts/setup-oracle.sh', kind: 'command'},
        {text: 'Creating user deepagents...', kind: 'output'},
        {text: 'Schema initialized.', kind: 'success'},
        {text: 'Setup complete.', kind: 'success'},
      ],
    },
    evidenceBullets: [
      'You made the local demo path easy to narrate and easy to repeat.',
      'That lowers the intimidation factor around Oracle immediately.',
      'A local-first story makes the rest of the campaign believable.',
    ],
    architectureTitle: 'What the local stack looks like',
    architecturePanel: {
      kind: 'code',
      title: 'docker-compose.yml',
      language: 'yaml',
      lines: [
        'image: container-registry.oracle.com/database/free:latest',
        'profiles: ["freepdb"]',
        'ports:',
        '  - "1521:1521"',
        'healthcheck: sqlplus ... FREEPDB1',
      ],
      highlight: [0, 1, 3, 4],
    },
    architectureBullets: [
      'The local story is not hand-wavy. It is codified in the repo.',
      'That gives the videos a real proof moment instead of a conceptual one.',
      'It also makes the fork easier to try and easier to demo live.',
    ],
    benefitTitle: 'Why local matters',
    benefitCards: [
      {eyebrow: 'Developer loop', title: 'Fast to stand up', body: 'One compose command and one setup script get the database ready.'},
      {eyebrow: 'Confidence', title: 'Easy to verify', body: 'You can show the persistence story on a laptop, not just in theory.'},
      {eyebrow: 'Reach', title: 'Lower barrier to entry', body: 'Oracle stops feeling distant when it runs next to the app.'},
    ],
    closingLine: 'The Oracle story starts on a laptop, not after a procurement meeting.',
    closingBadges: ['26AI FREE', 'DOCKER', 'LOCAL DEV'],
    closingCta: 'That is a better way to introduce a serious database to builders.',
  },
  {
    id: 'adb-bridge',
    renderId: 'Oracle07BridgeToADB',
    number: 7,
    theme: 'cloud',
    kicker: 'local to cloud',
    title: 'The bridge to Autonomous Database',
    claim: 'The cloud story works because it preserves the same mental model. You swap config, not the architecture.',
    introPoints: ['freepdb vs adb', 'Wallet support', 'Same code path'],
    evidenceTitle: 'Two modes, one shape',
    leftPanel: {
      kind: 'code',
      title: 'OracleConfig',
      language: 'python',
      lines: [
        'oracle_mode: Literal["freepdb", "adb"] = "freepdb"',
        'oracle_dsn: str | None = None',
        'oracle_wallet_path: str | None = None',
        'def get_dsn(self) -> str: ...',
      ],
      highlight: [0, 1, 2],
    },
    rightPanel: {
      kind: 'code',
      title: '.env transition',
      language: 'env',
      lines: [
        'DEEPAGENTS_ORACLE_MODE=freepdb',
        'DEEPAGENTS_ORACLE_HOST=localhost',
        '# switch environments',
        'DEEPAGENTS_ORACLE_MODE=adb',
        'DEEPAGENTS_ORACLE_DSN=(description=...)',
        'DEEPAGENTS_ORACLE_WALLET_PATH=/path/to/wallet',
      ],
      highlight: [0, 1, 3, 4, 5],
    },
    evidenceBullets: [
      'The local and cloud modes are visible right in config.',
      'That keeps the mental model stable while the target changes.',
      'A cloud bridge feels better when the code path barely flinches.',
    ],
    architectureTitle: 'The deployment ladder',
    architecturePanel: {
      kind: 'flow',
      title: 'One Oracle story, three steps',
      steps: [
        {
          eyebrow: 'Step 1',
          title: 'Build locally',
          body: 'Use Oracle Database 26ai Free with host, port, and service for the fastest dev loop.',
          tone: 'oracle',
        },
        {
          eyebrow: 'Step 2',
          title: 'Swap the connection mode',
          body: 'Flip `oracle_mode`, point at a DSN, and add the wallet path when needed.',
          tone: 'cloud',
        },
        {
          eyebrow: 'Step 3',
          title: 'Keep the same agent shape',
          body: 'The agent factory, backend wiring, and persistence story still look the same.',
          tone: 'vector',
        },
      ],
    },
    architectureBullets: [
      'The jump to ADB reads like a deployment step, not a new product.',
      'That makes the video stronger for both builders and platform audiences.',
      'It also gives the fork a much cleaner lifecycle story.',
    ],
    benefitTitle: 'Why that matters',
    benefitCards: [
      {eyebrow: 'Continuity', title: 'One architecture, two environments', body: 'The transition to cloud keeps the same shape.'},
      {eyebrow: 'Security', title: 'Wallet-aware path', body: 'ADB support is represented directly in config and connection handling.'},
      {eyebrow: 'Storytelling', title: 'Better product arc', body: 'The local demo naturally expands into a cloud deployment story.'},
    ],
    closingLine: 'The move to cloud feels like a config change because the architecture is already aligned.',
    closingBadges: ['ADB', 'WALLET SUPPORT', 'CONFIG-DRIVEN'],
    closingCta: 'That is the kind of bridge that keeps demos honest.',
  },
  {
    id: 'cleaner-architecture',
    renderId: 'Oracle08CleanerArchitecture',
    number: 8,
    theme: 'violet',
    kicker: 'fewer moving parts',
    title: 'Fewer moving parts, cleaner architecture',
    claim: 'A lot of agent stacks sprawl into too many services. Oracle lets this fork tell a more consolidated story.',
    introPoints: ['Less glue', 'One durable data plane', 'Cleaner operating story'],
    evidenceTitle: 'The contrast is visual',
    leftPanel: {
      kind: 'comparison',
      title: 'Architecture contrast',
      leftTitle: 'Typical stack',
      rightTitle: 'Oracle-centered stack',
      leftItems: ['agent harness', 'app database', 'vector store', 'cache layer', 'custom glue'],
      rightItems: ['agent harness', 'Oracle AI Database', 'persistent memory', 'semantic retrieval'],
    },
    rightPanel: {
      kind: 'table',
      title: 'What consolidates',
      columns: ['concern', 'typical', 'here'],
      rows: [
        ['memory', 'separate store', 'Oracle AI Database'],
        ['semantic retrieval', 'vector sidecar', 'Oracle AI Database'],
        ['health + pool', 'extra glue', 'OracleConnectionManager'],
      ],
      highlightRow: 1,
    },
    evidenceBullets: [
      'This is one of the best product-level messages in the whole fork.',
      'People feel architectural sprawl before they can name it.',
      'Oracle gives the system a tighter center of gravity.',
    ],
    architectureTitle: 'What the cleaner stack actually looks like',
    architecturePanel: {
      kind: 'flow',
      title: 'From sprawl to one durable data plane',
      steps: [
        {
          eyebrow: 'Before',
          title: 'Split persistence story',
          body: 'Memory, retrieval, and operational glue drift into separate boxes and separate seams.',
          tone: 'vector',
        },
        {
          eyebrow: 'Center',
          title: 'Oracle AI Database',
          body: 'Persistent memory, semantic retrieval, and durable state collapse into one stronger center.',
          tone: 'oracle',
        },
        {
          eyebrow: 'After',
          title: 'Cleaner operating shape',
          body: 'The app has fewer moving parts to explain, monitor, and keep stitched together.',
          tone: 'cloud',
        },
      ],
    },
    architectureBullets: [
      'The cleaner diagram is backed by real connection-pool and health-check code in the repo.',
      'That makes the simplification argument feel grounded instead of aspirational.',
      'Oracle is carrying both storage weight and architectural clarity here.',
    ],
    benefitTitle: 'Why teams care',
    benefitCards: [
      {eyebrow: 'Ops', title: 'Less glue to babysit', body: 'A tighter stack means fewer coordination points and fewer weird failures.'},
      {eyebrow: 'Narrative', title: 'Easier to explain internally', body: 'The architecture has a center of gravity.'},
      {eyebrow: 'Focus', title: 'More time on agent behavior', body: 'Less time disappears into stitching services together.'},
    ],
    closingLine: 'A cleaner memory story usually comes from a cleaner data story.',
    closingBadges: ['POOLING', 'HEALTH CHECK', 'LESS GLUE'],
    closingCta: 'Oracle gives the fork that center of gravity.',
  },
  {
    id: 'why-oracle-changes-it',
    renderId: 'Oracle09ProductionShaped',
    number: 9,
    theme: 'gold',
    kicker: 'broader product story',
    title: 'Why Oracle changes the Deep Agents story',
    claim: 'Deep Agents already ships planning, tools, sub-agents, and context management. Oracle makes that stack feel durable enough to carry further.',
    introPoints: ['Upstream strengths', 'Oracle as force multiplier', 'Production-shaped story'],
    evidenceTitle: 'What Deep Agents already brings',
    leftPanel: {
      kind: 'table',
      title: 'Upstream Deep Agents strengths',
      columns: ['capability', 'what it gives you'],
      rows: [
        ['planning', 'structured multi-step agent work'],
        ['tools', 'real interaction surface'],
        ['sub-agents', 'parallelizable execution'],
        ['context', 'memory and control scaffolding'],
      ],
      highlightRow: 2,
    },
    rightPanel: {
      kind: 'callout',
      title: 'What Oracle changes',
      body: 'Oracle upgrades the part agent demos usually undersell: durable memory, stronger retrieval, and a cleaner deployment narrative.',
      bullets: ['persistent memory', 'cross-session history', 'semantic retrieval', 'local-to-cloud continuity'],
      stat: 'The harness starts to feel production-shaped',
    },
    evidenceBullets: [
      'This is where the campaign widens from code proof to product narrative.',
      'Oracle does not replace the harness. It hardens the layer people usually trust the least.',
      'That changes how serious the whole stack feels.',
    ],
    architectureTitle: 'The capability stack after Oracle',
    architecturePanel: {
      kind: 'flow',
      title: 'How the system grows up',
      steps: [
        {
          eyebrow: 'Layer 1',
          title: 'Deep Agents harness',
          body: 'Planning, tools, sub-agents, and context management give the system its agent behavior.',
          tone: 'vector',
        },
        {
          eyebrow: 'Layer 2',
          title: 'Oracle memory layer',
          body: 'Durable memory and recoverable sessions give that behavior continuity across time.',
          tone: 'oracle',
        },
        {
          eyebrow: 'Layer 3',
          title: 'Semantic retrieval + cloud runway',
          body: 'Vector search and the bridge to ADB make the story sharper for real deployment conversations.',
          tone: 'cloud',
        },
      ],
    },
    architectureBullets: [
      'The Oracle layer makes the stack easier to position for real teams.',
      'It does not just add storage. It changes the perceived maturity of the system.',
      'That makes this a stronger setup for the final hero film.',
    ],
    benefitTitle: 'Why the broader audience cares',
    benefitCards: [
      {eyebrow: 'Product', title: 'Stronger story', body: 'The fork sounds more complete because the memory layer has weight.'},
      {eyebrow: 'Platform', title: 'Better operational narrative', body: 'There is a believable path from local demo to cloud deployment.'},
      {eyebrow: 'Advocacy', title: 'Better demos and content', body: 'The system now has sharper proof moments for every audience tier.'},
    ],
    closingLine: 'Oracle gives Deep Agents a more durable second act.',
    closingBadges: ['PLANNING', 'MEMORY', 'CLOUD RUNWAY'],
    closingCta: 'That is why this fork is more than a database adapter.',
  },
  {
    id: 'hero-film',
    renderId: 'Oracle10HeroFilm',
    number: 10,
    theme: 'vector',
    kicker: 'hero film',
    title: 'Deep Agents + Oracle AI Database',
    claim: 'A more durable way to build agent systems: memory that persists, retrieval that gets smarter, and a path from local to cloud that actually lines up.',
    introPoints: ['Durable memory', 'Semantic retrieval', 'Local to cloud'],
    evidenceTitle: 'The strongest proof beats, compressed',
    leftPanel: {
      kind: 'router',
      title: 'What Oracle carries',
      targets: ['Oracle AI Database', 'StateBackend'],
      routes: [
        {path: '/memory/', target: 'Oracle AI Database', tone: 'oracle'},
        {path: '/history/', target: 'Oracle AI Database', tone: 'cloud'},
        {path: '/skills/', target: 'Oracle AI Database', tone: 'oracle'},
        {path: 'semantic_grep()', target: 'Oracle AI Database', tone: 'vector'},
        {path: '/tmp/', target: 'StateBackend', tone: 'vector'},
      ],
      footer: 'Durable state and semantic retrieval converge on the same center.',
    },
    rightPanel: {
      kind: 'callout',
      title: 'What the campaign proved',
      body: 'Restart persistence, path routing, SQL-backed storage, vector search, one-line setup, local 26ai Free, and the bridge to ADB.',
      bullets: ['proof, not slogans', 'technical ladder', 'broader product close'],
      stat: '10-part campaign',
    },
    evidenceBullets: [
      'The best marketing angle here is earned credibility.',
      'Oracle is not a side note. It is the reason the memory story lands.',
      'That lets the final film feel broad without losing technical weight.',
    ],
    architectureTitle: 'The end-state message',
    architecturePanel: {
      kind: 'comparison',
      title: 'What this fork turns into',
      leftTitle: 'Interesting harness',
      rightTitle: 'Durable agent platform story',
      leftItems: ['ephemeral demos', 'memory skepticism', 'split persistence story'],
      rightItems: ['durable memory', 'semantic recall', 'local-to-cloud continuity'],
    },
    architectureBullets: [
      'The fork becomes easier to pitch to builders, platform teams, and Oracle audiences at the same time.',
      'That is what the whole ladder was designed to prove.',
      'The hero film should feel like a synthesis, not a recap.',
    ],
    benefitTitle: 'What Oracle brings to the table',
    benefitCards: [
      {eyebrow: 'Memory', title: 'Durability', body: 'The system remembers after the process disappears.'},
      {eyebrow: 'Retrieval', title: 'Semantic recall', body: 'The vector path raises the ceiling on useful memory.'},
      {eyebrow: 'Deployment', title: 'A real runway', body: 'The same story scales from local 26ai Free to ADB.'},
    ],
    closingLine: 'Deep Agents gets the harness. Oracle AI Database gives it staying power.',
    closingBadges: ['DEEP AGENTS', 'ORACLE AI DATABASE', 'MARKETING MASTER'],
    closingCta: 'That is the campaign. Now render it and ship it.',
    durationInFrames: 720,
  },
];

export const getVideoById = (videoId: string) => {
  const match = campaignVideos.find((video) => video.id === videoId);

  if (!match) {
    throw new Error(`Unknown video id: ${videoId}`);
  }

  return match;
};
