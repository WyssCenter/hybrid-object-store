// vendor
import React from 'react';
// components
import PatItem from '../PatItem';

export default {
  title: 'Pages/tokens/section/list/item/PatItem',
  component: PatItem,
  argTypes: {
    pat: {
      description: String,
      id: Number,
    },
    updateFetchId: Function,
  },

};



const Template = (args) => (<div>
  <div id="PatItem" />
    <PatItem
      title="PatItem Title"
      {...args}
    >
      <div>
        PatItem JSX Context
      </div>
    </PatItem>
</div>);

export const PatItemComp = Template.bind({});

PatItemComp.args = {
  pat: {
    description: 'Gigantum Client',
    id: 2,
  },
  updateFetchId: () => null
};


PatItemComp.parameters = {
  jest: ['PatItem.test.js'],
};
