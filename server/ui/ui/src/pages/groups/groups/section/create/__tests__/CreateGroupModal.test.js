// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
// components
import CreateGroupModal from '../CreateGroupModal';

jest.mock('Environment/createEnvironment', () => {
  return {
    post: (route, body) => new Promise((resolve) => {
      if (body.name === 'name') {
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


describe('CreateGroupModal', () => {
  it('Renders header divider', () => {
    render(
      <CreateGroupModal />
    );

    expect(screen.getByText(/Create Group/i)).toBeInTheDocument();
  })

  it('Values update in textboxes', () => {
    render(
      <CreateGroupModal />
    );
    const inputName = screen.queryAllByRole('textbox');
    fireEvent.keyUp(inputName[0], { target: { value: 'name' }});
    fireEvent.keyUp(inputName[1], { target: { value: 'this is a description' }});

    expect(inputName[0].value).toBe('name');
    expect(inputName[1].value).toBe('this is a description');
  });


  it('Submit Button Fires post and clears inputs', async () => {
    act(() => {render(
      <CreateGroupModal />
    );
    })
    const inputName = screen.queryAllByRole('textbox');
    fireEvent.keyUp(inputName[0], { target: { value: 'name' }});
    fireEvent.keyUp(inputName[1], { target: { value: 'this is a description' }});
    const button = screen.getByRole('button');
    fireEvent.click(button);
    await waitFor(() => jest.mock)

    expect(inputName[0].value).toBe('');
    expect(inputName[1].value).toBe('');
  });


  it('Submit Button Fires post with error and clears inputs', async () => {
    act(() => {render(
      <CreateGroupModal />
    );
    })
    const inputName = screen.queryAllByRole('textbox');
    fireEvent.keyUp(inputName[0], { target: { value: 'name sasd' }});
    fireEvent.keyUp(inputName[1], { target: { value: 'this is a description' }});
    const button = screen.getByRole('button');
    fireEvent.click(button);
    await waitFor(() => jest.mock)


    const errorElement = screen.getByText('Name value does not match spec');
    expect(errorElement).toBeInTheDocument();
  });

})
