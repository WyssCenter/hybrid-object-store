// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import Breadcrumbs from '../Breadcrumbs';

test('renders learn react link', () => {
  render(
    <Router>
      <Breadcrumbs />
    </Router>
  );
  const linkElement = screen.getByText(/Server/i);
  expect(linkElement).toBeInTheDocument();
});
