// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import DivisionText from '../DivisionText';

describe('DivisionText', () => {
  it('Renders text', () => {
    render(
      <DivisionText
        text="Header section"
      />);
    const divisionTextItem = screen.getByText('Header section');
    expect(divisionTextItem).toBeInTheDocument();
  });
})
