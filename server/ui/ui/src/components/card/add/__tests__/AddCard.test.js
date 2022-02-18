// vendor
import React from 'react';
import { render, screen, cleanup } from '@testing-library/react';
import userEvent from '@testing-library/user-event'
// components
import AddCard from '../AddCard';

describe('AddCard', () => {
  beforeEach(() => {
    cleanup();
  });

  it('renders namespace', async () => {
    render(
      <AddCard
        updateModalVisible={() => null}
        type="namespace"
      />
    );
    const linkElement = await screen.getByText(/Create Namespace/i);
    expect(linkElement).toBeInTheDocument();
  });


  it('renders dataset', async () => {
    render(
      <AddCard
        updateModalVisible={() => null}
        type="dataset"
      />
    );
    const linkElement = await screen.getByText(/Create Dataset/i);
    expect(linkElement).toBeInTheDocument();
  });


  it('test button', () => {
    const click = jest.fn();
    render(
      <AddCard
        updateModalVisible={click}
        type="namespace"
      />
    );
    userEvent.click(screen.getByRole('button'))
    expect(click).toHaveBeenCalledTimes(1);
  });

});
