// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event'
// components
import NamespaceFields from '../NamespaceFields';




describe('NamespaceFields', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders create namespace', () => {
    const keyboardEvent = jest.fn();
    render(
      <NamespaceFields
        handleBucketNameEvent={keyboardEvent}
        handleObjectChangeEvent={() => null}
        objectStoreList={['default', 'refault']}
      />
    );
    const createModalElement = screen.getByRole('textbox');

    userEvent.type(createModalElement, 'databucket');

    expect(createModalElement.value).toBe('databucket');
  });

  it('Clicking cancel fires handle close', () => {
    const keyboardEvent = jest.fn();
    const click = jest.fn();
    render(
      <NamespaceFields
        handleBucketNameEvent={keyboardEvent}
        handleObjectChangeEvent={click}
        objectStoreList={['default', 'numberone']}
      />
    );
    userEvent.click(screen.getByText(/default/));
    userEvent.click(screen.getByText(/numberone/));
    expect(click).toHaveBeenCalledTimes(2);
  });

})
