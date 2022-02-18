// vendor
import React from 'react';
// components
import NamespaceFields from '../NamespaceFields';

export default {
  title: 'Shared/modal/create/NamespaceFields',
  component: NamespaceFields,
  argTypes: {
    handleBucketNameEvent: Function,
    handleObjectChangeEvent: Function,
    objectStoreList: Array,
  },

};



const Template = (args) => (<div>
  <div id="NamespaceFields" />
    <NamespaceFields
      title="NamespaceFields Title"
      {...args}
    />
</div>);

export const NamespaceFieldsComp = Template.bind({});

NamespaceFieldsComp.args = {
  handleBucketNameEvent: () => null,
  handleObjectChangeEvent: () => null,
  objectStoreList: ['default', 'refault'],
};


NamespaceFieldsComp.parameters = {
  jest: ['CreateModal.test.js'],
};
