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
import { filecoin, filecoinCalibration } from 'wagmi/chains';
import {
  oracleCalibrationContractAddress,
  oracleMainNetContractAddress, oraclePowerCalibrationContractAddress,
  powerVotingCalibrationContractAddress,
  powerVotingMainNetContractAddress,
  web3StorageEmail
} from "../common/consts";
import { create } from '@web3-storage/w3up-client';

export const stringToBase64Url = (str: string) => {
  const base64 = btoa(str);

  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

export const bigNumberToFloat = (value: number | string, decimals: number = 18) => {
  return Number(value) ? (Number(value) / (10 ** decimals)).toFixed(2) : '0';
}

/**
 * Convert value to Byte String
 * @param bytes
 */
export const convertBytes = (bytes: number | string) => {
  // Define an array of unit strings representing different byte units
  const units = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  // Initialize the unit index to 0
  let unitIndex = 0;
  // Convert input bytes to a number
  let remainder = Number(bytes);

  // Loop until the remainder is greater than or equal to 1024 and the unit index is within range
  while (remainder >= 1024 && unitIndex < units.length - 1) {
    // Divide the remainder by 1024 to convert to the next unit
    remainder /= 1024;
    unitIndex++;
  }

  // Return a string representing the converted bytes with appropriate unit
  // If remainder is truthy (i.e., not zero), return the formatted string; otherwise, return '0'
  return remainder ? `${remainder.toFixed(2)} ${units[unitIndex]}` : '0';
}

/**
 * Handle MetaMask error message
 * @param errorMsg
 */
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

/**
 * Check duplicate value
 * @param array
 */
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

/**
 * Check whether value isn't empty
 * @param value
 */
export const validateValue = (value: string) => {
  return value?.trim() !== '';
};


const contractAddresses: any = {
  powerVoting: {
    [filecoin.id]: powerVotingMainNetContractAddress,
    [filecoinCalibration.id]: powerVotingCalibrationContractAddress,
  },
  oracle: {
    [filecoin.id]: oracleMainNetContractAddress,
    [filecoinCalibration.id]: oracleCalibrationContractAddress,
  },
  oraclePower: {
    [filecoin.id]: oraclePowerCalibrationContractAddress,
    [filecoinCalibration.id]: oraclePowerCalibrationContractAddress,
  },
};

/**
 * Get smart contract address
 * @param chainId
 * @param type
 */
export const getContractAddress = (chainId: number, type: string) => {
  const chainContracts = contractAddresses[type] || {};
  return chainContracts[chainId];
};

/**
 * Store data into web3.storage
 * @param params for the data
 */
export const getWeb3IpfsId = async (params: object | string) => {
  const client = await create();
  // first time setup!
  if (!Object.keys(client.accounts()).length) {
    // waits for you to click the link in your email to verify your identity
    const account = await client.login(web3StorageEmail);
    // create a space for your uploads
    const space = await client.createSpace('power-voting');
    // save the space to the store, and set as "current"
    await space.save();
    // associate this space with your account
    await account.provision(space.did());
  }

  const json = JSON.stringify(params);
  const data = new Blob([json]);
  const cid = await client.uploadFile(data);
  return cid.toString();
}