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

// as per the spec:
// RFC 4648, Section 4
// without = padding characters (sometimes referred to as "raw" or "unpadded" base64)
export function unpaddedBase64(buf: Uint8Array | string): string {
    const encodedBuf = Buffer.from(buf).toString("base64")

    let lastIndex = encodedBuf.length - 1

    while (encodedBuf[lastIndex] === "=") {
        lastIndex--
    }

    return encodedBuf.slice(0, lastIndex + 1)
}

export function unpaddedBase64Buffer(buf: Uint8Array | string): Buffer {
    return Buffer.from(unpaddedBase64(buf), "base64")
}

/*
    e.g. chunked("hello world", 2, ".") returns
    ["he.", "ll.", "o .", "wo.", "rl.", "d."]
 */
export function chunked(input: string, chunkSize: number, suffix = ""): Array<string> {
    const output = []
    let currentChunk = ""
    for (let i = 0, chunks = 0; i < input.length; i++) {
        currentChunk += input[i]

        const posInChunk = i - (chunks * chunkSize)

        if (posInChunk === chunkSize - 1) {
            output.push(currentChunk + suffix)
            currentChunk = ""
            chunks++
        } else if (i === input.length - 1) {
            output.push(currentChunk + suffix)
        }
    }

    return output
}

// slices the input string up to and including the first
// occurrence of the string provided in `searchTerm`
// returns the whole string if it's not found
// e.g. sliceUntil("hello world", "ll") will return "hell"
export function sliceUntil(input: string, searchTerm: string) {
    let lettersMatched = 0
    let inputPointer = 0

    while (inputPointer < input.length && lettersMatched < searchTerm.length) {
        if (input[inputPointer] === searchTerm[lettersMatched]) {
            ++lettersMatched
        } else if (input[inputPointer] === searchTerm[0]) {
            lettersMatched = 1
        } else {
            lettersMatched = 0
        }

        ++inputPointer
    }

    return input.slice(0, inputPointer)
}

