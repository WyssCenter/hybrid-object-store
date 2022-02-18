// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import AccountLabel from '../AccountLabel';

describe('AccountLabel', () => {
  it('Renders AccountLabel', async () => {
    render(
      <AccountLabel
        label="Username"
        value="admin"
      />
    );
    const label = screen.getByText('Username:');

    expect(label).toBeInTheDocument();

  });

  it('Account label will not render if value is *', async () => {
    render(
      <AccountLabel
        label="Username"
        value="*"
      />
    );
    const label = screen.queryAllByText('Username:');
    expect(label.length).toBe(0);

  });
});
