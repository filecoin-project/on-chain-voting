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

import { create } from 'zustand';
interface StoringCidState {
  storingCid: string[];
}

interface StoringHashState {
  storingHash: string[];
}

export const useVoterInfo = create(set => ({
  voterInfo: [],
  setVoterInfo: (newVoterInfo: []) => set({ voterInfo: newVoterInfo }),
}));

export const useCurrentTimezone = create(set => ({
  timezone: '',
  setTimezone: (newTimezone: '') => set({ timezone: newTimezone }),
}));
export const useVotingList = create(set => ({
  votingData: {
    votingList: [],
    totalPage: 10,
  },
  setVotingList: (newData: any) => set({ votingData: newData }),
}));

export const useProposalStatus = create(set => ({
  status: '',
  setStatusList: (status: 0) => set({ status: status }),
}));

export const useStoringCid = create<StoringCidState>((set, get) => ({
  storingCid: localStorage.getItem('storingCid') ? JSON.parse(localStorage.getItem('storingCid')!) : [],

  addStoringCid: (newStoringCid: any[]) => {
    const updatedCid = [...newStoringCid, ...get().storingCid];
    set({ storingCid: updatedCid });
    localStorage.setItem('storingCid', JSON.stringify(updatedCid));
  },

  setStoringCid: (newStoringCid: any[]) => {
    set({ storingCid: newStoringCid });
    localStorage.setItem('storingCid', JSON.stringify(newStoringCid));
  },
}));

export const useStoringHash = create<StoringHashState>((set, get) => ({
  storingHash: localStorage.getItem('storingHash') ? JSON.parse(localStorage.getItem('storingHash')!) : [],

  addStoringHash: (newStoringHash: any[]) => {
    const updatedHash = [...newStoringHash, ...get().storingHash];
    set({ storingHash: updatedHash });
    localStorage.setItem('storingHash', JSON.stringify(updatedHash));
  },

  setStoringHash: (newStoringHash: any[]) => {
    set({ storingHash: newStoringHash });
    localStorage.setItem('storingHash', JSON.stringify(newStoringHash));
  },
}));
export const useSearchValue = create(set => ({
  searchValue: '',
  setSearchValue: (value: '') => set({ searchValue: value }),
}));

export const useFipList = create(set => ({
  data: localStorage.getItem('fipList') ? JSON.parse(localStorage.getItem('fipList')!) : {
    fipList: [],
    totalSize: 0,
    isFipEditorAddress: false,
  },
  setFipList: (value: any[], address: '') => {
    const isHas = value.filter((item: any) => item.editor === address)[0];
    const data = {
      fipList: value,
      totalSize: value.length,
      isFipEditorAddress: isHas ? true : false,
    }
    set({ data });
    localStorage.setItem('fipList', JSON.stringify(data));
  },
}));

export const useTransactionHash = create<any>((set, get) => ({
  transactionHash: localStorage.getItem('transactionHash') ? JSON.parse(localStorage.getItem('transactionHash')!) : {},
  setTransactionHash: (newStoringHash: any) => {
    const updatedHash = { ...get().storingHash, ...newStoringHash };
    set({ transactionHash: updatedHash });
    localStorage.setItem('transactionHash', JSON.stringify(updatedHash));
  },
}));
