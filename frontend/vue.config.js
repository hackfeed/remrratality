module.exports = {
  css: {
    loaderOptions: {
      sass: {
        prependData: `
            @import "@/assets/scss/colors.scss";
            `,
      },
    },
  },
  devServer: {
    proxy: "http://backend_dev:8080",
  },
};
