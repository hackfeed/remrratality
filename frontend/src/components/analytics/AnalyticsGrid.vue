<template>
  <h2>{{ grid.title }}</h2>
  <div id="wrapper"></div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { Grid } from "gridjs";
import { Grid as GridType } from "@/interfaces/grid";
import "gridjs/dist/theme/mermaid.css";

@Options({
  props: {
    grid: Object as () => GridType,
  },
  emits: ["upload-data", "upload-new"],
})
export default class AnalyticsGrid extends Vue {
  readonly grid!: GridType;

  mounted(): void {
    const wrapper = document.getElementById("wrapper");
    if (wrapper) {
      new Grid({
        columns: this.grid.cols,
        data: this.grid.rows,
        search: true,
      }).render(wrapper);
    }
  }
}
</script>
