import { create } from 'zustand';

interface StoringCidState {
  storingCid: string[];
}

export const useVoterInfo = create(set => ({
  voterInfo: [],
  setVoterInfo: (newVoterInfo: []) => set({ voterInfo: newVoterInfo }),
}));

export const useCurrentTimezone = create(set => ({
  timezone: '',
  setTimezone: (newTimezone: '') => set({ timezone: newTimezone }),
}));

export const useStoringCid = create<StoringCidState>((set, get) => ({
  storingCid: localStorage.getItem('storingCid') ? JSON.parse(localStorage.getItem('storingCid')!) : [],

  addStoringCid: (newStoringCid: string[]) => {
    const updatedCid = [...newStoringCid, ...get().storingCid];
    set({ storingCid: updatedCid });
    localStorage.setItem('storingCid', JSON.stringify(updatedCid));
  },

  setStoringCid: (newStoringCid: string[]) => {
    set({ storingCid: newStoringCid });
    localStorage.setItem('storingCid', JSON.stringify(newStoringCid));
  },
}));