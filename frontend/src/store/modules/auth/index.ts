import { AuthState } from "@/interfaces/state";
import actions from "@/store/modules/auth/actions";
import getters from "@/store/modules/auth/getters";
import mutations from "@/store/modules/auth/mutations";

const state = {
  userId: null,
  token: null,
  didAutoLogout: false,
} as AuthState;

export default {
  state,
  actions,
  getters,
  mutations,
};
