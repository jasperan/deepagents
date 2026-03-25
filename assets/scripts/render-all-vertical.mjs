import {execSync} from 'node:child_process';
import {mkdirSync} from 'node:fs';

const renders = [
  ['OracleVertical01MemorySurvives', 'oracle-vertical-01-memory-survives.mp4'],
  ['OracleVertical02Routing', 'oracle-vertical-02-routing.mp4'],
  ['OracleVertical03SystemOfRecord', 'oracle-vertical-03-system-of-record.mp4'],
  ['OracleVertical04VectorSearch', 'oracle-vertical-04-vector-search.mp4'],
  ['OracleVertical05OneLineSetup', 'oracle-vertical-05-one-line-setup.mp4'],
  ['OracleVertical06LocalWith26ai', 'oracle-vertical-06-local-with-26ai.mp4'],
  ['OracleVertical07BridgeToADB', 'oracle-vertical-07-bridge-to-adb.mp4'],
  ['OracleVertical08CleanerArchitecture', 'oracle-vertical-08-cleaner-architecture.mp4'],
  ['OracleVertical09ProductionShaped', 'oracle-vertical-09-production-shaped.mp4'],
  ['OracleVertical10HeroFilm', 'oracle-vertical-10-hero-film.mp4'],
];

const dryRun = process.argv.includes('--dry-run');
mkdirSync('renders/vertical', {recursive: true});

for (const [id, output] of renders) {
  const cmd = `npx remotion render src/index.ts ${id} renders/vertical/${output}`;
  console.log(cmd);
  if (!dryRun) {
    execSync(cmd, {stdio: 'inherit'});
  }
}
