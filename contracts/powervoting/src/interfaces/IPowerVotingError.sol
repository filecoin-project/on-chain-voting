// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

interface IPowerVotingError {
    // time error
    error TimeError(string);

    // status error
    error StatusError(string);

    // option length error
    error OptionLengthError(string);

    // already voted error
    error AlreadyVotedError(string);

    // permission error
    error PermissionError(string);

    // address max error
    error AddressMaxError(string);
}