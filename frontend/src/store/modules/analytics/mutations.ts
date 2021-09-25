import { BarData } from "@/interfaces/bar";
import { File } from "@/interfaces/file";
import { Grid } from "@/interfaces/grid";
import { AnalyticsState } from "@/interfaces/state";

export default {
  setPeriodStart(state: AnalyticsState, data: string): void {
    state.periodStart = data;
  },
  setPeriodEnd(state: AnalyticsState, data: string): void {
    state.periodEnd = data;
  },
  setFile(state: AnalyticsState, data: File): void {
    state.file = data;
  },
  setData(state: AnalyticsState, data: BarData): void {
    state.data = data;
  },
  setGrid(state: AnalyticsState, data: Grid): void {
    state.grid = data;
  },
  setFiles(state: AnalyticsState, data: File[]): void {
    state.files = data;
  },
};
