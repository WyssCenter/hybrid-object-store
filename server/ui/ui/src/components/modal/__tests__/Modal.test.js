// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event'
// components
import Modal from '../Modal';


describe('Modal', () => {
  beforeAll(() => {
    const modalRoot = global.document.createElement('div');
    modalRoot.setAttribute('id', 'modal');
    const body = global.document.querySelector('body');
    body.appendChild(modalRoot);
  })

  it('modal renders', () => {
    render(
      <Modal
        handleClose={jest.fn()}
        header="Modal Header"
        overflow={false}
        icon=""
        size="medium"
        subheader="Modal Subheader"
      >
        <p>Children</p>
      </Modal>
    );
    const linkElement = screen.getByText('Children');
    expect(linkElement).toBeInTheDocument();
  });


  it('modal header matches', () => {
    render(
      <Modal
        handleClose={jest.fn()}
        header="Modal Header"
        overflow={false}
        icon=""
        size="medium"
        subheader="Modal Subheader"
      >
        <p>Children</p>
      </Modal>
    );
    const linkElement = screen.getByText('Modal Header');
    expect(linkElement).toBeInTheDocument();
  });



  it('modal subheader matches', () => {
    render(
      <Modal
        handleClose={jest.fn()}
        header="Modal Header"
        overflow={false}
        icon=""
        size="medium"
        subheader="Modal Subheader"
      >
        <p>Children</p>
      </Modal>
    );
    const linkElement = screen.getByText("Modal Subheader");
    expect(linkElement).toBeInTheDocument();
  });



  it('Modal handleClose fires an event', () => {
    const click = jest.fn()
    render(
      <Modal
        handleClose={click}
        header="Modal Header"
        overflow={false}
        icon=""
        size="medium"
        subheader="Modal Subheader"
      >
        <p>Children</p>
      </Modal>
    );
    userEvent.click(screen.getByRole('presentation'))
    expect(click).toHaveBeenCalledTimes(1);
  });
});
