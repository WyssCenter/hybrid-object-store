// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import Toolbar from '../Toolbar';

test('renders learn react link', () => {
  render(
    <Router>
      <Toolbar />
    </Router>
  );
  const linkElement = screen.getByText(/namespace/i);
  expect(linkElement).toBeInTheDocument();
});
