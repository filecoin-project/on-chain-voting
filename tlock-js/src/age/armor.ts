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

import {chunked} from "./utils"

const header = "-----BEGIN AGE ENCRYPTED FILE-----"
const footer = "-----END AGE ENCRYPTED FILE-----"

// takes some payload and encodes it as armor with the AGE armor headers in lines of size `chunkSize`
export function encodeArmor(input: string, chunkSize = 64): string {
    const base64Input = Buffer.from(input, "binary").toString("base64")
    const columnisedInput = chunked(base64Input, chunkSize).join("\n")

    // if the last line is exactly 64 columns, add an extra newline
    let paddedFooter = footer
    if (columnisedInput.length > 0 && columnisedInput[columnisedInput.length - 1].length === 64) {
        paddedFooter = "\n" + footer
    }
    return `${header}\n${columnisedInput}\n${paddedFooter}\n`
}

// takes an armored payload and decodes it if it is an AGE armor payload
// and it satisfies some security properties
export function decodeArmor(armor: string, chunkSize = 64): string {
    // could start/end with space or newlines, let's strip them
    armor = armor.trimStart()
    const lengthBeforeEndTrim = armor.length
    armor = armor.trimEnd()
    const lengthAfterTrim = armor.length

    // for compliance with the go age implementation, we deny more than 1024 whitespace chars:
    // see: https://github.com/FiloSottile/age/blob/8e3f74c283b2e9b3cd0ec661fa4008504e536d20/armor/armor.go#L104
    const trimmedWhitespace = lengthBeforeEndTrim - lengthAfterTrim
    if (trimmedWhitespace > 1024) {
        throw Error("too much whitespace at the end of the armor payload")
    }

    if (!armor.startsWith(header)) {
        throw Error(`Armor cannot be decoded if it does not start with a header! i.e. ${header}`)
    }

    if (!armor.endsWith(footer)) {
        throw Error(`Armor cannot be decoded if it does not end with a footer! i.e. ${footer}`)
    }

    const base64Payload = armor.slice(header.length, armor.length - footer.length)
    const lines = base64Payload.split("\n")
    if (lines.some(line => line.length > chunkSize)) {
        throw Error(`Armor to decode cannot have lines longer than ${chunkSize} (configurable) in order to stop padding attacks`)
    }

    if (lines[lines.length - 1].length >= chunkSize) {
        throw Error(`The last line of an armored payload must be less than ${chunkSize} (configurable) to stop padding attacks`)
    }

    return Buffer.from(base64Payload, "base64").toString("binary")
}

export function isProbablyArmored(input: string): boolean {
    return input.startsWith(header)
}
