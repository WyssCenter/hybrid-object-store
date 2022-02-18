// vendor
import React from 'react';
import { Story, Meta } from '@storybook/react/types-6-0';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from '../../../../../AppContext';
// component
import User from '../User';
// data
import context from './UserData';

export default {
  title: 'Layout/Header/Toolbar/User',
  component: User,
  argTypes: {},
} as Meta;


const UserWrapper = () => {
  return (
    <MemoryRouter>
      <AppContext.Provider value={context}>
        <User
          send={(stateValue) => null}
        />
      </AppContext.Provider>
    </MemoryRouter>
  )
}



const Template: Story = (args) => <UserWrapper {...args} />;

export const UserComp = Template.bind({});

UserComp.args = {};


UserComp.parameters = {
  jest: ['User.test.js'],
};
