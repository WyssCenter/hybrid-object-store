// vendor
import React from 'react';
// components
import ListPat from '../ListPat';

export default {
  title: 'Pages/tokens/section/list/ListPat',
  component: ListPat,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="ListPat" />
    <ListPat
      title="ListPat Title"
      {...args}
    >
      <div>
        ListPat JSX Context
      </div>
    </ListPat>
</div>);

export const ListPatComp = Template.bind({});

ListPatComp.args = {
  handleClose: (evt) => evt,
  isVisible: true,
  modalType: "namespace",
  postRoute: "namespace",
  updateNamespaceFetchId: (evt) => evt,
};


ListPatComp.parameters = {
  jest: ['ListPat.test.js'],
};
