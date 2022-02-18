// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import TooltipConfirm from '../TooltipConfirm';


describe('TooltipConfirm', () => {
  it('renders learn react link', () => {
    render(<TooltipConfirm />);
    const linkElement = screen.getByRole('presentation');
    expect(linkElement).toBeInTheDocument();
  });
});
