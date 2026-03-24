import {execSync} from 'node:child_process';
import {mkdirSync} from 'node:fs';

const renders = [
  ['Oracle01MemorySurvives', 'oracle-01-memory-survives.mp4'],
  ['Oracle02Routing', 'oracle-02-routing.mp4'],
  ['Oracle03SystemOfRecord', 'oracle-03-system-of-record.mp4'],
  ['Oracle04VectorSearch', 'oracle-04-vector-search.mp4'],
  ['Oracle05OneLineSetup', 'oracle-05-one-line-setup.mp4'],
  ['Oracle06LocalWith26ai', 'oracle-06-local-with-26ai.mp4'],
  ['Oracle07BridgeToADB', 'oracle-07-bridge-to-adb.mp4'],
  ['Oracle08CleanerArchitecture', 'oracle-08-cleaner-architecture.mp4'],
  ['Oracle09ProductionShaped', 'oracle-09-production-shaped.mp4'],
  ['Oracle10HeroFilm', 'oracle-10-hero-film.mp4'],
];

const dryRun = process.argv.includes('--dry-run');
mkdirSync('renders', {recursive: true});

for (const [id, output] of renders) {
  const cmd = `npx remotion render src/index.ts ${id} renders/${output}`;
  console.log(cmd);
  if (!dryRun) {
    execSync(cmd, {stdio: 'inherit'});
  }
}
