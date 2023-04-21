import moment from "moment";

export const secondsToYears = (seconds: number) => {
  return Math.round(moment.duration(seconds, "seconds").as("years"));
};

export const yearsToSeconds = (years: number) => {
  return Math.round(moment.duration(years, "years").as("seconds"));
};
