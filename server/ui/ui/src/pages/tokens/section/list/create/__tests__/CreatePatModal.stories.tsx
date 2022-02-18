// vendor
import React from 'react';
// components
import CreatePatModal from '../CreatePatModal';

export default {
  title: 'Pages/tokens/section/create/CreatePatModal',
  component: CreatePatModal,
  argTypes: {
  },

};


const Template = (args) => (<div>
  <div id="CreatePatModal" />
    <CreatePatModal
      title="CreatePat Title"
      {...args}
    >
      <div>
        CreatePat JSX Context
      </div>
    </CreatePatModal>
</div>);

export const CreatePatModalComp = Template.bind({});

CreatePatModalComp.args = {
  handleClose: (evt) => evt,
  isVisible: true,
  modalType: "namespace",
  postRoute: "namespace",
  updateNamespaceFetchId: (evt) => evt,
};


CreatePatModalComp.parameters = {
  jest: ['CreatePatModal.test.js'],
};
