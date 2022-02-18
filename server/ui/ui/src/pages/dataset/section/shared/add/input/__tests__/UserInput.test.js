// vendor
import fetchMock, { enableFetchMocks } from 'jest-fetch-mock';
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
// components
import UserInput from '../UserInput';

enableFetchMocks();


describe('UserInput', () => {
  it('Renders Users', () => {
    render(
      <UserInput
        inputRef={{ current: <input />}}
        permissionType="user"
        updateName={jest.fn()}
      />
    );
    const addUserText = screen.getByText(/Add User/i);
    expect(addUserText).toBeInTheDocument();
  });


  it('Input callback updates', () => {
    const updateName = jest.fn();
    render(
      <UserInput
        inputRef={{ current: <div />}}
        permissionType="user"
        updateName={updateName}
      />

    );
    const input = screen.getByRole('textbox');

    fireEvent.keyUp(input, { target: { value: 'name' }});


    expect(input.value).toBe('name');
    expect(updateName).toHaveBeenCalledTimes(1);
  });

})
