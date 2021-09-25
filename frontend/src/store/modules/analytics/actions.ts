import { ActionContext } from "vuex";

import { File } from "@/interfaces/file";
import { AnalyticsState, RootState } from "@/interfaces/state";

export default {
  async uploadData(context: ActionContext<AnalyticsState, RootState>, data: string): Promise<void> {
    const response = await fetch("http://localhost:8081/files/upload", {
      method: "POST",
      body: data,
      headers: {
        token: localStorage.getItem("token")!,
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to fill data");
      throw error;
    }

    context.commit("setFile", responseData.filename);
  },
  async loadData(
    context: ActionContext<AnalyticsState, RootState>,
    data: Record<string, unknown>
  ): Promise<void> {
    const response = await fetch("http://localhost:8081/analytics", {
      method: "POST",
      body: JSON.stringify(data),
      headers: {
        token: localStorage.getItem("token")!,
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to load analytics.");
      throw error;
    }

    const curData = { ...context.rootGetters["analytics/data"] };
    const mrr = responseData.mrr;
    const months = responseData.months;

    curData.labels = months;

    for (const el of curData.datasets) {
      if (el.label === "New") {
        el.data = mrr.New;
      }
      if (el.label === "Old") {
        el.data = mrr.Old;
      }
      if (el.label === "Expansion") {
        el.data = mrr.Expansion;
      }
      if (el.label === "Reactivation") {
        el.data = mrr.Reactivation;
      }
      if (el.label === "Contraction") {
        el.data = mrr.Contraction;
      }
      if (el.label === "Churn") {
        el.data = mrr.Churn;
      }
    }

    const grid = {
      title: "Monthly Reccuring Revenue (Table)",
      cols: [""].concat(months),
      rows: [
        ["New"].concat(mrr.New),
        ["Old"].concat(mrr.Old),
        ["Expansion"].concat(mrr.Expansion),
        ["Reactivation"].concat(mrr.Reactivation),
        ["Contraction"].concat(mrr.Contraction),
        ["Churn"].concat(mrr.Churn),
        ["MRR"].concat(mrr.Total),
      ],
    };

    context.commit("setGrid", grid);
    context.commit("setData", curData);
  },
  async loadFiles(context: ActionContext<AnalyticsState, RootState>): Promise<void> {
    const response = await fetch("http://localhost:8081/files/load", {
      headers: {
        token: localStorage.getItem("token")!,
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to load files");
      throw error;
    }

    context.commit("setFiles", responseData.files);
  },
  async deleteFile(
    context: ActionContext<AnalyticsState, RootState>,
    data: Record<string, unknown>
  ): Promise<void> {
    const reqData = {
      name: data,
    };

    const response = await fetch("http://localhost:8081/files/delete", {
      method: "POST",
      body: JSON.stringify(reqData),
      headers: {
        token: localStorage.getItem("token")!,
      },
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to delete file");
      throw error;
    }

    const files = context.rootGetters["analytics/files"].filter(
      (file: File) => file.name != (reqData as Record<string, unknown>).name
    );

    context.commit("setFiles", files);
  },
};
