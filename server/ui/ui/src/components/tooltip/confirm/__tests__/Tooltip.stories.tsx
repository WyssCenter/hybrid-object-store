import React from 'react';

import TooltipConfirm from '../TooltipConfirm';

export default {
  title: 'Components/tooltip/TooltipConfirm',
  component: TooltipConfirm,
  argTypes: {
  },

};




const Template = (args) => (<div>
  <div id="TooltipConfirm" />
    <TooltipConfirm {...args} />
</div>);

export const TooltipConfirmComp = Template.bind({});

TooltipConfirmComp.args = {};


TooltipConfirmComp.parameters = {
  jest: ['TooltipConfirm.test.js'],
};
