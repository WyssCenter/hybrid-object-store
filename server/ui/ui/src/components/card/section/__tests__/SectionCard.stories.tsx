import React from 'react';

import SectionCard from '../SectionCard';

export default {
  title: 'Components/card/SectionCard',
  component: SectionCard,
  argTypes: {
    verticalHeight: string,
  },
};



const Template = (args) => (<div>
  <div id="SectionCard" />
    <SectionCard
      {...args}
    >
      <div>
        SectionCard JSX Context
      </div>
    </SectionCard>
</div>);

export const SectionCardComp = Template.bind({});

SectionCardComp.args = {
  verticalHeight: 'grid-v-3'
};


SectionCardComp.parameters = {
  jest: ['SectionCard.test.js'],
};
