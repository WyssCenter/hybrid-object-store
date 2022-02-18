 const path = require('path');
 const requireContext = require('require-context');
 const react = require('@storybook/react');

 const req = requireContext('../../src/', true, /\.stories\.js$/); // <- import all the stories at once

 react.configure(req, module);

module.exports = {
  webpackFinal: async (config, { configType }) => {
    // `configType` has a value of 'DEVELOPMENT' or 'PRODUCTION'
    // You can change the configuration based on that.
    // 'PRODUCTION' is used when building the static version of storybook.

    // Make whatever fine-grained changes you need
    config.module.rules.push({
      test: /\.scss$/,
      use: ['style-loader', 'css-loader', 'sass-loader'],
      include: path.resolve(__dirname, '../'),
    });


    config.resolve.alias = {
      // with-react-native-for-web/
      'react-native': 'react-native-web',
      'Styles': path.resolve(__dirname, '../src/styles'),
      'Hooks': path.resolve(__dirname, '../src/hooks'),
      'Shared': path.resolve(__dirname, '../src/shared'),
      'Images': path.resolve(__dirname, '../src/images'),
      'Fonts': path.resolve(__dirname, '../src/fonts'),
      'Components': path.resolve(__dirname, '../src/components'),
      'Layout': path.resolve(__dirname, '../src/layout'),
      'Environment': path.resolve(__dirname, '../src/environment'),
      'Pages': path.resolve(__dirname, '../src/pages'),
      'Src': path.resolve(__dirname, '../src'),
    };

    // Return the altered config
    return config;
  },
  "stories": [
    "../src/**/*.stories.mdx",
    "../src/**/*.stories.@(js|jsx|ts|tsx)"
  ],
  "addons": [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    "@storybook/addon-jest",
    "@storybook/addon-storyshots",
    "@storybook/addon-storysource",
    "@storybook/addon-a11y",
  ]
}
