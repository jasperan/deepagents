import {execSync} from 'node:child_process';
import {mkdirSync} from 'node:fs';

const posters = [
  ['OracleVertical01MemorySurvives', 150, 'oracle-vertical-01-memory-survives.png'],
  ['OracleVertical02Routing', 180, 'oracle-vertical-02-routing.png'],
  ['OracleVertical03SystemOfRecord', 210, 'oracle-vertical-03-system-of-record.png'],
  ['OracleVertical04VectorSearch', 240, 'oracle-vertical-04-vector-search.png'],
  ['OracleVertical05OneLineSetup', 150, 'oracle-vertical-05-one-line-setup.png'],
  ['OracleVertical06LocalWith26ai', 180, 'oracle-vertical-06-local-with-26ai.png'],
  ['OracleVertical07BridgeToADB', 210, 'oracle-vertical-07-bridge-to-adb.png'],
  ['OracleVertical08CleanerArchitecture', 210, 'oracle-vertical-08-cleaner-architecture.png'],
  ['OracleVertical09ProductionShaped', 225, 'oracle-vertical-09-production-shaped.png'],
  ['OracleVertical10HeroFilm', 270, 'oracle-vertical-10-hero-film.png'],
];

const dryRun = process.argv.includes('--dry-run');
mkdirSync('renders/posters-vertical', {recursive: true});

for (const [id, frame, output] of posters) {
  const cmd = `npx remotion still src/index.ts ${id} renders/posters-vertical/${output} --frame=${frame}`;
  console.log(cmd);
  if (!dryRun) {
    execSync(cmd, {stdio: 'inherit'});
  }
}
