<script lang="ts">
import { Options, Vue, mixins } from "vue-class-component";
import { BarData, BarDataOptions } from "@/interfaces/bar";
import { Bar } from "vue3-chart-v2";
import { BarOptions } from "chart.js";

class BarTS extends mixins(Bar) {}

@Options({
  props: {
    barData: Object as () => BarData,
    barOptions: Object as () => BarDataOptions,
  },
})
export default class AnalyticsChart extends mixins(Vue, BarTS) {
  readonly barData!: BarData;
  readonly barOptions!: BarOptions;

  mounted(): void {
    this.renderChart(this.barData, this.barOptions);
  }
}
</script>
