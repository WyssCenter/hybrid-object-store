import React from 'react';

import DivisionText from '../DivisionText';

export default {
  title: 'Components/division/DivisionText',
  component: DivisionText,
  argTypes: {
    text: String
  },

};



const Template = (args) => (<div>
  <div id="DivisionText" />
    <DivisionText {...args} />
</div>);

export const DivisionTextComp = Template.bind({});

DivisionTextComp.args = {
  text: 'Header Section'
};


DivisionTextComp.parameters = {
  jest: ['DivisionText.test.js'],
};
