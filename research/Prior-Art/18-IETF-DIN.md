# IETF DIN (the IRTF decentralization work)

## What it does

The "DIN draft" in the reading list points to the IRTF DINRG, a research group under the Internet Research Task Force chartered in 2017. It carries the acronym DINRG from its original name, the Decentralized Internet Infrastructure Research Group, and its current chartered name is the Decentralization of the Internet Research Group, co-chaired by Dirk Kutscher and Lixia Zhang. Its purpose is to investigate open research issues in decentralizing infrastructure services such as trust management, identity, name resolution, and resource discovery. This entry is standards and research context, not a product that won or stalled.

## Verdict

Not a product; an active research forum with thin document output. The value to DII is vocabulary, cautionary evidence, and existing standards to align with, not a specification to implement off the shelf.

## Why

DINRG is genuinely active as a venue, meeting at essentially every IETF gathering through early 2026, but it is a research group, so it produces discussion, workshop reports, and the occasional informational document rather than binding standards. It has no adopted group drafts, and the single titled draft most associated with it, McFadden's "A Taxonomy of Internet Consolidation," is an individual draft that expired without becoming an RFC. So the specific "DIN draft" is minor and effectively dormant as a document, and DII should not build a roadmap around it.

The durable value sits in the surrounding body of work. RFC 9518, "Centralization, Decentralization, and Internet Standards" (2023, informational), is the single most useful piece: it defines centralization as the ability of a single entity or small group to observe, capture, control, or extract rent from an internet function, surveys decentralization strategies and their limits, and warns about re-centralization. The group's curated research adds hard cautionary evidence, notably the IPFS measurement study "The Cloud Strikes Back," which found a nominally decentralized network re-centralized onto a few cloud hosts, and the Named Data Networking line on signed, named objects, which is the closest architectural cousin to signed, gossiped, DNS-like discovery.

## How it compares to DII

DII plans a federation of trusted pods with signed, gossiped, DNS-like discovery and no central index, and this literature is the frame DII should be conversant in to be credible and to avoid known failure modes. Three things are worth taking directly. First, RFC 9518's definition of centralization is a usable acceptance test: for each DII function, discovery, identity, trust, and request routing, ask whether a single entity or small group could observe, control, or extract rent from it, and design out the cases where the answer is yes. Second, the re-centralization evidence is the highest-value warning, because it shows decentralized protocols routinely re-centralize at the deployment layer, on shared clouds, dominant bootstrap nodes, or a single reference implementation, which is exactly the gravity the federation batch also observed in Mastodon and Matrix. DII's "no central index" is a protocol property, and it must also watch pod hosting, seed nodes, the trust root that admits pods, and client-implementation dominance.

Third, DII should align identity and naming with existing standards rather than invent them: W3C Decentralized Identifiers give a registry-free identifier and key-rotation model well suited to pod identity, and DNSSEC and DANE offer a way to root identity in DNS if DII wants a DNS-like trust anchor. What DII cannot take off the shelf is a specification for its signed, gossiped federation protocol; that layer DII builds itself, informed by the NDN signed-object designs.

## Borrow / Avoid

Borrow: RFC 9518's centralization definition as a design acceptance test; the documented re-centralization failure modes as a checklist to design against; W3C DIDs and DNSSEC/DANE as identity and naming standards to align with; and the NDN signed-named-object work as the closest prior art for signed discovery.

Avoid: expecting an IETF or IRTF specification to implement for signed gossiped discovery, since none exists at that maturity; and building anything around the specific expired "DIN draft," which is minor and dormant.

## Sources

- [IRTF DINRG charter and status](https://datatracker.ietf.org/rg/dinrg/about/)
- [IRTF DINRG official page (chartered 2017)](https://www.irtf.org/dinrg.html)
- [RFC 9518, Centralization, Decentralization, and Internet Standards](https://www.rfc-editor.org/rfc/rfc9518.html)
- [DINRG wiki (curated research corpus)](https://wiki.ietf.org/group/dinrg)

Note: DINRG has no adopted drafts, and the McFadden consolidation-taxonomy draft's expired status is inferred from search results and the group wiki rather than a directly fetched datatracker page. RFC 9518 is an Independent Submission, informational, and not an IETF standard.
