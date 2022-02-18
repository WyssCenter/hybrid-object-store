// vendor
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
// components
import PermissionsItem from '../PermissionsItem';

jest.mock('Environment/createEnvironment', () => {
  return {
    put: () => new Promise((resolve) => {
      resolve(
        {
          ok: true,
          json: () => new Promise((resolve) => { resolve({ success: true})})
        }
      )
    })
  }
});

const adminUserItem = {
  group: {
    group_name: 'admin-hoss-deafult-group'
  },
  permission: 'rw',
}

const regularUserItem = {
  group: {
    group_name: 'user-hoss-deafult-group'
  },
  permission: 'r',
}


const adminGroupItem = {
  group: {
    group_name: 'admin-group'
  },
  permission: 'rw',
}

const userGroupItem = {
  group: {
    group_name: 'user-group'
  },
  permission: 'r',
}


describe('Delete', () => {
  it('Renders Admin user with rw', () => {
    render(
      <PermissionsItem
        item={adminUserItem}
        sectionType="user"
      />

    );
    const name = screen.getByText('rw');
    expect(name).toBeInTheDocument();
  });


  it('Renders Regular user with r', () => {
    render(
      <PermissionsItem
        item={regularUserItem}
        sectionType="user"
      />

    );
    const name = screen.getByText('r');
    expect(name).toBeInTheDocument();
  });

  it('Renders Admin group with rw', () => {
    render(
      <PermissionsItem
        item={adminGroupItem}
        sectionType="user"
      />

    );
    const name = screen.getByText('rw');
    expect(name).toBeInTheDocument();
  });


  it('Renders user group with r', () => {
    render(
      <PermissionsItem
        item={userGroupItem}
        sectionType="user"
      />
    );
    const name = screen.getByText('r');
    expect(name).toBeInTheDocument();
  });


  it('Clicking dropdown opens menu', async () => {
    render(
      <PermissionsItem
        item={userGroupItem}
        sectionType="user"
      />
    );
    const dropdownButton = screen.getByText('r');
    fireEvent.click(dropdownButton);

    expect(screen.getByText(/rw/i)).toBeInTheDocument();
  });


  it('clicking menu item updates permissions', async () => {
    render(
      <PermissionsItem
        item={userGroupItem}
        sectionType="user"
      />
    );
    const dropdownButton = screen.getByText('r');
    fireEvent.click(dropdownButton);

    const menuButton = screen.getByText('rw');
    fireEvent.click(menuButton);

    expect(screen.getByText(/rw/i)).toBeInTheDocument();
  });


})
