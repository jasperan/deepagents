# Oracle Campaign Voiceover Guide

This pass adds **voiceover timing and recording scripts**, not generated narration.

`ELEVENLABS_API_KEY` is not set in the environment, so no ElevenLabs TTS audio was created. The timing below is ready for either human recording or future TTS generation.

## Timing model

Default 22s videos:
- 0.0s to 4.4s: hook
- 4.0s to 9.2s: proof
- 9.2s to 14.4s: architecture
- 14.4s to 18.4s: value
- 18.4s to 22.0s: close

Hero film, 24s:
- 0.0s to 4.4s: hook
- 4.0s to 9.2s: proof
- 9.2s to 14.4s: architecture
- 14.4s to 18.4s: value
- 18.4s to 24.0s: close

---

## Oracle01MemorySurvives

**0.0s to 4.4s**  
Most agents forget as soon as the process dies.

**4.0s to 9.2s**  
Here, the agent saves its findings to Oracle, shuts down, restarts, and keeps going from the same memory.

**9.2s to 14.4s**  
Memory, history, and skills all land in the `da_files` table as durable state.

**14.4s to 18.4s**  
That turns a fragile demo loop into real cross-session continuity.

**18.4s to 22.0s**  
Kill the process. Keep the memory.

---

## Oracle02Routing

**0.0s to 4.4s**  
The Oracle story works because persistence is routed cleanly.

**4.0s to 9.2s**  
`/memory`, `/history`, and `/skills` go to Oracle. Scratch paths stay in memory.

**9.2s to 14.4s**  
That split lives inside `CompositeBackend`, so the integration stays additive to upstream Deep Agents.

**14.4s to 18.4s**  
You get durable state where it matters without dragging temporary work through the database.

**18.4s to 22.0s**  
Persistence becomes a routing decision, not a rewrite.

---

## Oracle03SystemOfRecord

**0.0s to 4.4s**  
This memory layer is persuasive because it has a concrete storage model.

**4.0s to 9.2s**  
Writes resolve into Oracle SQL, with MERGE-based upserts and timestamps on every update.

**9.2s to 14.4s**  
The `da_files` schema gives memory a durable address, a namespace, and a traceable history.

**14.4s to 18.4s**  
That means agent state stops feeling like a cache and starts feeling like application data.

**18.4s to 22.0s**  
Oracle becomes the system of record.

---

## Oracle04VectorSearch

**0.0s to 4.4s**  
Keyword search is a weak memory model for agents.

**4.0s to 9.2s**  
`OracleVectorBackend` stores embeddings on write and queries them with `VECTOR_DISTANCE`.

**9.2s to 14.4s**  
Now the agent can retrieve related memory by meaning, not just by exact phrasing.

**14.4s to 18.4s**  
And it all stays inside Oracle AI Database instead of spilling into another vector service.

**18.4s to 22.0s**  
That is where Oracle changes the ceiling of the product.

---

## Oracle05OneLineSetup

**0.0s to 4.4s**  
Power matters more when the setup stays short.

**4.0s to 9.2s**  
`create_oracle_deep_agent()` gives you an Oracle-backed Deep Agent in one clean entrypoint.

**9.2s to 14.4s**  
The helper wires config, connection management, backend routing, memory, and skills for you.

**14.4s to 18.4s**  
That makes the first success fast without hiding the architecture from serious users.

**18.4s to 22.0s**  
A stronger storage story lands better when the first line stays simple.

---

## Oracle06LocalWith26ai

**0.0s to 4.4s**  
The Oracle story starts locally.

**4.0s to 9.2s**  
A Docker profile brings up Oracle Database 26ai Free, and one setup script initializes the user and schema.

**9.2s to 14.4s**  
That gives you a real persistence loop on a laptop, not just a cloud-shaped promise.

**14.4s to 18.4s**  
It lowers the barrier to trying the fork and makes demos much more honest.

**18.4s to 22.0s**  
Oracle feels a lot closer when it runs right next to the app.

---

## Oracle07BridgeToADB

**0.0s to 4.4s**  
The bridge to cloud should feel like a deployment step, not a rewrite.

**4.0s to 9.2s**  
Here, local mode uses host, port, and service. Cloud mode switches to DSN and wallet-aware config.

**9.2s to 14.4s**  
The agent factory and backend wiring still look the same, which keeps the mental model stable.

**14.4s to 18.4s**  
That makes the move to Autonomous Database easier to explain and easier to trust.

**18.4s to 22.0s**  
Same architecture. Bigger runway.

---

## Oracle08CleanerArchitecture

**0.0s to 4.4s**  
A lot of agent stacks sprawl into too many services.

**4.0s to 9.2s**  
This fork pulls persistent memory, semantic retrieval, and operational glue toward one durable center.

**9.2s to 14.4s**  
Oracle AI Database becomes the data plane instead of one more box in a crowded diagram.

**14.4s to 18.4s**  
That means fewer seams, less glue, and a story teams can explain more easily.

**18.4s to 22.0s**  
Cleaner memory usually starts with cleaner data architecture.

---

## Oracle09ProductionShaped

**0.0s to 4.4s**  
Deep Agents already gives you a strong harness.

**4.0s to 9.2s**  
Planning, tools, sub-agents, and context management are already there.

**9.2s to 14.4s**  
Oracle hardens the layer agent demos usually undersell: durable memory, better retrieval, and cloud continuity.

**14.4s to 18.4s**  
That changes how serious the whole stack feels in front of builders and platform teams.

**18.4s to 22.0s**  
Oracle gives Deep Agents a more durable second act.

---

## Oracle10HeroFilm

**0.0s to 4.4s**  
Deep Agents gives you the harness.

**4.0s to 9.2s**  
Oracle AI Database gives it durable memory, cross-session continuity, and semantic retrieval.

**9.2s to 14.4s**  
The same story runs from local 26ai Free to Autonomous Database in the cloud.

**14.4s to 18.4s**  
That makes the product easier to demo, easier to explain, and easier to take seriously.

**18.4s to 24.0s**  
Deep Agents gets the harness. Oracle gives it staying power.
