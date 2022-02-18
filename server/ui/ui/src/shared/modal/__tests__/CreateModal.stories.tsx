import React from 'react';

import CreateModal from '../CreateModal';

export default {
  title: 'Shared/modal/create/CreateModal',
  component: CreateModal,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="CreateModal" />
    <CreateModal
      title="CreateModal Title"
      {...args}
    >
      <div>
        CreateModal JSX Context
      </div>
    </CreateModal>
</div>);

export const CreateModalComp = Template.bind({});

CreateModalComp.args = {
  handleClose: (evt) => evt,
  isVisible: true,
  modalType: "namespace",
  postRoute: "namespace",
  updateNamespaceFetchId: (evt) => evt,
};


CreateModalComp.parameters = {
  jest: ['CreateModal.test.js'],
};
