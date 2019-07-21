const webpack = require("webpack")
const path = require("path")
const MiniCssExtractPlugin = require("mini-css-extract-plugin")
const { CleanWebpackPlugin } = require("clean-webpack-plugin")
const HtmlWebpackPlugin = require("html-webpack-plugin")
const TerserJSPlugin = require("terser-webpack-plugin")
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin")
const PACKAGE = require("./package.json")

const PATHS = {
  SRC: path.resolve(__dirname, "src"),
  DIST: path.resolve(__dirname, "dist"),
  NODE_MODULES: path.resolve(__dirname, "node_modules"),
  PUBLIC: path.resolve(__dirname, "/"),
}

module.exports = (env, argv) => {
  const mode = argv.mode
  const prod = mode === "production"

  return {
    entry: `${PATHS.SRC}/index.js`,
    resolve: {
      modules: [PATHS.SRC, PATHS.NODE_MODULES],
      alias: {
        svelte: path.resolve("node_modules", "svelte")
      },
      extensions: [".mjs", ".js", ".svelte"],
      mainFields: ["svelte", "browser", "module", "main"]
    },
    output: {
      path: __dirname + "/dist",
      filename: "[name].js",
      chunkFilename: "[name].js"
    },
    module: {
      rules: [
        {
          test: /\.svelte$/,
          use: {
            loader: "svelte-loader",
            options: {
              emitCss: true,
              hotReload: true
            }
          }
        },
        {
          test: /\.(sa|sc|c)ss$/,
          use: [
            {
              loader: MiniCssExtractPlugin.loader,
              options: {
                hmr: !prod,
              },
            },
            "css-loader",
            "sass-loader",
          ],
        }
      ]
    },
    plugins: [
      new CleanWebpackPlugin(),
      new HtmlWebpackPlugin({
        title: "kube-vault",
        template: `${PATHS.SRC}/index.html`,
        filename: "index.html",
        inject: true,
        minify: {
          collapseWhitespace: false,
          collapseInlineTagWhitespace: false,
          removeComments: true,
          removeRedundantAttributes: true,
        },
        inlineSource: /^.*$/,
      }),
      new MiniCssExtractPlugin({
        filename: "[name].css"
      }),
      new webpack.DefinePlugin({
        "process.env": {
          VERSION: JSON.stringify(PACKAGE.version),
          API_HOST: JSON.stringify(!prod ? "http://localhost:8080" : ""),
          NAMESPACE: JSON.stringify(process.env.NAMESPACE ? process.env.NAMESPACE : "default"),
        },
      }),
    ],
    optimization: {
      minimizer: [new TerserJSPlugin({}), new OptimizeCSSAssetsPlugin({})],
      splitChunks: {
        cacheGroups: {
          vendor: {
            test: /node_modules/,
            chunks: "initial",
            name: "vendor",
            enforce: true,
          },
          main: {
            test: /src/,
            chunks: "initial",
            name: "main",
            enforce: true,
          },
        },
      },
    },
    mode,
    devtool: prod ? false : "source-map",
    devServer: {
      port: 8081,
      host: "localhost",
    },
  }
}