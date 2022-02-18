import React from 'react';

import Dropdown from '../Dropdown';

export default {
  title: 'Components/dropdown/Dropdown',
  component: Dropdown,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="Dropdown" />
    <Dropdown {...args} />
</div>);

export const DropdownComp = Template.bind({});

DropdownComp.args = {};


DropdownComp.parameters = {
  jest: ['Dropdown.test.js'],
};
