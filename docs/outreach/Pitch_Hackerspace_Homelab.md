# For Hackerspaces, Makerspaces, and Homelabbers

Status: v0.1 (2026-07-04). One-page pitch. Background in Outreach_Guide.md; hardware in Node_Requirements.md.

## The short version

DII is a volunteer, non-commercial project, built on the Tor model, to keep a dependable floor of open-weight AI running on hardware people own, so everyday access to capable AI cannot be switched off from above. In June 2026 a frontier model was disabled worldwide by a government order for about two weeks, then turned back on. That was the proof that the AI most people lean on is permitted for now, not owned. DII is the resilient answer: local-first nodes that run open models offline and optionally federate into small trusted pods, so a community keeps a working floor no matter what happens upstream.

## Why you, first

You already do the hard part. You self-host, you run open-source infrastructure, and you probably have spare GPU capacity sitting idle right now. The culture that runs Tor relays and home servers is exactly the culture DII is built on, so there is no ideological distance to cover. A hackerspace is also the ideal shape for a pod: a small group of people and machines that already trust each other, in a physical space that can keep one box always on. Standing up the first real pod here is the fastest path from idea to something that runs.

## The ask

One anchor node. A single capable machine, hosted in the space, left on, that serves the ~30B-class open-weight floor to the pod. It can be a members' rig, a donated box, or something the space buys. Around it, any member with a decent GPU can run a lighter edge node from home, and members with no capable hardware can still use the pod as consumers. Requirements are one page, in Node_Requirements.md; the short version is a 24 GB GPU or an Apple Silicon machine with 32 GB or more, plus power and a network drop. It degrades gracefully and turns off whenever you want, and a node cut off from the network still works locally, so the space gets a private local AI box out of it even in isolation.

## What you get back

A shared, private, always-available AI node the space owns outright, with no subscription and no upstream switch. A concrete, interesting systems project to build on: the Phase-1 code is Go, everything speaks an OpenAI-compatible endpoint, and the routing and federation work is genuinely open. And the standing you would expect from being one of the first pods in a project other communities will point to.

## The honest caveats

This is early. The near-term goal is a working local node and a two-node pod that prove the floor is real, not a finished network. Throughput on one consumer GPU is a handful of tokens per second, and the workloads are chosen to tolerate that. There is no token, no marketplace, and no money in it anywhere; people run nodes because they want the floor to exist, the same reason anyone runs a Tor relay. Identity, cross-pod trust, and abuse resistance are still open design questions, and we would rather build them with a real pod than pretend they are solved.

## Next step

Ten minutes to walk a couple of members through the setup, then spin up one node on hardware you already have and see it serve. That is the whole first step.
