import { createStore } from "vuex";

import { RootState } from "@/interfaces/state";

import analytics from "./modules/analytics";
import auth from "./modules/auth";

export default createStore<RootState>({
  state: new RootState(),
  modules: {
    auth,
    analytics,
  },
});
