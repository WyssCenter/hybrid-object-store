import React from 'react';

import Progress from '../Progress';

export default {
  title: 'Components/loader/Progress',
  component: Progress,
  argTypes: {
    isCanceling: boolean,
    isComplete: boolean,
    error: string,
    percentageComplete: number,
    text: string,
  },

};



const Template = (args) => (<div>
  <div id="Progress" />
    <Progress {...args} />
</div>);

export const ProgressComp = Template.bind({});

ProgressComp.args = {
  isCanceling: false,
  isComplete: false,
  error: 'Error',
  percentageComplete: 33,
  text: 'Updating data',
};


ProgressComp.parameters = {
  jest: ['Progress.test.js'],
};
