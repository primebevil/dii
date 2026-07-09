# ADR-0010: Voluntary Sponsorship, Not Paid Private Use

Status: Accepted, 2026-07-06. Extends ADR-0008 and ADR-0009, and resolves the private-pod sponsorship question raised in the batch-4 design review.

## Context

ADR-0008 opened funding into a non-commercial steward, and ADR-0009 established that private pods serve only their own members and self-fund. The review then raised a fair question: a private pod, a law firm, a medical-research group, a corporate team, gets real value from the software and the network DII built, even though it shares only within itself. Could DII ask such pods for money to fund the public open resources, and does that break the non-commercial stance?

The framing that resolves it: the network stays free and ungated regardless, because a private pod runs its own hardware in its own closed zone and consumes no public capacity, so any contribution is for the software and the mission, not for access to the public floor. The remaining question is only whether the contribution is asked for or required.

## Decision

DII invites voluntary sponsorship and donation to fund public-serving pods and the steward. It does not charge a required fee or license for private or commercial use. The voluntary contribution may be socially normed, a "suggested sustaining contribution" for well-resourced or commercial private pods, but it is never enforced and never gates use. The justification is a cross-subsidy for the mission: private beneficiaries funding public access.

This keeps ADR-0002 and ADR-0005 intact. The network is still free, nothing is sold, and no one is charged to use it; money flows into the non-commercial steward, which passes the direction-and-destination test in ADR-0008. It sits on the funding ladder from ADR-0009 as a voluntary top-up: private pods self-fund their own operation and are additionally invited, not required, to sponsor the public good. Because no one becomes a paying customer, it avoids the mission-pull that the review flagged as the main risk of charging.

The alternative, a required fee or license for private or commercial use, is explicitly not adopted and is kept in reserve. It would not violate "the network is free," since it charges for the software rather than for public access, but it introduces paying customers and their pull, enforcement needs such as a non-commercial or source-available license that would turn "free software for all" into "free for some," and added complexity. It is revisited only if voluntary funding does not sustain, and would be its own decision.

## Consequences

This gives the steward a concrete, on-mission ask to make of institutional and commercial private-pod operators, and it is expected to capture most of the realistic funding from the parties best able to pay without the costs of a paid tier. It starts the funding conversation from the safest end, consistent with the whole project's stance of choosing the version that does not compromise the mission and reserving the harder version for later if needed.

## Open questions

The mechanics of the donation and sponsorship channel remain open and tie to ADR-0008's individual-donation on-ramp and the still-unchosen legal vehicle that can receive the money. Whether to publish suggested contribution levels, and how to frame the ask so it reads as mission support rather than a disguised price, is left for when the steward exists.
