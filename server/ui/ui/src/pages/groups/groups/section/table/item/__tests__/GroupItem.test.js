// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import GroupsContext from '../../../../GroupsContext';
// components
import GroupItem from '../GroupItem';

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


const membership = {"group":{"group_name":"group-1","description":"this is a description"}}


describe('GroupItem', () => {
  it('Renders group data', () => {
    const send = jest.fn();
    render(
      <GroupsContext.Provider value={{ send }}>
        <MemoryRouter>
          <GroupItem
            membership={membership}
            authorized
          />
        </MemoryRouter>
      </GroupsContext.Provider>
    );

    expect(screen.getByText(/group-1/i)).toBeInTheDocument();
    expect(screen.getByText(/this is a description/i)).toBeInTheDocument();
  })

  it('Tooltip appears after delete button click', () => {
    const send = jest.fn();
    render(
      <GroupsContext.Provider value={{ send }}>
        <MemoryRouter>
          <GroupItem
            membership={membership}
            authorized
          />
        </MemoryRouter>
      </GroupsContext.Provider>
    );
    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/Are you sure/i)).toBeInTheDocument();
  });


  it('Cancel clears tooltip', async () => {
    const send = jest.fn();
    render(
      <GroupsContext.Provider value={{ send }}>
        <MemoryRouter>
          <GroupItem
            membership={membership}
            authorized
          />
        </MemoryRouter>
      </GroupsContext.Provider>
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
          <GroupsContext.Provider value={{ send }}>
            <MemoryRouter>
              <GroupItem
                membership={membership}
                authorized
              />
            </MemoryRouter>
          </GroupsContext.Provider>
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

  it('Authorized user has access to delete button', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <GroupsContext.Provider value={{ send }}>
            <MemoryRouter>
              <GroupItem
                membership={membership}
                authorized
              />
            </MemoryRouter>
          </GroupsContext.Provider>
        );
      }
    );
    expect(screen.queryAllByRole('button').length).toBe(1)
  });

  it('Unauthorized user does not have access to delete button', async () => {
    const send = jest.fn();
    act(() => {
        render(
          <GroupsContext.Provider value={{ send }}>
            <MemoryRouter>
              <GroupItem
                membership={membership}
                authorized={false}
              />
            </MemoryRouter>
          </GroupsContext.Provider>
        );
      }
    );
    expect(screen.queryAllByRole('button').length).toBe(0)
  });

})
