# On Chain Voting for FIPs
This repo contains a complete copy of the on chain voting tool being developed
by the Filecoin Foundation, in partnership with StorSwift, for use in the Filecoin Improvement Proposal (FIP) process. More
information about FIPs can be found in [the FIPs
repo](https://github.com/filecoin-project/FIPs)

If you'd like to try out a deployed version of the tool, you can check it out at [vote.fil.org](https://vote.fil.org).

Please keep in mind, the tool is still under active development, and may not currently meet all the requirements of our community. If you find someplace where we are falling short (i.e. buggy code, clunky UI, unclear process), please submit an issue to this repo and bring it to our attention.

### Specification
The specification for the on chain voting tool can be found [in this Google Doc](https://docs.google.com/document/d/13910NE-O3mUQ6rztt6f3xe7hwW_aS-xaPW_zHuTpBW4/edit)

### Repo structure
The system at a high level is composed of a few core components, each listed in a separate folder:

- [**docs**](https://github.com/filecoin-project/on-chain-voting/tree/main/docs): Documentation
- [**snapshot**](https://github.com/filecoin-project/on-chain-voting/tree/main/snapshot): Contains code for the snapshot service, used to take and store snapshots of key data for later use in calculating vote power.
- [**frontend**](https://github.com/filecoin-project/on-chain-voting/tree/main/frontend): Contains the frontend dApp code that runs the UI, timelock encypts vote data, and stores proposals and ecnypted data in web3.storage.
- [**backend**](https://github.com/filecoin-project/on-chain-voting/tree/main/backend): The powervoting backend contains code necessary to sync proposals and votes, decode all votes once the timelock encryption has expired, and calculate the final result based on the power of each respective vote. Everything this service does should be independently verifiable using on chain data.
- [**contracts**](https://github.com/filecoin-project/on-chain-voting/tree/main/contracts): Contains smart contracts used for the core Power Voting functionality, including managing FIP Editors, creating proposals, and casting votes.

### Testing
Each component of the overall system is tested using a go test suite (or [`forge`](https://github.com/foundry-rs/foundry), in the case of smart contracts). The test suite for each component is run by CI anytime the code in that subfolder changes. You can view test history on theGithub actions page: https://github.com/filecoin-project/on-chain-voting/actions

### Community Code Review
*“Given enough eyeballs, all bugs are shallow” - Linus’s Law*

We are requesting the community help us to engage in a community code review of the current implementation. This is to make sure of a few things:
- There are no errors with the specification
- The current implementation (and associated tests) agree with the specification
- The current test suite is testing the code correctly
- There are no critical pathways that are not being tested)
- The current tool is usable by all members of the community
- The current tool is not biased towards any one particular part of the community

To share your feedback, please file an issue on this repository, or reach out to ian@fil.org directly with more substaintial questions.
