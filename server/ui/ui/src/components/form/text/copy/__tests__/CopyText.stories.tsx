// vendor
import React from 'react';
import { faCopy } from '@fortawesome/free-solid-svg-icons'
// components
import CopyText from '../CopyText';

export default {
  title: 'Components/form/text/CopyText',
  component: CopyText,
  argTypes: {
    icon: []
  },
};



const Template = (args) => (<div>
  <div id="CopyText" />
    <CopyText {...args} />
</div>);

export const CopyTextComp = Template.bind({});

CopyTextComp.args = {
  icon: faCopy,
};


CopyTextComp.parameters = {
  jest: ['CopyText.test.js'],
};
