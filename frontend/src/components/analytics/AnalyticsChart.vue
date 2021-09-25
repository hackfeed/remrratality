<template>
  <BarChart :chartData="barData" :options="barOptions" />
</template>

<script lang="ts">
import { BarData, BarOptions } from "@/interfaces/bar";
import { Vue, Options } from "vue-class-component";
import { BarChart } from "vue-chart-3";
import { Chart, registerables } from "chart.js";

Chart.register(...registerables);

@Options({
  props: {
    barData: Object as () => BarData,
    barOptions: Object as () => BarOptions,
  },
  components: { BarChart },
})
export default class AnalyticsChart extends Vue {
  readonly barData!: BarData;
  readonly barOptions!: BarOptions;
}
</script>
