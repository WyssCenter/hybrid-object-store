import React from 'react';

import Search from '../Search';

export default {
  title: 'Components/form/text/Search',
  component: Search,
  argTypes: {},
};

const Template = (args) => <Search {...args} />;

export const SearchComp = Template.bind({});

SearchComp.args = {
  placeholder: 'Search Datasets',
  list: [],
  updateList: () => null,
};

SearchComp.parameters = {
  jest: ['Search.test.js'],
};
