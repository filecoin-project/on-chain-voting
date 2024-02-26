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

import {hkdf} from "@noble/hashes/hkdf"
import {sha256} from "@noble/hashes/sha256"
import {hmac} from "@noble/hashes/hmac"

export function createMacKey(fileKey: Uint8Array, macMessage: string, headerText: string): Uint8Array {
    // empty string salt as per the spec!
    const hmacKey = hkdf(sha256, fileKey, "", Buffer.from(macMessage, "utf8"), 32)
    return Buffer.from(hmac(sha256, hmacKey, Buffer.from(headerText, "utf8")))
}

// returns a string of n bytes read from a CSPRNG like /dev/urandom.
export async function random(n: number): Promise<Uint8Array> {
    if (typeof window === "object" && "crypto" in window) {
        return window.crypto.getRandomValues(new Uint8Array(n))
    }

    // parcel likes to resolve polyfills for things even if they aren't used
    // so this indirection tricks it into not doing it and not complaining :)
    const x = "crypto"
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    const bytes = require(x).randomBytes(n)
    return new Uint8Array(bytes.buffer, bytes.byteOffset, bytes.byteLength)
}
