<template>
  <div>
    <section class="centered">
      <base-card>
        <h2>Get full control of your app's money flow with mrratality</h2>
        <base-button link to="/analytics">Try now</base-button>
      </base-card>
      <base-card>
        <h2>Quick MRR report in chart style</h2>
        <img src="@/assets/images/mrr.png" alt="Quick MRR Report" />
      </base-card>
      <base-card>
        <h2>... and in table style</h2>
        <img src="@/assets/images/mrrtbl.png" alt="Quick MRR Report" />
      </base-card>
      <base-card>
        <h2>Flexible period adjustment</h2>
        <img src="@/assets/images/periods.png" alt="Flexible period adjustment" />
      </base-card>
      <base-card>
        <h2>Simple data export from CSV</h2>
        <img src="@/assets/images/csv.png" alt="Simple data export from CSV" />
      </base-card>
      <base-card>
        <h2>Upload your data now and get carte blanche of your money flow</h2>
        <base-button link to="/analytics">Upload data</base-button>
      </base-card>
    </section>
  </div>
</template>

<script lang="ts">
import { Vue } from "vue-class-component";

export default class Index extends Vue {}
</script>

<style lang="scss" scoped>
img {
  height: 20rem;
}

.centered {
  text-align: center;
}
</style>
