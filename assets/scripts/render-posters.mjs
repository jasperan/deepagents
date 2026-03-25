import {execSync} from 'node:child_process';
import {mkdirSync} from 'node:fs';

const posters = [
  ['Oracle01MemorySurvives', 210, 'oracle-01-memory-survives.png'],
  ['Oracle02Routing', 390, 'oracle-02-routing.png'],
  ['Oracle03SystemOfRecord', 390, 'oracle-03-system-of-record.png'],
  ['Oracle04VectorSearch', 210, 'oracle-04-vector-search.png'],
  ['Oracle05OneLineSetup', 210, 'oracle-05-one-line-setup.png'],
  ['Oracle06LocalWith26ai', 210, 'oracle-06-local-with-26ai.png'],
  ['Oracle07BridgeToADB', 390, 'oracle-07-bridge-to-adb.png'],
  ['Oracle08CleanerArchitecture', 390, 'oracle-08-cleaner-architecture.png'],
  ['Oracle09ProductionShaped', 390, 'oracle-09-production-shaped.png'],
  ['Oracle10HeroFilm', 610, 'oracle-10-hero-film.png'],
];

const dryRun = process.argv.includes('--dry-run');
mkdirSync('renders/posters', {recursive: true});

for (const [id, frame, output] of posters) {
  const cmd = `npx remotion still src/index.ts ${id} renders/posters/${output} --frame=${frame}`;
  console.log(cmd);
  if (!dryRun) {
    execSync(cmd, {stdio: 'inherit'});
  }
}
