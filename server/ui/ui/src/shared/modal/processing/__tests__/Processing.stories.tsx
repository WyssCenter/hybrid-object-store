import React from 'react';

import Processing from '../Processing';

export default {
  title: 'Shared/modal/create/Processing',
  component: Processing,
  argTypes: {
    modalType: 'namespace',
    name: 'namespace-12'
  },

};



const Template = (args) => (<div>
  <div id="Processing" />
    <Processing
      title="Processing Title"
      {...args}
    >
      <div>
        Processing JSX Context
      </div>
    </Processing>
</div>);

export const ProcessingComp = Template.bind({});

ProcessingComp.args = {
  modalType: 'namespace',
  name: 'namespace-12'
};


ProcessingComp.parameters = {
  jest: ['Processing.test.js'],
};
