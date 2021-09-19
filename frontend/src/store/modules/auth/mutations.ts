import { AuthState } from "@/interfaces/state";

export default {
  setUser(state: AuthState, payload: Record<string, string>): void {
    state.token = payload.token;
    state.userId = payload.userId;
    state.didAutoLogout = false;
  },
  setAutoLogout(state: AuthState): void {
    state.didAutoLogout = true;
  },
};
