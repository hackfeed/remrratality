<template>
  <div>
    <base-dialog :show="!!error" title="An error occured!" @close="handleError">
      <p>{{ error }}</p>
    </base-dialog>
    <base-card class="centered">
      <div v-if="!uploadNew && !isUploaded && !isLoading && analyticsFiles.length > 0">
        <analytics-files
          :files="analyticsFiles"
          @upload-new="setUploadNew"
          @choose-file="setFile"
          @delete-file="deleteFile"
          @is-uploaded="setIsUploaded"
        ></analytics-files>
      </div>
      <div v-else-if="!isUploaded && !isLoading">
        <h2>Please upload CSV file below</h2>
        <analytics-form
          :files-not-empty="analyticsFilesEmpty"
          @upload-data="uploadData"
          @upload-new="setUploadNew"
        ></analytics-form>
      </div>
      <div v-else-if="isLoading">
        <base-spinner></base-spinner>
      </div>
      <div v-else-if="isUploaded">
        <h2>Choose periods for MRR report</h2>
        <analytics-periods-form
          @load-data="loadData"
          @upload-new="setUploadNew"
        ></analytics-periods-form>
      </div>
    </base-card>
    <div v-if="isLoaded" class="centered">
      <base-card>
        <analytics-chart
          :bar-data="analyticsData"
          :bar-options="analyticsOptions"
        ></analytics-chart>
      </base-card>
      <base-card>
        <analytics-grid :grid="analyticsGrid"></analytics-grid>
      </base-card>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

import AnalyticsFiles from "@/components/analytics/AnalyticsFiles.vue";
import AnalyticsForm from "@/components/analytics/AnalyticsForm.vue";
import AnalyticsPeriodsForm from "@/components/analytics/AnalyticsPeriodsForm.vue";
import AnalyticsChart from "@/components/analytics/AnalyticsChart.vue";
import AnalyticsGrid from "@/components/analytics/AnalyticsGrid.vue";

import { BarData, BarOptions } from "@/interfaces/bar";

@Options({
  components: {
    AnalyticsFiles,
    AnalyticsForm,
    AnalyticsPeriodsForm,
    AnalyticsChart,
    AnalyticsGrid,
  },
  watch: {
    file() {
      this.$store.commit("analytics/setPeriodStart", "2021-01");
      this.$store.commit("analytics/setPeriodEnd", "2021-01");
    },
  },
})
export default class Analytics extends Vue {
  uploadNew = false;
  isUploaded = false;
  isLoading = false;
  isLoaded = false;
  file: string | null = null;
  error: string | null = null;

  get analyticsFiles(): string[] {
    return this.$store.getters["analytics/files"];
  }

  get analyticsData(): BarData {
    return this.$store.getters["analytics/data"];
  }

  get analyticsGrid(): Record<string, unknown> {
    return this.$store.getters["analytics/grid"];
  }

  get analyticsOptions(): BarOptions {
    return this.$store.getters["analytics/dataOptions"];
  }

  get analyticsFilesEmpty(): boolean {
    return this.analyticsFiles.length === 0;
  }

  async uploadData(data: Record<string, unknown>): Promise<void> {
    this.isLoading = true;

    try {
      await this.$store.dispatch("analytics/uploadData", data);
      this.isUploaded = true;
      this.file = this.$store.getters["analytics/file"];
      await this.loadFiles();
    } catch (error) {
      this.error = error.message || "Something went wrong!";
    }

    this.isLoading = false;
  }

  async loadData(data: Record<string, unknown>): Promise<void> {
    this.isLoaded = false;
    this.isLoading = true;

    data = {
      filename: this.file,
      period_start: data.periodStart + "-01",
      period_end: data.periodEnd + "-01",
    };
    try {
      await this.$store.dispatch("analytics/loadData", data);
      this.isLoaded = true;
    } catch (error) {
      this.error = error.message || "Something went wrong!";
    }

    this.isLoading = false;
  }

  async loadFiles(): Promise<void> {
    this.isLoading = true;

    try {
      await this.$store.dispatch("analytics/loadFiles");
    } catch (error) {
      this.error = error.message || "Something went wrong!";
    }

    this.isLoading = false;
  }

  async deleteFile(data: Record<string, unknown>): Promise<void> {
    try {
      await this.$store.dispatch("analytics/deleteFile", data);
      await this.loadFiles();
    } catch (error) {
      this.error = error.message || "Something went wrong!";
    }
  }

  setUploadNew(data: boolean): void {
    if (data === false) {
      this.isUploaded = false;
      this.isLoaded = false;
      this.file = null;
    }
    this.uploadNew = data;
  }

  setFile(data: string): void {
    this.file = data;
  }

  setIsUploaded(data: boolean): void {
    this.isUploaded = data;
  }

  handleError(): void {
    this.error = null;
  }

  created(): void {
    this.loadFiles();
  }
}
</script>

<style lang="scss" scoped>
.centered {
  text-align: center;
}
</style>
