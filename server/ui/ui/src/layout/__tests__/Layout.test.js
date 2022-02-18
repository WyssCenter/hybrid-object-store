// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// components
import Layout from '../Layout';

test('renders learn react link', () => {
  render(
    <Router>
      <Layout />
    </Router>
  );
  const linkElement = screen.getByAltText(/Gigantum HOS Logo/i);
  expect(linkElement).toBeInTheDocument();
});
