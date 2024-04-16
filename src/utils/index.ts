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

export const stringToBase64Url = (str: string) => {
  const base64 = btoa(str);

  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

export const bigNumberToFloat = (value: number | string, decimals: number = 18) => {
  return Number(value) ? (Number(value) / (10 ** decimals)).toFixed(2) : '0';
}

export const convertBytes = (bytes: number | string) => {
  const units = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let unitIndex = 0;
  let remainder = Number(bytes);

  while (remainder >= 1024 && unitIndex < units.length - 1) {
    remainder /= 1024;
    unitIndex++;
  }

  return remainder ? `${remainder.toFixed(2)} ${units[unitIndex]}` : '0';
}