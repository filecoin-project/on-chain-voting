# On Chain Voting for FIPs
This repo contains a complete copy of the on chain voting tool being developed
by StorSwift for use in the Filecoin Improvement Proposal (FIP) process. More
information about FIPs can be found in [the FIPs
repo](https://github.com/filecoin-project/FIPs)

### Specification
The specification for the on chain voting tool can be found [in this Google Doc](https://docs.google.com/document/d/13910NE-O3mUQ6rztt6f3xe7hwW_aS-xaPW_zHuTpBW4/edit)

### Components
The current implementation is spread across multiple repositories. This repo is
intended to serve as a single point of reference for the implementation, and
will include technical documentation and details around deployment and testing,
as well as importing all components as submodules to ensure a copy of the code
remains in the `filecoin-project` org.

The component repositories are:
- [Power Voting](https://github.com/black-domain/power-voting)
- [Power Voting back-end](https://github.com/black-domain/powervoting-backend)
- [Power Voting contracts](https://github.com/black-domain/powervoting-contracts)
- [Oracle contracts](https://github.com/black-domain/power-oracle-contracts)
- [Power Oracle node](https://github.com/black-domain/power-oracle-node)
- [UCAN utils](https://github.com/black-domain/ucan-utils)
