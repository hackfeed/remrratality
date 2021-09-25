import { BarData, BarOptions } from "@/interfaces/bar";
import { Grid } from "@/interfaces/grid";
import { AnalyticsState } from "@/interfaces/state";
import actions from "@/store/modules/analytics/actions";
import getters from "@/store/modules/analytics/getters";
import mutations from "@/store/modules/analytics/mutations";

const state = {
  periodStart: "2021-01",
  periodEnd: "2021-01",
  files: [],
  file: null,
  grid: {
    title: null,
    cols: null,
    rows: null,
  } as Grid,
  data: {
    labels: [],
    datasets: [
      {
        label: "New",
        data: [],
        backgroundColor: "#027436",
      },
      {
        label: "Old",
        data: [],
        backgroundColor: "#09a776",
      },
      {
        label: "Expansion",
        data: [],
        backgroundColor: "#62da9a",
      },
      {
        label: "Reactivation",
        data: [],
        backgroundColor: "#707fd7",
      },
      {
        label: "Contraction",
        data: [],
        backgroundColor: "#ff8700",
      },
      {
        label: "Churn",
        data: [],
        backgroundColor: "#8f0239",
      },
    ],
  } as BarData,
  dataOptions: {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: true,
        text: "Monthly Reccuring Revenue (Chart)",
        fontSize: 24,
        fontColor: "black",
      },
      tooltips: {
        backgroundColor: "#17BF62",
      },
    },
    scales: {
      x: {
        stacked: true,
        grid: {
          display: false,
        },
      },
      y: {
        stacked: true,
        ticks: {
          beginAtZero: true,
        },
        grid: {
          display: false,
        },
      },
    },
  } as BarOptions,
} as AnalyticsState;

export default {
  namespaced: true,
  state,
  actions,
  getters,
  mutations,
};
