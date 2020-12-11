const presets = ['module:metro-react-native-babel-preset']
const plugins = ['@babel/plugin-transform-flow-strip-types']

plugins.push([
  'module-resolver',
  {
    root: ['./src'],
    extensions: ['.js', '.json'],
    alias: {
      '@': './src',
    },
  },
])

module.exports = {
  presets,
  plugins,
}
