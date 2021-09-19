<template>
  <section>
    <div>
      <label for="file" class="flat">Browse</label>
      <input type="file" id="file" ref="file" @change="handleFileUpload" />
    </div>
    <div>
      <base-button v-if="!filesNotEmpty" mode="flat" @click="uploadNew(false)"
        >Back to files</base-button
      >
      <base-button @click="submitFile">Upload</base-button>
    </div>
  </section>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({
  props: {
    filesNotEmpty: String,
  },
  emits: ["upload-data", "upload-new"],
})
export default class AnalyticsForm extends Vue {
  file: File | string = "";
  readonly filesNotEmpty!: string;

  handleFileUpload(): void {
    const files = this.$refs.file as HTMLInputElement;
    if (files && files.files) {
      this.file = files.files[0];
    }
  }

  submitFile(): void {
    const formData = new FormData();
    formData.append("file", this.file);
    this.$emit("upload-data", formData);
  }

  uploadNew(toUpload: boolean): void {
    this.$emit("upload-new", toUpload);
  }
}
</script>

<style lang="scss" scoped>
label {
  margin-bottom: 1rem;
}

#file {
  opacity: 0;
  position: absolute;
  z-index: -1;
}

.flat {
  text-decoration: none;
  padding: 0.75rem 1.5rem;
  font: inherit;
  background-color: transparent;
  color: $color-primary;
  border: none;
  cursor: pointer;
  border-radius: 30px;
  margin-right: 0.5rem;
  display: inline-block;
  &:hover,
  &:active {
    background-color: $color-primary-light;
    color: $color-neutral-primary;
  }
}
</style>
