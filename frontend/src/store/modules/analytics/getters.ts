import { BarData, BarDataOptions } from "@/interfaces/bar";
import { File } from "@/interfaces/file";
import { Grid } from "@/interfaces/grid";
import { AnalyticsState } from "@/interfaces/state";

export default {
  periodStart(state: AnalyticsState): string {
    return state.periodStart;
  },
  periodEnd(state: AnalyticsState): string {
    return state.periodEnd;
  },
  files(state: AnalyticsState): File[] | null {
    return state.files;
  },
  file(state: AnalyticsState): File | null {
    return state.file;
  },
  data(state: AnalyticsState): BarData {
    return state.data;
  },
  grid(state: AnalyticsState): Grid {
    return state.grid;
  },
  dataOptions(state: AnalyticsState): BarDataOptions {
    return state.dataOptions;
  },
};
