import React from 'react';

import Success from '../Success';

export default {
  title: 'Shared/modal/create/Success',
  component: Success,
  argTypes: {
    modalType: string,
    namespace: string,
  },

};



const Template = (args) => (<div>
  <div id="Success" />
    <Success
      title="Success Title"
      {...args}
    >
      <div>
        Success JSX Context
      </div>
    </Success>
</div>);

export const SuccessComp = Template.bind({});

SuccessComp.args = {
  modalType: 'namespace',
  namespace: 'namespace-12'
};


SuccessComp.parameters = {
  jest: ['Success.test.js'],
};
