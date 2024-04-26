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

export const extractRevertReason = (errorMsg: string) => {
  const revertReasonIndex = errorMsg?.indexOf('revert reason: ');
  if (revertReasonIndex === -1) {
    return null;
  }

  const endOfRevertReasonIndex = errorMsg?.indexOf(', vm error:', revertReasonIndex);
  if (endOfRevertReasonIndex === -1) {
    return null;
  }

  const startIndex = revertReasonIndex + 'revert reason: '.length;
  return errorMsg?.slice(startIndex, endOfRevertReasonIndex)?.trim();
}

export const hasDuplicates = (array: string[]) => {
  return new Set(array).size !== array.length;
}

export const markdownToText = (markdownString: string) => {
  // Remove links ([...](...))
  let noLinks = markdownString.replace(/\[(.+?)\]\(.+?\)/g, '$1');

  // Remove images (![...](...))
  let noImages = noLinks.replace(/!\[(.+?)\]\(.+?\)/g, '$1');

  // Remove inline code blocks（`...`）
  let noInlineCode = noImages.replace(/`([^`]+)`/g, '$1');

  // Remove bold (**...** or __...__)
  let noBold = noInlineCode.replace(/(?:\*{2}|_{2})(.*?)(?:\*{2}|_{2})/g, '$1');

  // Remove italic (*...* or _..._)
  let noItalic = noBold.replace(/(?:\*|_)(.*?)(?:\*|_)/g, '$1');

  // Remove headings (#...，##...，###..., etc.)
  let noHeadings = noItalic.replace(/^#+\s*(.*?)\s*#*$/gm, '$1');

  // Remove unordered lists（-... or *...）
  let noUnorderedList = noHeadings.replace(/^[\s]*[-\*][\s]+(.*?)[\s]*$/gm, '$1');

  // Remove ordered lists（1. ...）
  let noOrderedList = noUnorderedList.replace(/^[\s]*\d+\.[\s]+(.*?)[\s]*$/gm, '$1');

  // Remove strikethrough (~~...~~)
  let noStrikethrough = noOrderedList.replace(/~~(.*?)~~/g, '$1');

  // Remove blockquotes（> ...）
  let noBlockquote = noStrikethrough.replace(/^\s*>\s*(.*?)[\s]*$/gm, '$1');

  // Remove horizontal rules（---）
  let noHorizontalRule = noBlockquote.replace(/^[\s]*[-*_][\s]*[-*_][\s]*[-*_][\s]*$/gm, '');

  // Remove HTML comments rules（<!--...-->）
  let noHTMLComments = noHorizontalRule.replace(/<!--[\s\S]*?-->/g, '');

  // Remove HTML anchor rules（(#...)）
  let noHTMLAnchor = noHTMLComments.replace(/\(#([^)]+)\)/g, '');

  // Remove blank lines
  let finalString = noHTMLAnchor.replace(/^\s*[\r\n]/gm, '');

  return finalString;
}