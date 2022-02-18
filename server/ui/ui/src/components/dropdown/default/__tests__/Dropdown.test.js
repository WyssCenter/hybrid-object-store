// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import Dropdown from '../Dropdown';

describe('Dropdown', () => {
  it('Renders dropdown', () => {
    render(
      <Dropdown
        customStyle="item"
        itemAction={jest.fn()}
        label="dropdown"
        listAction={jest.fn()}
        listItems={['rat terrier', 'pug', 'beagle']}
        visibility={true}
      />);
    const dropdownItem = screen.getByText('pug');
    expect(dropdownItem).toBeInTheDocument();
  });
})
