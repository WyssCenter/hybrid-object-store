// vendor
import React from 'react';
// components
import PatSection from '../PatSection';

export default {
  title: 'Pages/tokens/section/PatSection',
  component: PatSection,
  argTypes: {
  },

};



const Template = (args) => (
  <div>
    <div id="PatSection" />
      <PatSection
        title="PatSection Title"
        {...args}
      >
        <div>
          PatSection JSX Context
        </div>
      </PatSection>
  </div>
);

export const PatSectionComp = Template.bind({});

PatSectionComp.args = {
};


PatSectionComp.parameters = {
  jest: ['PatSection.test.js'],
};
