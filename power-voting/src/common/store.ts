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
    searchKey:''
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