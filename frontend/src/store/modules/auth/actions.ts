import { ActionContext } from "vuex";

import { AuthState, RootState } from "@/interfaces/state";

let timer: number;

export default {
  async auth(
    context: ActionContext<AuthState, RootState>,
    payload: Record<string, unknown>
  ): Promise<void> {
    const mode = payload.mode;
    let url = "http://localhost:8081/signup";
    if (mode === "login") {
      url = "http://localhost:8081/login";
    }
    const response = await fetch(url, {
      method: "POST",
      body: JSON.stringify({
        email: payload.email,
        password: payload.password,
      }),
    });
    const responseData = await response.json();

    if (!response.ok) {
      const error = new Error(responseData.message || "Failed to signup. Check your signup data");
      throw error;
    }

    const expirationDate = +responseData.expiresAt;
    const expiresIn = expirationDate * 1000 - new Date().getTime();

    localStorage.setItem("token", responseData.idToken);
    localStorage.setItem("userId", responseData.localId);
    localStorage.setItem("tokenExpiration", expirationDate.toString());

    timer = setTimeout(() => {
      context.dispatch("autoLogout");
    }, expiresIn);

    context.commit("setUser", {
      token: responseData.idToken,
      userId: responseData.localId,
    });
  },
  tryLogin(context: ActionContext<AuthState, RootState>): void {
    const token = localStorage.getItem("token");
    const userId = localStorage.getItem("userId");
    const tokenExpiration = localStorage.getItem("tokenExpiration");

    const expiresIn = +tokenExpiration! * 1000 - new Date().getTime();

    if (expiresIn < 0) {
      return;
    }

    timer = setTimeout(() => {
      context.dispatch("autoLogout");
    }, expiresIn);

    if (token && userId) {
      context.commit("setUser", {
        token: token,
        userId: userId,
      });
    }
  },
  async login(
    context: ActionContext<AuthState, RootState>,
    payload: Record<string, unknown>
  ): Promise<any> {
    return context.dispatch("auth", { ...payload, mode: "login" });
  },
  async signup(
    context: ActionContext<AuthState, RootState>,
    payload: Record<string, unknown>
  ): Promise<any> {
    return context.dispatch("auth", { ...payload, mode: "signup" });
  },
  logout(context: ActionContext<AuthState, RootState>): void {
    localStorage.removeItem("token");
    localStorage.removeItem("userId");
    localStorage.removeItem("tokenExpiration");

    clearTimeout(timer);

    context.commit("setUser", {
      token: null,
      userId: null,
    });
    context.commit("analytics/setFile", null);
  },
  autoLogout(context: ActionContext<AuthState, RootState>): void {
    context.dispatch("logout");
    context.commit("setAutoLogout");
  },
};
