// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
import GroupContext from '../../../../GroupContext';
// components
import UserItem from '../UserItem';

jest.mock('Environment/createEnvironment', () => {
  return {
    del: () => new Promise((resolve) => {
      resolve(
        {
          ok: true,
          json: () => new Promise((resolve) => { resolve({ success: true})})
        }
      )
    })
  }
});

const adminRole = {
  profile: {
    name: 'admin',
    role: 'admin'
  }
}

const privilegedRole = {
  profile: {
    name: 'privileged',
    role: 'privileged'
  }
}

const userRole = {
  profile: {
    name: 'user',
    role: 'user'
  }
}

const user = {"username":"admin","full_name":"admin","given_name":"*","family_name":"*","email":"admin@example.com","email_verified":true,"role":"admin"}

const groupname = 'group-1';


describe('UserItem', () => {
  it('Renders group data', () => {
    const send = jest.fn();
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <GroupContext.Provider value={{ send, groupname }}>
          <MemoryRouter>
            <UserItem userItem={user} />
          </MemoryRouter>
        </GroupContext.Provider>
      </AppContext.Provider>
    );
    expect(screen.queryAllByText(/admin/i).length).toBe(3);
  })

  it('Tooltip appears after delete button click', () => {
    const send = jest.fn();
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <GroupContext.Provider value={{ send, groupname }}>
          <MemoryRouter>
            <UserItem userItem={user}  />
          </MemoryRouter>
        </GroupContext.Provider>
      </AppContext.Provider>
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();
  });


  it('Cancel clears tooltip', async () => {
    const send = jest.fn();
    render(
      <AppContext.Provider value={{ user: adminRole }}>
        <GroupContext.Provider value={{ send, groupname }}>
          <MemoryRouter>
            <UserItem userItem={user} />
          </MemoryRouter>
        </GroupContext.Provider>
      </AppContext.Provider>
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();

    const buttons = screen.queryAllByRole('button');
    expect(buttons.length).toBe(3);

    fireEvent.click(buttons[1]);

    const buttonsAfterCancel = screen.queryAllByRole('button');

    expect(buttonsAfterCancel.length).toBe(1);
  });


  it('Submit Button Fires open tooltip and confirm fires delete, send is called to reset table', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <AppContext.Provider value={{ user: adminRole }}>
            <GroupContext.Provider value={{ send, groupname }}>
              <MemoryRouter>
                <UserItem userItem={user} />
              </MemoryRouter>
            </GroupContext.Provider>
          </AppContext.Provider>
        );
      }
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();

    const confirmButton = screen.getByText(/confirm/i);
    fireEvent.click(confirmButton);

    await waitFor(() => jest.mock);

    expect(send).toBeCalledTimes(1);
  });


  it('Admin role has delete option', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <AppContext.Provider value={{ user: adminRole }}>
            <GroupContext.Provider value={{ send, groupname }}>
              <MemoryRouter>
                <UserItem userItem={user} />
              </MemoryRouter>
            </GroupContext.Provider>
          </AppContext.Provider>
        );
      }
    );
    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(1);
  });


  it('Privileged role has delete option', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <AppContext.Provider value={{ user: privilegedRole }}>
            <GroupContext.Provider value={{ send, groupname }}>
              <MemoryRouter>
                <UserItem userItem={user} />
              </MemoryRouter>
            </GroupContext.Provider>
          </AppContext.Provider>
        );
      }
    );
    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(1);
  });


  it('Regular user role does not have delete option', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <AppContext.Provider value={{ user: userRole }}>
            <GroupContext.Provider value={{ send, groupname }}>
              <MemoryRouter>
                <UserItem userItem={user} />
              </MemoryRouter>
            </GroupContext.Provider>
          </AppContext.Provider>
        );
      }
    );
    const buttons = screen.queryAllByRole('button');

    expect(buttons.length).toBe(0);
  });

})
