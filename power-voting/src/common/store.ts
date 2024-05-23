import { create } from 'zustand';

export const useVoterInfo = create((set) => ({
  voterInfo: [],
  setVoterInfo: (newVoterInfo: []) => set({ voterInfo: newVoterInfo }),
}));

export const useCurrentTimezone = create((set) => ({
  timezone: '',
  setTimezone: (newTimezone: '') => set({ timezone: newTimezone }),
}));