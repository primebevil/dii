# Folding@home

## What it does

Folding@home is a distributed-computing project, founded at Stanford in 2000 by Vijay Pande and now run by Greg Bowman's lab at the University of Pennsylvania as a multi-lab consortium, where volunteers donate GPU and CPU compute to simulate protein folding and molecular dynamics for disease research. It uses its own client and servers rather than BOINC. It is the closest precedent to DII in this batch because it mobilized GPU donors specifically, the same expensive, rivalrous compute DII needs, rather than the free CPU cycles most volunteer projects rely on.

## Verdict

Won on scientific mission and historic scale, stalled on sustaining the surge. A real, permanent success that also demonstrates how quickly crisis-driven participation recedes.

## Why

The COVID spike is the story. Folding@home entered 2020 as a mature project of about 30,000 active devices. After it announced a SARS-CoV-2 pivot in early March, more than 400,000 new devices joined within two weeks, and by March 25 it reached 768 native petaflops, about 1.5 x86 exaflops, making it the first exaflop-scale computing system of any kind, briefly faster than the top 500 supercomputers combined. Press put the mid-April aggregate near 2.4 exaflops, on more than 280,000 GPUs and 4.8 million CPU cores, with over a million active devices and four million total participants by June 2020. The science was real, including cryptic drug-binding pockets on the viral proteome and antiviral leads that reached animal testing.

Then it receded. Active participation fell within months of the peak even as cumulative sign-ups stayed high, and by late October 2025 throughput was down to 12.9 native petaflops, roughly 98 percent below peak and near or below where it started before COVID. The surge cohort was not retained. The durable baseline is the same population that was always there, hardware enthusiasts and long-time altruists, many with personal disease connections, a committed core in the tens of thousands. The lesson is that acute-crisis motivation and durable everyday motivation are different fuels: the first produces a spike, only the second sets the floor, and here the floor never ratcheted up.

## How it compares to DII

Folding@home is the most direct evidence in the batch for DII's central question, because it tested GPU donation at scale. The good news is that the demand-side potential is real: millions of GPU owners can be mobilized for a legible mission. The hard news is retention. GPU donation carries a real cost in electricity, wear, and the opportunity cost of the card, and mission goodwill spiked on crisis but did not durably outbid that cost for most people, so almost none of the surge stayed. DII should read its own participation the same way: plan for the committed core as the baseline and treat any surge as a windfall to capture, never as the operating assumption.

What DII offers that a pure crisis mobilization does not is a reason to stay once urgency fades, because a DII node delivers a standing benefit to its own operator and pod, a reliable floor of intelligence they use every day, rather than only contributing to a distant collective goal. That everyday, self-interested utility is the durable fuel Folding@home mostly lacked outside the emergency. The caution DII must hold is not to budget on spike-level participation, and to measure active contributors rather than cumulative sign-ups, which is exactly where Folding@home's numbers looked healthier than the reality.

## Borrow / Avoid

Borrow: mission salience under a GPU-relevant framing, which is the strongest lever observed here; a pre-built amplification and scaling stack so a surge can be captured rather than choke the servers, as nearly happened in 2020; deliberate cultivation of the enthusiast and overclocker core, which is the twenty-year durable asset; and shipping verifiable results to keep the mission legible.

Avoid: budgeting on a spike, since Folding@home's receded about 98 percent; assuming mission alone retains GPU donors against a real cost, since it did not; and reporting cumulative sign-ups as if they were active donors.

## Sources

- [Folding@home, Wikipedia](https://en.wikipedia.org/wiki/Folding@home)
- [Folding@home 2020 in review, Greg Bowman](https://foldingathome.org/2021/01/05/2020-in-review-and-happy-new-year-2021/)
- [Over 4 million computers joined Folding@home, HPCwire](https://www.hpcwire.com/off-the-wire/over-4-million-computers-worldwide-joined-foldinghome-to-aid-in-coronavirus-research/)

Note: the roughly 2.4 exaflops peak is an x86-equivalent press figure; the project's own precisely dated numbers are 470 native petaflops on March 20 and 768 native petaflops on March 25, 2020. A clean official mid-2026 active-device count could not be verified from the JavaScript-rendered stats page; the 12.9 native petaflops figure is the project's own metric as of October 31, 2025.
