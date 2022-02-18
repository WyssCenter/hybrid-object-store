// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import HierarchyHeader from '../HierarchyHeader';

describe('HierarchyHeader', () => {
  it('Renders button', () => {
    render(
      <HierarchyHeader
        header="Header"
        subheader="Subheader"
      />
    );
    const headerElement = screen.getByText('Header');
    const subheaderElement = screen.getByText('Subheader');
    expect(headerElement).toBeInTheDocument();
    expect(subheaderElement).toBeInTheDocument();
  });

});
