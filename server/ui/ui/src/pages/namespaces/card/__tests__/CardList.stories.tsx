import React from 'react';

import CardList from '../CardList';

export default {
  title: 'Pages/namespaces/card/CardList',
  component: CardList,
  argTypes: {
  },

};


const namespacesList = [{
    bucketName: 'data',
    description: 'this is a namespace',
    name: 'namespace1'
  },
  {
    bucketName: 'data',
    description: 'this is aanother namespace',
    name: 'namespace2'
  }]


const Template = (args) => (<div>
  <div id="CardList" />
    <CardList
      title="CardList Title"
      {...args}
    >
      <div>
        CardList JSX Context
      </div>
    </CardList>
</div>);

export const CardListComp = Template.bind({});

CardListComp.args = {
  namespacesList: namespacesList
};


CardListComp.parameters = {
  jest: ['CardList.test.js'],
};
