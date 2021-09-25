import { BarData, BarDataOptions } from "./bar";
import { File } from "./file";
import { Grid } from "./grid";

export class RootState {}

export interface AnalyticsState {
  periodStart: string;
  periodEnd: string;
  files: File[] | null;
  file: File | null;
  grid: Grid;
  data: BarData;
  dataOptions: BarDataOptions;
}

export interface AuthState {
  userId: string | null;
  token: string | null;
  didAutoLogout: boolean;
}
