// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {Stanza} from "./age-encrypt-decrypt"

const noOpType = "no-op"

// if you wish to encrypt with AGE but simply pass the filekey in the recipient stanza, then use this
// protip: you probably don't!
class NoOpEncDec {
    static async wrap(filekey: Uint8Array): Promise<Array<Stanza>> {
        return [{
            type: noOpType,
            args: [],
            body: filekey
        }]
    }

    static async unwrap(recipients: Array<Stanza>): Promise<Uint8Array> {
        if (recipients.length !== 1) {
            throw Error("NoOpEncDec only expects a single stanza!")
        }

        if (recipients[0].type !== noOpType) {
            throw Error(`NoOpEncDec expects the type of the stanza to be ${noOpType}`)
        }

        return recipients[0].body
    }
}

export {NoOpEncDec}
