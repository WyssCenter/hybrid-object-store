import React from 'react';

import Card from '../Card';

export default {
  title: 'Components/card/Card',
  component: Card,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="Card" />
    <Card
      title="Card Title"
      {...args}
    >
      <div>
        Card JSX Context
      </div>
    </Card>
</div>);

export const CardComp = Template.bind({});

CardComp.args = {};


CardComp.parameters = {
  jest: ['Card.test.js'],
};
