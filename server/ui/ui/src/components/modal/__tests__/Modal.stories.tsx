import React from 'react';

import Modal from '../Modal';

export default {
  title: 'Components/modal/Modal',
  component: Modal,
  argTypes: {},

};

const modalRoot = global.document.createElement('div');
modalRoot.setAttribute('id', 'modal');
const body = global.document.querySelector('body');
body.appendChild(modalRoot);



const Template = (args) => (<div>
  <div id="modal" />
  <Modal
    header="Modal Header"
    subheader="Modal Subheader"
    size="medium"
    handleClose={() => { return }}
    {...args}
  >
    <div>Sample Modal Content</div>
  </Modal>
</div>);

export const ModalComp = Template.bind({
  header: 'Test'
});

ModalComp.args = {
  handleClose: jest.fn(),
  header: 'Add Namespace',
  overflow: true,
  icon: '',
  size: 'medium',
  subheader: 'Create a new namespace here',
};


ModalComp.parameters = {
  jest: ['Modal.test.js'],
};
