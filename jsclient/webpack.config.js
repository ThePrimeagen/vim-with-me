const path = require('path');

const Dotenv = require('dotenv-webpack')
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
    entry: {
        server_test_libwebsockets: './src/server/libwebsocket.test.ts',
        client: './src/index.ts',
    },
    watch: true,
    plugins: [
        new Dotenv(),
    ],
    module: {
        rules: [{
            test: /.tsx|ts$/,
            exclude: /node_modules/,
            use: {
                loader: 'ts-loader'
            }
        }]
    },

    resolve: {
      extensions: [ '.tsx', '.ts', '.js' ],
    },

    optimization: {
		// We no not want to minimize our code.
		minimize: false
	},

    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'static')
    }
};

