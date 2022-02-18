import React from 'react';

import Tooltip from '../Tooltip';

export default {
  title: 'Components/tooltip/Tooltip',
  component: Tooltip,
  argTypes: {
  },

};




const Template = (args) => (<div>
  <div id="Tooltip" />
    <Tooltip {...args} />
</div>);

export const TooltipComp = Template.bind({});

TooltipComp.args = {};


TooltipComp.parameters = {
  jest: ['Tooltip.test.js'],
};
