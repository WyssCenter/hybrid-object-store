// vendor
import React from 'react';
// component
import Loader from '../Loader';

export default {
  title: 'Components/loader/Loader',
  component: Loader,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="Loader" role="img"/>
    <Loader {...args} />
</div>);

export const LoaderComp = Template.bind({});

LoaderComp.args = {
  nested: false,
};


LoaderComp.parameters = {
  jest: ['Loader.test.js'],
};
