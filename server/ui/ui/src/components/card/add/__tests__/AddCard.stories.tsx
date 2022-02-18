import React from 'react';

import AddCard from '../AddCard';

export default {
  title: 'Components/card/AddCard',
  component: AddCard,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="AddCard" />
    <AddCard
      title="AddCard Title"
      {...args}
    >
      <div>
        AddCard JSX Context
      </div>
    </AddCard>
</div>);

export const AddCardComp = Template.bind({});

AddCardComp.args = {};


AddCardComp.parameters = {
  jest: ['AddCard.test.js'],
};
