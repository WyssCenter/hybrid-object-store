// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { MemoryRouter } from 'react-router-dom';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
// data
import userData from '../../../__tests__/GroupsData'
// components
import GroupsTable from '../GroupsTable';


describe('GroupsTable', () => {
  it('Renders header divider', () => {
    render(
      <MemoryRouter>
        <GroupsTable
          user={userData}
          authorized
        />
      </MemoryRouter>
    );

    expect(screen.getByText(/Groups/i)).toBeInTheDocument();
  })

  it('Row has all elements', () => {
    render(
      <MemoryRouter>
        <GroupsTable
          user={userData}
          authorized
        />
      </MemoryRouter>
    );


    expect(screen.getByText('group-1')).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeInTheDocument();
    expect(screen.getByText('this is a description')).toBeInTheDocument();
  });

  it('Authorized user has Actions header', async () => {
    act(() => {
        render(
          <MemoryRouter>
            <GroupsTable
              user={userData}
              authorized
            />
          </MemoryRouter>
        );
      }
    );
    expect(screen.queryAllByText('Actions').length).toBe(1)
  });

  it('Unauthorized user does not have Actions header', async () => {
    act(() => {
        render(
          <MemoryRouter>
            <GroupsTable
              user={userData}
              authorized={false}
            />
          </MemoryRouter>
        );
      }
    );
    expect(screen.queryAllByText('Actions').length).toBe(0)
  });

})
