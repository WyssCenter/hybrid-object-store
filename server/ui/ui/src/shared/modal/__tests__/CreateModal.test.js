// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event'
// components
import CreateModal from '../CreateModal';




describe('CreateModal', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders create namespace', () => {
    render(
      <CreateModal
        handleClose={jest.fn}
        isVisible={true}
        modalType="namespace"
        postRoute="namespace"
        updateNamespaceFetchId={jest.fn}
      />
    );
    const createModalElement = screen.getByText(/Create a new namespace here/);
    expect(createModalElement).toBeInTheDocument();
  });

  it('Clicking cancel fires handle close', () => {
    const click = jest.fn();
    render(
      <CreateModal
        handleClose={click}
        isVisible={true}
        modalType="namespace"
        postRoute="namespace"
        updateNamespaceFetchId={jest.fn}
      />
    );
    userEvent.click(screen.getByText(/Cancel/));
    expect(click).toHaveBeenCalledTimes(1);
  });

})
