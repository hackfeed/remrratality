<template>
  <div>
    <base-dialog :show="!!error" title="An error occured" @close="handleError">
      <p>{{ error }}</p>
    </base-dialog>
    <base-dialog :show="isLoading" title="Authenticating..." fixed>
      <base-spinner></base-spinner>
    </base-dialog>
    <base-card>
      <form @submit.prevent="submitForm">
        <div class="form-control">
          <label for="email">E-Mail</label>
          <input type="email" name="email" id="email" v-model.trim="email" />
        </div>
        <div class="form-control">
          <label for="password">Password</label>
          <input type="password" name="password" id="password" v-model.trim="password" />
        </div>
        <p v-if="!formIsValid">
          Please enter a valid email and password (must be at least 5 characters long)
        </p>
        <base-button>{{ submitButtonCaption }}</base-button>
        <base-button type="button" mode="flat" @click="switchAuthMode">{{
          switchModeButtonCaption
        }}</base-button>
      </form>
    </base-card>
  </div>
</template>

<script lang="ts">
import { Vue } from "vue-class-component";

export default class UserAuth extends Vue {
  email = "";
  password = "";
  formIsValid = true;
  mode = "login";
  isLoading: boolean | null = null;
  error: string | null = null;

  get submitButtonCaption(): string {
    if (this.mode === "login") {
      return "Login";
    } else {
      return "Signup";
    }
  }

  get switchModeButtonCaption(): string {
    if (this.mode === "login") {
      return "Signup";
    } else {
      return "Login";
    }
  }

  async submitForm(): Promise<void> {
    this.formIsValid = true;

    if (this.email === "" || !this.email.includes("@") || this.password.length < 5) {
      this.formIsValid = false;
      return;
    }

    this.isLoading = true;
    const actionPayload = { email: this.email, password: this.password };
    try {
      if (this.mode === "login") {
        await this.$store.dispatch("login", actionPayload);
      } else {
        await this.$store.dispatch("signup", actionPayload);
      }
      const redirectUrl = "/" + (this.$route.query.redirect || "analytics");
      this.$router.replace(redirectUrl);
    } catch (error) {
      this.error = error.message || "Failed to authenticate, try later!";
    }

    this.isLoading = false;
  }

  switchAuthMode(): void {
    if (this.mode === "login") {
      this.mode = "signup";
    } else {
      this.mode = "login";
    }
  }

  handleError(): void {
    this.error = null;
  }
}
</script>

<style lang="scss" scoped>
form {
  margin: 1rem;
  padding: 1rem;
}

label {
  font-weight: bold;
  margin-bottom: 0.5rem;
  display: block;
}

input,
textarea {
  display: block;
  width: 100%;
  font: inherit;
  border: 1px solid #ccc;
  padding: 0.15rem;
  :focus {
    border-color: #389948;
    background-color: #faf6ff;
    outline: none;
  }
}

.form-control {
  margin: 0.5rem 0;
}
</style>
