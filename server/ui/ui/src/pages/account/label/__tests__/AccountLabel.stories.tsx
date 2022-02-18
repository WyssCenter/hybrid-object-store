import React from 'react';

import AccountLabel from '../AccountLabel';

export default {
  title: 'Pages/dataset/AccountLabel',
  component: AccountLabel,
  argTypes: {},
};

const Template = (args) => <AccountLabel {...args} />;

export const AccountLabelComp = Template.bind({});

AccountLabelComp.args = {};


AccountLabelComp.parameters = {
  jest: ['AccountLabel.test.js'],
};
