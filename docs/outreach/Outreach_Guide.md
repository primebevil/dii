# Outreach Guide: How to Ask Without Sounding Insane

Status: v0.1 (2026-07-04). Companion to the audience one-pagers in this folder and to Node_Requirements.md. Framing is drawn from Why DII Exists, the Architecture Overview, and the Project Charter.

This guide is for the person doing the asking, not the person being asked. It covers how to frame DII, what to request from each kind of host, the objections you will hit, and a short script you can adapt. You do not need to be a salesman. You need to sound like an engineer describing a sober, already-precedented idea, because that is what this is.

## The single most useful move

Lead with a real event and a real precedent, in that order, before you say anything about your own project. Two facts do almost all the work of making you sound credible.

The first is the June 2026 Fable 5 episode. A frontier model was switched off worldwide by government directive within hours, then switched back on about two weeks later. This is documented in mainstream press (Time, Forbes, Al Jazeera) and in Anthropic's own statement. You are not speculating about a hypothetical. You are pointing at something that already happened and asking a reasonable follow-up question: what runs when the thing you depend on is switched off from above?

The second is that the fix already has a template. Tor keeps privacy infrastructure alive through volunteers who run relays, and the Library Freedom Project got public libraries to host Tor relays starting around 2015. DII is the same shape applied to open-weight AI: volunteers and institutions run nodes so a dependable floor of capable intelligence stays available. When you say "an institution already did this for privacy, we are doing it for access to everyday AI," the idea stops sounding like a manifesto and starts sounding like the next obvious instance of a pattern people already accept.

If you open with those two points, the person across from you is oriented before you ask for anything. Everything below is detail.

## The one-sentence version

"DII is a volunteer, non-commercial project, on the Tor model, to keep a dependable floor of open-weight AI running on hardware people own, so that everyday access to capable AI cannot be switched off from above."

Say that, then stop. Let them ask the next question. The most common failure mode for a non-salesman is to keep talking past the point where the other person is already interested.

## The two-tier ask

The thing that trips people up about volunteer compute is the assumption that everyone needs a powerful machine. DII's architecture removes that assumption, and understanding it is what lets you make a small, specific ask instead of a vague large one.

The router spills work outward in one direction: local first, then to a trusted pod, then to a wider federation. A pod is a small group of nodes that already trust each other, such as a hackerspace, a lab, a library, or a family. Because of this shape you do not need every participant to own a GPU. You need one capable anchor per pod, surrounded by lighter edge nodes and by consumers who host nothing at all.

That splits your outreach into two clean asks. The anchor ask goes to whoever has real hardware or a small budget: a research-computing group with idle cluster time, a hackerspace with a rig, or a library willing to buy and host one good box. The edge ask goes to individuals, your friends-with-a-decent-GPU pattern, who add capacity but are not load-bearing. Most conversations you start will be about anchors, because the anchor is what makes a pod exist for a whole community.

Naming which ask you are making, out loud, is reassuring. "I am not asking everyone here to run this. I am asking whether this space would host one node that anchors a pod for the group" is a much easier thing to say yes to than an open-ended request for participation.

## What to actually ask for

Keep the request concrete and small, and offer a menu so the host can pick their comfort level. In rough order from lightest to heaviest:

Run an edge node on a machine you already own. This is the friends-and-members ask. Low commitment, nothing bought, turn it off whenever.

Donate idle cycles on hardware that already exists. This is the research-computing ask. Clusters and lab machines run well below capacity at night and between grant cycles, and "let a public-good floor use the idle time, within your policy" is an easy frame for a sympathetic director.

Host an always-on anchor node. This is the core institutional ask. The host provides space, power, and a network drop for one capable box that stays on and serves a pod. The box can be donated, already owned, or bought.

Fund or buy an anchor box. This is the "community data center" version, scaled down to one machine to start. A library, department, or community foundation puts money toward hardware the community then owns and runs. The precedent here is community and municipal broadband co-ops and mesh networks like Freifunk, where a group pooled money to own shared infrastructure rather than rent it.

Lend endorsement or space. Some hosts cannot run hardware but can offer a room, a mailing list, a workshop slot, or their name. For a young project, a library or university saying "we host this" is worth as much as the compute.

You are not asking for all of these at once. You are naming the menu and letting the host self-select. Most will start at the lightest option that fits them, which is fine.

## The opening script

A thirty-second version you can adapt to any of the three audiences:

"You probably saw that in June a frontier AI model got switched off worldwide by a government order for about two weeks, then switched back on. That showed the access we all lean on is something we are permitted to use for now, not something we own. There is a small volunteer project, non-commercial, modeled on Tor, that keeps a dependable floor of open-weight AI running on hardware people own, so that basic access cannot be revoked from above. Libraries already do this for privacy by running Tor relays. I am looking for a few places willing to host one node that anchors a pod for their community. Would it be worth ten minutes to walk you through what that involves?"

Adjust the middle sentence for the room. For a hackerspace, drop the explaining and lean on the Tor and self-hosting culture they already share. For a research-computing director, foreground idle-capacity donation and the open-weight, quantized workload. For a library, foreground digital equity and the Library Freedom Project precedent.

## Objections you will hit, and short answers

"Is this legal / are you helping people dodge export controls?" No. DII only uses open-weight models that are already legally published, and it explicitly does not try to serve frontier or dual-use capability. It concedes that some control over the most dangerous capabilities is legitimate. The floor is everyday capability, closer to electricity than to weapons. This is stated plainly in Why DII Exists and the "What this is not" sections of the charter.

"Is this crypto / is there a token?" No. There is no token, no tradable credit, no marketplace, and no business model. Monetization is retired, not paused. People run nodes for the same reason they run Tor relays, because they want the floor to exist. If the crypto question comes up, answer it flatly and early, because it is the fastest way to lose a mission-driven host's trust.

"What does it cost us to run?" A node is one machine, its power, and a network drop. Requirements are in Node_Requirements.md. It degrades gracefully and can be turned off at any time with no obligation, and a node that loses the network still works locally for its own users, so the host gets something even in isolation.

"Who is liable for what runs on it?" A node serves capabilities within policy the host controls: who it accepts work from, who it sends work to, and where data may travel. Cross-boundary work is opt-in and policy-gated. This is a real design area, still being worked, and you should say so rather than overclaim. Honesty about open questions reads as competence here, not weakness.

"Why not just use the cloud / a local app like Ollama?" For a single person, a local app already delivers most of the promise, and DII's own node runtime is exactly that at its core. The network adds two things a single app cannot: it serves people who have no capable device of their own, and it has no central switch anyone can throw. The point is continuity for a community, not a better single-user tool.

"Isn't this a huge undertaking?" It is deliberately modest. The near-term goal is a working local node and a two-node pod that prove the floor is real, not a global system. You are asking someone to be one of the first few pods, not to join something that must already be large to work.

## What not to say

Avoid grandiosity. Do not describe millions of strangers' machines becoming one giant computer; that is the exact framing the project's own red-team teardown rejected, and it makes you sound naive. Talk about small trusted pods.

Do not oversell the model's power. The pitch is a reliable floor of good-enough intelligence, not a frontier model in every garage. Overpromising invites an easy dismissal when the demo runs a 14B-class model instead of GPT-scale output.

Do not lead with ideology. The June 2026 event and the Tor precedent carry the argument. If you open with a speech about freedom and control, a cautious institution hears an activist. If you open with a documented event and an existing library program, they hear an engineer describing infrastructure.

Do not hide the open questions. Identity, cross-boundary verification, and abuse resistance are genuinely unsettled. Saying "here is what we have decided and here is what we are still working out" builds more trust than a frictionless story.

## Where to start

Of the three audiences, a hackerspace or makerspace is usually the easiest first yes, and a reasonable place to get your first real pod running. The people there already self-host, already share the Tor and open-source culture the pitch rests on, and often have spare hardware and a physical space that can host an always-on box. You will spend less time establishing that the idea is sane and more time on the practical setup. A working hackerspace pod then becomes the proof you point to when you walk into a research-computing office or a library, which are more cautious rooms. The three one-pagers in this folder are written to be handed to each audience directly.
