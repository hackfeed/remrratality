<template>
  <section>
    <h2>Uploaded files</h2>
    <ul>
      <div @click="chooseFile(file)" class="files-selection" v-for="file in files" :key="file">
        <p class="files-selection__filename">{{ file.name }}</p>
        <p class="files-selection__uploadtime">Uploaded at {{ parseDate(file.uploaded_at) }}</p>
        <img @click.stop="deleteFile(file)" src="@/assets/images/remove.png" alt="Remove file" />
      </div>
    </ul>
    <base-button @click="uploadNew(true)">Upload new</base-button>
  </section>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { File } from "@/interfaces/file";

@Options({
  props: {
    files: Object as () => File[],
  },
  emits: ["upload-new", "choose-file", "delete-file", "is-uploaded"],
})
export default class AnalyticsFiles extends Vue {
  readonly files!: File[];

  parseDate(date: string): string {
    const unixTime = Date.parse(date);
    const hsm = new Date(unixTime).toLocaleTimeString("ru-RU");
    const dt = new Date(unixTime).toLocaleDateString("ru-RU");
    return `${hsm} ${dt}`;
  }

  uploadNew(toUpload: boolean): void {
    this.$emit("upload-new", toUpload);
  }

  chooseFile(file: File): void {
    this.$emit("choose-file", file.name);
    this.$emit("is-uploaded", true);
  }

  deleteFile(file: File): void {
    this.$emit("delete-file", file.name);
  }
}
</script>

<style lang="scss" scoped>
ul {
  padding: 0;
}

img {
  width: 1rem;
}

.files-selection {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 0 0 0.5rem;
  padding: 0 0.2rem;
  border: 1px $color-border-primary solid;
  border-radius: 0.5rem;
  &:hover {
    cursor: pointer;
    border-color: none;
    background: $color-neutral-hover;
  }
}

.files-selection__filename {
  color: $color-primary;
  font-weight: bold;
}

.files-selection__uploadtime {
  color: $color-neutral-secondary;
}
</style>
