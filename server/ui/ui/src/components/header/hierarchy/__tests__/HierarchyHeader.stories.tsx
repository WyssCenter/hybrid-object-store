// vendor
import React from 'react';
// components
import HierarchyHeader from '../HierarchyHeader';

export default {
  title: 'Components/header/HierarchyHeader',
  component: HierarchyHeader,
  argTypes: {
    header: String,
    subheader: String
  },
};



const Template = (args) => (<div>
  <div id="HierarchyHeader" />
    <HierarchyHeader {...args} />
</div>);

export const HierarchyHeaderComp = Template.bind({});

HierarchyHeaderComp.args = {
  header: 'Header',
  subheaeder: 'Subheader'
};


HierarchyHeaderComp.parameters = {
  jest: ['HierarchyHeader.test.js'],
};
