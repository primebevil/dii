# SETI@home

## What it does

SETI@home was a volunteer computing project from UC Berkeley, public from May 1999, that had people donate idle CPU cycles through a screensaver to scan radio-telescope data for narrowband signals that might indicate extraterrestrial intelligence. Data came piggyback from Arecibo and later Green Bank. It was not the first internet volunteer-computing project, but it is the one that popularized the model, and it later spun out BOINC as reusable middleware. Its two goals were to do real SETI science and to prove volunteer computing works; the second succeeded completely, the first found nothing confirmed.

## Verdict

Won the proof of concept, stalled on durability. A landmark that mobilized millions and then stopped distributing work in March 2020, not a failure but not perpetual infrastructure either.

## Why

The scale was real. Over 5.2 million participants took part across its lifetime, the most of any volunteer-computing project at the time. Guinness recorded it as the largest computation in history, and by mid-2013 it was running at over 668 teraflops across roughly 145,000 active machines. The active peak, reported around a million volunteers, tracked the roughly one year of heavy press after launch, then declined gradually as coverage faded. By the March 2020 hibernation the project was down to about 91,000 active users on 145,000 hosts.

It ended for two reasons the project stated plainly: it had analyzed as much data as it judged useful for now, a demand-side limit, and managing the distributed processing was a lot of work for a small team that wanted to move to back-end analysis. Thin funding was a chronic background pressure, there was no government money for the science, and Arecibo's budget was shrinking before the dish collapsed in December 2020. The deeper point is that SETI@home was built as a science experiment with a finite question, not as durable infrastructure, and had no plan for the day the useful work ran out.

## How it compares to DII

The gap that matters is the cost of the donation. SETI@home reached millions precisely because giving was almost free: spare CPU cycles otherwise wasted on a screensaver, no felt electricity cost, no foregone use, no opportunity cost. Even at that near-zero burden, sustained participation still eroded to under 100,000 active over two decades. DII asks for something categorically heavier, expensive and rivalrous GPU inference, where the card cannot serve DII and do the owner's own work at once and the electricity is real. So SETI@home's million-volunteer peak is not a reachable analog for DII's ceiling; that scale was subsidized by a burden DII does not have.

What DII offers that SETI@home did not is a renewing reason to exist. SETI@home ran out of useful work; DII's work is inference demand from real users, which does not deplete the way a fixed telescope backlog does, and DII's benefit accrues to a live user constituency in the present rather than to a someday-detection. DII also spreads stewardship across federated pods rather than resting continuity on one lab's staffing and one telescope, which is the single-institution fragility that ended SETI@home.

## Borrow / Avoid

Borrow: mission as the primary recruiter, since a legible mission moved a million people; credit, competition, and team identity, which are what retained volunteers after novelty faded and which matter more, not less, when the donation has real cost; cross-validation by sending each unit to more than one node, which deterred cheating and is directly relevant to any contribution or reputation credit; and spinning the reusable platform out from the mission, since BOINC, not the SETI campaign, is what endured.

Avoid: single-institution dependence for funding, operations, and work supply; running out of useful work with no durability plan; and building the sustainability model on press attention, which is a spike that decays. Above all, do not read SETI@home's scale as achievable at DII's burden.

## Sources

- [SETI@home, Wikipedia](https://en.wikipedia.org/wiki/SETI@home)
- [SETI@home going into hibernation, SETI Institute (Mar 2020)](https://www.seti.org/news/seti-at-home-going-into-hibernation/)
- [SETI@home stops after 20 years, Information Age / ACS](https://ia.acs.org.au/article/2020/seti-home-stops-after-20-years.html)

Note: the roughly one-million active-peak figure comes from summaries of the Anderson and Korpela Annual Reviews account, not a primary page fetched in full; the 5.2M lifetime figure, the 2020 hibernation, and the stated reasons are from the sources above.
