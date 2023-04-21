export const sortAddress = (addr?: string, len: number = 4): string => {
  if (!addr) return "0x";
  if (addr.length <= len * 2 + 2) {
    return addr;
  }

  return `${addr.slice(0, len + 2)}...${addr.slice(-len)}`;
};

export const isEmptyAddress = (addr: string) => {
  return addr == "0x0000000000000000000000000000000000000000" || addr == "0x" || addr == "0" || addr == "";
};
