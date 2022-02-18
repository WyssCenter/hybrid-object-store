// vendor
import React from 'react';
// components
import NewPat from '../NewPat';

const pat = {
  token: "r4nD0mStr1ngF0rTh3T0k3n",
  id: 2,
}

export default {
  title: 'Pages/tokens/section/create/new/NewPat',
  component: NewPat,
  argTypes: {
    dissmissPat: Function,
    pat: Object,
  },

};


const Template = (args) => (<div>
  <div id="NewPat" />
    <NewPat
      title="NewPat Title"
      {...args}
    >
      <div>
        NewPat JSX Context
      </div>
    </NewPat>
</div>);

export const NewPatComp = Template.bind({});

NewPatComp.args = {
  dismissPat: () => null,
  pat,

};


NewPatComp.parameters = {
  jest: ['NewPat.test.js'],
};
