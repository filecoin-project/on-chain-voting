import create from 'zustand';

export const useVoterInfo = create((set) => ({
  voterInfo: [],
  setVoterInfo: (newVoterInfo: []) => set({ voterInfo: newVoterInfo }),
}));