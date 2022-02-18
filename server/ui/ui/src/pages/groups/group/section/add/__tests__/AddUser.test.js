// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
// context
import AppContext from 'Src/AppContext'
import GroupContext from '../../../GroupContext'
// components
import AddUser from '../AddUser';



const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin',
    groups: 'blah,fah,group-1'
  }
}

const privilegedRoleWithoutAccess = {
  profile: {
    name: 'privileged',
    role: 'privileged',
    groups: 'blah,fah'
  }
};



const privilegedRoleWithAccess = {
  profile: {
    name: 'privileged',
    role: 'privileged',
    groups: 'blah,fah,group-1'
  }
}


const groupname = 'group-1';

jest.mock('Environment/createEnvironment', () => {
  return {
    put: (route) => new Promise((resolve) => {
      if (route.indexOf('error') === -1) {
        resolve(
          {
            json: () => new Promise((resolve) => { resolve({success: true})})
          }
        )
      } else {
        resolve(
          {
            json: () => new Promise((resolve) => { resolve({error: 'Name value does not match spec'})})
          }
        )
      }
    })
  }
});


describe('AddUser', () => {
  it('Renders component', () => {
    const send = jest.fn();
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <GroupContext.Provider value={{ send, groupname }}>
          <AddUser />
        </GroupContext.Provider>
      </AppContext.Provider>
    );

    expect(screen.getByText(/Who do you want/i)).toBeInTheDocument();
  })

  it('Values update in textboxes', () => {
    const send = jest.fn();
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <GroupContext.Provider value={{ send, groupname }}>
          <AddUser />
        </GroupContext.Provider>
      </AppContext.Provider>
    );
    const inputName = screen.queryAllByRole('textbox');
    fireEvent.keyUp(inputName[0], { target: { value: 'name' }});

    expect(inputName[0].value).toBe('name');
  });


  it('Submit Button Fires post and clears inputs', async () => {
    const send = jest.fn();
    act(() => {
      render(
        <AppContext.Provider value={{ user: adminRole }}>
          <GroupContext.Provider value={{ send, groupname }}>
            <AddUser />
          </GroupContext.Provider>
        </AppContext.Provider>
      );
    });
    const inputName = screen.queryAllByRole('textbox');
    fireEvent.keyUp(inputName[0], { target: { value: 'name' }});
    const button = screen.getByRole('button');
    fireEvent.click(button);
    await waitFor(() => jest.mock)

    expect(inputName[0].value).toBe('');
  });


  it('Submit Button Fires post with error and clears inputs', async () => {
    const send = jest.fn();
    act(() => {
      render(
        <AppContext.Provider value={{ user: adminRole }}>
          <GroupContext.Provider value={{ send, groupname }}>
            <AddUser />
          </GroupContext.Provider>
        </AppContext.Provider>
      );
    });

    const inputName = screen.queryAllByRole('textbox');
    fireEvent.keyUp(inputName[0], { target: { value: 'error -sdsd' }});
    const button = screen.getByRole('button');
    fireEvent.click(button);
    await waitFor(() => jest.mock)


    const errorElement = screen.getByText(/Name value does not match spec/i);
    expect(errorElement).toBeInTheDocument();
  });



  it('Priveleged user that is a member can edit', async () => {
    const send = jest.fn();
    act(() => {
      render(
        <AppContext.Provider value={{ user: privilegedRoleWithAccess }}>
          <GroupContext.Provider value={{ send, groupname }}>
            <AddUser />
          </GroupContext.Provider>
        </AppContext.Provider>
      );
    });

    const inputName = screen.getByRole('textbox');
    const buttons = screen.queryAllByRole('button');
    const warningElement = screen.queryAllByText(/privileged users must be a member of a group to have write access/i);


    expect(inputName).not.toBeDisabled();
    expect(buttons[0]).not.toBeDisabled();
    expect(warningElement.length).toBe(0);

  });


  it('Priveleged user that is not a member cannot edit', async () => {
    const send = jest.fn();
    act(() => {
      render(
        <AppContext.Provider value={{ user: privilegedRoleWithoutAccess }}>
          <GroupContext.Provider value={{ send, groupname }}>
            <AddUser />
          </GroupContext.Provider>
        </AppContext.Provider>
      );
    });

    const inputName = screen.getByRole('textbox');
    const buttons = screen.queryAllByRole('button');
    const warningElement = screen.getByText(/privileged users must be a member of a group to have write access/i);

    expect(inputName).toBeDisabled();
    expect(buttons[0]).toBeDisabled();
    expect(warningElement).toBeInTheDocument();

  });

})
