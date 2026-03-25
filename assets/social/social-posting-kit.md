# Oracle Campaign Social Posting Kit

This kit is for the **native 9:16 cuts**.

Use the vertical masters in `assets/renders/vertical/` and the matching posters in `assets/renders/posters-vertical/`.

## Recommended publishing order

Start with the easiest proof, then widen the frame:

1. `OracleVertical01MemorySurvives`
2. `OracleVertical05OneLineSetup`
3. `OracleVertical06LocalWith26ai`
4. `OracleVertical04VectorSearch`
5. `OracleVertical02Routing`
6. `OracleVertical03SystemOfRecord`
7. `OracleVertical07BridgeToADB`
8. `OracleVertical08CleanerArchitecture`
9. `OracleVertical09ProductionShaped`
10. `OracleVertical10HeroFilm`

---

## 01. Memory survives restart

**Video:** `assets/renders/vertical/oracle-vertical-01-memory-survives.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-01-memory-survives.png`

**Primary title**  
Your agent kept its memory after the restart

**Alt titles**
- Oracle-backed agents don’t lose context
- Durable agent memory, finally
- Restart safe memory for Deep Agents

**Thumbnail headline**  
Memory survives restart

**Caption**  
Most agent demos forget everything when the process dies.

This fork pushes `/memory`, `/history`, and `/skills` into Oracle AI Database, so the next run keeps going from stored context instead of starting cold.

**Short CTA**  
Deep Agents + Oracle AI Database.

**Hashtags**  
#AI #Agents #Oracle #Python #Remotion

---

## 02. Where the memory actually lives

**Video:** `assets/renders/vertical/oracle-vertical-02-routing.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-02-routing.png`

**Primary title**  
Persistence is a routing decision

**Alt titles**
- The cleanest part of this Oracle fork
- How we routed durable state into Oracle
- CompositeBackend, but actually useful

**Thumbnail headline**  
Route memory on purpose

**Caption**  
The Oracle story works because the persistence split is explicit.

`/memory`, `/history`, and `/skills` go durable. Scratch work stays fast and local. That keeps the architecture clean and keeps Oracle central without dragging every temp path through the database.

**Short CTA**  
Additive architecture. Better memory.

**Hashtags**  
#AIEngineering #Agents #Oracle #Python

---

## 03. Oracle as the system of record

**Video:** `assets/renders/vertical/oracle-vertical-03-system-of-record.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-03-system-of-record.png`

**Primary title**  
Agent memory needs a real storage model

**Alt titles**
- Oracle as the system of record for agent memory
- Durable memory needs more than vibes
- The `da_files` table makes this story real

**Thumbnail headline**  
Real storage model

**Caption**  
This memory layer lands because it has a concrete shape.

Writes resolve into Oracle SQL. The state lands in `da_files`. You can point at rows, timestamps, and namespaces instead of treating memory like some spooky side cache.

**Short CTA**  
Memory with a real spine.

**Hashtags**  
#Database #Oracle #AIAgents #Python

---

## 04. Semantic memory with Oracle AI Vector Search

**Video:** `assets/renders/vertical/oracle-vertical-04-vector-search.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-04-vector-search.png`

**Primary title**  
Semantic memory inside Oracle AI Database

**Alt titles**
- Keyword search is a weak memory model
- Oracle AI Vector Search changes the ceiling
- Meaning-aware retrieval for agent memory

**Thumbnail headline**  
Semantic memory

**Caption**  
Keyword search misses too much.

The vector path in this fork stores embeddings on write, then uses `VECTOR_DISTANCE` to pull memory back by meaning. Same Oracle-centered architecture, much smarter recall.

**Short CTA**  
Smarter memory, no extra vector sidecar.

**Hashtags**  
#VectorSearch #Oracle #AI #Agents #RAG

---

## 05. One line to get an Oracle-backed agent

**Video:** `assets/renders/vertical/oracle-vertical-05-one-line-setup.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-05-one-line-setup.png`

**Primary title**  
One line gets you an Oracle-backed agent

**Alt titles**
- Fast start, durable memory
- The shortest path into the Oracle fork
- Good defaults matter more than speeches

**Thumbnail headline**  
One-line setup

**Caption**  
The storage story is serious. The setup still stays short.

`create_oracle_deep_agent()` wires the Oracle config, connection manager, backend routing, and memory defaults together so the first success happens fast.

**Short CTA**  
Fast start. Durable context.

**Hashtags**  
#Python #Agents #Oracle #DeveloperExperience

---

## 06. Oracle Database 26ai Free on your laptop

**Video:** `assets/renders/vertical/oracle-vertical-06-local-with-26ai.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-06-local-with-26ai.png`

**Primary title**  
This Oracle story starts on a laptop

**Alt titles**
- Oracle Database 26ai Free for local agent dev
- Local-first Oracle memory for Deep Agents
- Docker, Oracle, and a real memory loop

**Thumbnail headline**  
Oracle on your laptop

**Caption**  
The nicest part of this fork is that it doesn’t start in a cloud slide deck.

A Docker profile brings up Oracle Database 26ai Free locally, then one setup script creates the user and schema so the memory path is live right away.

**Short CTA**  
Local first. Oracle-backed.

**Hashtags**  
#Oracle #Docker #Python #AIAgents

---

## 07. The bridge to Autonomous Database

**Video:** `assets/renders/vertical/oracle-vertical-07-bridge-to-adb.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-07-bridge-to-adb.png`

**Primary title**  
The move to cloud should feel boring

**Alt titles**
- From local Oracle to ADB without changing the model
- Same agent shape, bigger runway
- The clean bridge to Autonomous Database

**Thumbnail headline**  
Local to ADB

**Caption**  
The best deployment bridges feel uneventful.

Local mode uses host, port, and service. Cloud mode swaps to DSN and wallet-aware config. The agent factory and backend wiring still look the same, which is exactly what you want.

**Short CTA**  
Same model. Bigger runway.

**Hashtags**  
#Oracle #ADB #Cloud #Agents #Python

---

## 08. Fewer moving parts, cleaner architecture

**Video:** `assets/renders/vertical/oracle-vertical-08-cleaner-architecture.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-08-cleaner-architecture.png`

**Primary title**  
A cleaner memory story starts with cleaner data architecture

**Alt titles**
- Fewer moving parts for agent persistence
- Oracle gives the stack a center of gravity
- Less glue, stronger memory layer

**Thumbnail headline**  
Less glue

**Caption**  
Agent stacks sprawl fast.

This fork pulls persistent memory, semantic retrieval, and durable state toward one stronger center in Oracle AI Database. Fewer seams, less glue, easier to explain.

**Short CTA**  
Cleaner stack. Better memory.

**Hashtags**  
#Architecture #Oracle #AIAgents #PlatformEngineering

---

## 09. Why Oracle changes the Deep Agents story

**Video:** `assets/renders/vertical/oracle-vertical-09-production-shaped.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-09-production-shaped.png`

**Primary title**  
Oracle gives Deep Agents a more durable second act

**Alt titles**
- The part Oracle hardens the most
- Deep Agents was already strong, Oracle changes the feel
- The memory layer is what makes this stack grow up

**Thumbnail headline**  
A more durable stack

**Caption**  
Deep Agents already gives you planning, tools, sub-agents, and context handling.

Oracle hardens the layer people usually trust the least: durable memory, continuity across sessions, stronger retrieval, and a cleaner path from local to cloud.

**Short CTA**  
Stronger harness. Stronger memory layer.

**Hashtags**  
#AI #Agents #Oracle #Python #Platform

---

## 10. Hero film

**Video:** `assets/renders/vertical/oracle-vertical-10-hero-film.mp4`  
**Poster:** `assets/renders/posters-vertical/oracle-vertical-10-hero-film.png`

**Primary title**  
Deep Agents + Oracle AI Database

**Alt titles**
- Durable memory for agent systems
- Oracle gives Deep Agents staying power
- The full Oracle story in 16 seconds

**Thumbnail headline**  
Staying power

**Caption**  
Deep Agents gives you the harness.

Oracle AI Database gives it staying power: durable memory, semantic recall, and a path from local 26ai Free to Autonomous Database that actually lines up.

**Short CTA**  
Built around Oracle.

**Hashtags**  
#Oracle #AI #Agents #VectorSearch #Python
