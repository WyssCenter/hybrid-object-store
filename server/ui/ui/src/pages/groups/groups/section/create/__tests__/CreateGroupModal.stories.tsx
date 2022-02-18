// vendor
import React from 'react';
// component
import CreateGroupModal from '../CreateGroupModal';

export default {
  title: 'Pages/groups/groups/section/create/CreateGroup',
  component: CreateGroupModal,
  argTypes: {},
};

const Template = (args) => <CreateGroupModal {...args} />;

export const CreateGroupModalComp = Template.bind({});

CreateGroupModalComp.args = {
  sectionType: "user"
};


CreateGroupModalComp.parameters = {
  jest: ['CreateGroupModal.test.js'],
};
