<template>
  <section>
    <div class="periods-selection">
      <div class="periods-selection__period">
        <label for="period-start">Period Start</label>
        <input type="month" name="period-start" id="period-start" v-model="periodStart" />
      </div>
      <div class="periods-selection__period">
        <label for="period-end">Period End</label>
        <input type="month" name="period-end" id="period-end" v-model="periodEnd" />
      </div>
    </div>
    <div>
      <base-button mode="flat" @click="uploadNew(false)">Back to files</base-button>
      <base-button @click="submitPeriods">Load report</base-button>
    </div>
  </section>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  emits: ["upload-data", "upload-new"],
})
export default class AnalyticsPeriodsForm extends Vue {
  periodStart: string = this.$store.getters["analytics/periodStart"];
  periodEnd: string = this.$store.getters["analytics/periodEnd"];

  submitPeriods(): void {
    this.$store.commit("analytics/setPeriodStart", this.periodStart);
    this.$store.commit("analytics/setPeriodEnd", this.periodEnd);
    this.$emit("load-data", {
      periodStart: this.periodStart,
      periodEnd: this.periodEnd,
    });
  }

  uploadNew(toUpload: boolean): void {
    this.$emit("upload-new", toUpload);
  }
}
</script>

<style lang="scss" scoped>
label {
  font-weight: bold;
  display: block;
  margin-bottom: 0.5rem;
}

input[type="month"] {
  border: 1px solid #389948;
  border-radius: 5px;
  color: black;
  text-align: center;
  font-family: "Roboto", sans-serif;
  cursor: pointer;
  &:focus {
    outline: none;
  }
}

.periods-selection {
  display: flex;
}

.periods-selection__period {
  margin: 1rem auto;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
</style>
