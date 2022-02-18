// vendor
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
// components
import AddSection from '../AddSection';


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


describe('AddSection', () => {
  it('Renders Add User', () => {
    render(
      <AddSection
        sectionType="user"
      />
    );
    const addUserText = screen.queryAllByText(/Add user/i);
    expect(addUserText[0]).toBeInTheDocument();
    expect(addUserText[1]).toBeInTheDocument();
  });


  it('Input callback updates', () => {
    render(
      <AddSection
        sectionType="user"
      />

    );
    const input = screen.getByRole('textbox');

    fireEvent.keyUp(input, { target: { value: 'name' }});


    expect(input.value).toBe('name');
  });


  it('Add Button is not disabled', () => {
    render(
      <AddSection
        sectionType="user"
      />

    );
    const input = screen.getByRole('textbox');

    fireEvent.keyUp(input, { target: { value: 'name' }});

    const button = screen.getByRole('button');

    expect(button.disabled).toBe(false);
  });


  it('Enter Fires Submit', async () => {
    render(
      <AddSection
        sectionType="user"
      />

    );
    const input = screen.getByRole('textbox');

    fireEvent.keyUp(input, { key: 'Enter', target: { value: 'name' }});

    await waitFor(()=> jest.mock);


    expect(input.value).toBe('');
  });

})
