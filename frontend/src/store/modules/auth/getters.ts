import { AuthState } from "@/interfaces/state";

export default {
  userId(state: AuthState): string | null {
    return state.userId;
  },
  token(state: AuthState): string | null {
    return state.token;
  },
  isAuthenticated(state: AuthState): boolean {
    return !!state.token;
  },
  didAutoLogout(state: AuthState): boolean {
    return state.didAutoLogout;
  },
};
