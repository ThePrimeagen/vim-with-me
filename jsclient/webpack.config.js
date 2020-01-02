const path = require('path');

const Dotenv = require('dotenv-webpack')
const CopyPlugin = require('copy-webpack-plugin');

const htmlCopy = [
    'player',
    'commander',
    'index',
    'canvasFun',
    'editor',
];
const folderCopy = [
    './resources'
];

const staticPath = path.join(__dirname, 'static');
const htmlPath = path.join(__dirname, 'html');

function getCopyPaths() {
    return htmlCopy.
        map(url => ({
            from: path.join(htmlPath, `${url}.html`),
            to: staticPath
        })).
        concat(folderCopy.map(p => ({
            from: path.join(__dirname, p),
            to: path.join(staticPath, 'resources')
        })));
}

module.exports = {
    entry: {
        player: './lib/player/index.ts',
        commander: './lib/commander/index.tsx',
        home: './lib/home/index.tsx',
        editor: './lib/editor/index.tsx',
        canvasFun: './lib/canvasFun/index.ts',
    },
    watch: true,
    plugins: [
        new CopyPlugin(getCopyPaths()),
        new Dotenv(),
    ],
    module: {
        rules: [{
            test: /.tsx|ts$/,
            exclude: /node_modules/,
            use: {
                loader: 'ts-loader'
            }
        }, {
            test: /.wasm$/,
            exclude: /node_modules/,
            use: {
                loader: 'wasm-loader'
            }
        }]
    },

    resolve: {
      extensions: [ '.tsx', '.ts', '.js', '.wasm' ],
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

