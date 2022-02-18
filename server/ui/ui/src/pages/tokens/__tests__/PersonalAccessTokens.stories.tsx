// vendor
import React from 'react';
// components
import PersonalAccessTokens from '../PersonalAccessTokens';

export default {
  title: 'Pages/tokens/PersonalAccessTokens',
  component: PersonalAccessTokens,
  argTypes: {
  },

};



const Template = (args) => (<div>
  <div id="PersonalAccessTokens" />
    <PersonalAccessTokens
      title="PersonalAccessTokens Title"
      {...args}
    >
      <div>
        PersonalAccessTokens JSX Context
      </div>
    </PersonalAccessTokens>
</div>);

export const PersonalAccessTokensComp = Template.bind({});

PersonalAccessTokensComp.args = {
};


PersonalAccessTokensComp.parameters = {
  jest: ['PersonalAccessTokens.test.js'],
};
