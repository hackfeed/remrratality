import { BarData, BarOptions } from "./bar";
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
  dataOptions: BarOptions;
}

export interface AuthState {
  userId: string | null;
  token: string | null;
  didAutoLogout: boolean;
}
